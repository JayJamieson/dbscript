package mysql

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type mysqlPosition struct {
	pos   mysql.Position
	force bool
}

type RowChangeEvent struct {
	Database          string         `json:"database"`
	Table             string         `json:"table"`
	Type              string         `json:"type"`
	TimeStamp         uint32         `json:"ts"`
	Position          string         `json:"position"`
	ServerID          string         `json:"server_id"`
	PrimaryKey        []any          `json:"pk"`
	PrimaryKeyColumns []string       `json:"pk_columns"`
	Before            map[string]any `json:"before"`
	After             map[string]any `json:"after"`
}

func (l *BinlogListener) OnRow(event *canal.RowsEvent) error {
	var err error
	var events []RowChangeEvent

	switch event.Action {
	case canal.InsertAction:
		events, err = makeInsertEvent(event)
	case canal.DeleteAction:
		events, err = makeDeleteEvent(event)
	case canal.UpdateAction:
		events, err = makeUpdateEvent(event)
	default:
		err = fmt.Errorf("invalid rows action %s", event.Action)
	}

	if err != nil {
		l.cancel()
		return fmt.Errorf("action type %s err %w, closing listener", event.Action, err)
	}

	l.eventCh <- events

	return l.ctx.Err()
}

func (l *BinlogListener) OnXID(header *replication.EventHeader, nextPos mysql.Position) error {
	return nil
}

func (l *BinlogListener) OnGTID(header *replication.EventHeader, gtidEvent mysql.BinlogGTIDEvent) error {
	return nil
}

func (l *BinlogListener) OnRotate(event *replication.EventHeader, rotateEvent *replication.RotateEvent) error {
	pos := mysql.Position{
		Name: string(rotateEvent.NextLogName),
		Pos:  uint32(rotateEvent.Position),
	}

	l.mysqlPositionSaveCh <- mysqlPosition{pos, true}

	return l.ctx.Err()
}

func (l *BinlogListener) OnTableChanged(*replication.EventHeader, string, string) error {
	return nil
}

func (l *BinlogListener) OnDDL(event *replication.EventHeader, pos mysql.Position, _ *replication.QueryEvent) error {
	l.mysqlPositionSaveCh <- mysqlPosition{pos, true}

	return l.ctx.Err()
}

func (l *BinlogListener) OnPosSynced(header *replication.EventHeader, pos mysql.Position, gtid mysql.GTIDSet, force bool) error {
	return nil
}

func (l *BinlogListener) OnRowsQueryEvent(*replication.RowsQueryEvent) error {
	return nil
}

func (l *BinlogListener) String() string {
	return "BinlogListener"
}

func makeUpdateEvent(e *canal.RowsEvent) ([]RowChangeEvent, error) {
	// create variable to hold slice of RowChangeEvent
	events := make([]RowChangeEvent, 0)

	// extract schema and table name from e.Table
	schema := e.Table.Schema
	table := e.Table.Name

	// Process rows in pairs (before, after)
	for i := 0; i < len(e.Rows); i += 2 {
		if i+1 >= len(e.Rows) {
			return nil, fmt.Errorf("missing after row for update event")
		}

		// create Before and After map structures using e.Table.Columns
		before := make(map[string]any)
		after := make(map[string]any)

		// e.Rows[i] is before state
		// e.Rows[i+1] is after state
		// order of column names is order of values on e.Rows
		beforeRow := e.Rows[i]
		afterRow := e.Rows[i+1]

		for colIdx, col := range e.Table.Columns {
			if colIdx < len(beforeRow) {
				before[col.Name] = beforeRow[colIdx]
			}
			if colIdx < len(afterRow) {
				after[col.Name] = afterRow[colIdx]
			}
		}

		// Extract primary key information
		primaryKey := make([]any, 0)
		primaryKeyColumns := make([]string, 0)

		for _, pkIdx := range e.Table.PKColumns {
			if pkIdx < len(e.Table.Columns) {
				primaryKeyColumns = append(primaryKeyColumns, e.Table.Columns[pkIdx].Name)
				if pkIdx < len(beforeRow) {
					primaryKey = append(primaryKey, beforeRow[pkIdx])
				}
			}
		}

		event := RowChangeEvent{
			Database:          schema,
			Table:             table,
			Type:              "UPDATE",
			TimeStamp:         e.Header.Timestamp,
			Position:          fmt.Sprintf("%d", e.Header.LogPos),
			ServerID:          fmt.Sprintf("%d", e.Header.ServerID),
			PrimaryKey:        primaryKey,
			PrimaryKeyColumns: primaryKeyColumns,
			Before:            before,
			After:             after,
		}

		events = append(events, event)
	}

	return events, nil
}

func makeDeleteEvent(e *canal.RowsEvent) ([]RowChangeEvent, error) {
	// create variable to hold slice of RowChangeEvent
	events := make([]RowChangeEvent, 0)

	// extract schema and table name from e.Table
	schema := e.Table.Schema
	table := e.Table.Name

	// Process each deleted row
	for _, row := range e.Rows {
		// create Before map structure using e.Table.Columns
		before := make(map[string]any)

		// For delete events, only before state exists
		// order of column names is order of values on e.Rows
		for colIdx, col := range e.Table.Columns {
			if colIdx < len(row) {
				before[col.Name] = row[colIdx]
			}
		}

		// Extract primary key information
		primaryKey := make([]any, 0)
		primaryKeyColumns := make([]string, 0)

		for _, pkIdx := range e.Table.PKColumns {
			if pkIdx < len(e.Table.Columns) {
				primaryKeyColumns = append(primaryKeyColumns, e.Table.Columns[pkIdx].Name)
				if pkIdx < len(row) {
					primaryKey = append(primaryKey, row[pkIdx])
				}
			}
		}

		event := RowChangeEvent{
			Database:          schema,
			Table:             table,
			Type:              "DELETE",
			TimeStamp:         e.Header.Timestamp,
			Position:          fmt.Sprintf("%d", e.Header.LogPos),
			ServerID:          fmt.Sprintf("%d", e.Header.ServerID),
			PrimaryKey:        primaryKey,
			PrimaryKeyColumns: primaryKeyColumns,
			Before:            before,
			After:             nil, // No after state for delete
		}

		events = append(events, event)
	}

	return events, nil
}

func makeInsertEvent(e *canal.RowsEvent) ([]RowChangeEvent, error) {
	// create variable to hold slice of RowChangeEvent
	events := make([]RowChangeEvent, 0)

	// extract schema and table name from e.Table
	schema := e.Table.Schema
	table := e.Table.Name

	// Process each inserted row
	for _, row := range e.Rows {
		// create After map structure using e.Table.Columns
		after := make(map[string]any)

		// For insert events, only after state exists
		// order of column names is order of values on e.Rows
		for colIdx, col := range e.Table.Columns {
			if colIdx < len(row) {
				after[col.Name] = row[colIdx]
			}
		}

		// Extract primary key information
		primaryKey := make([]any, 0)
		primaryKeyColumns := make([]string, 0)

		for _, pkIdx := range e.Table.PKColumns {
			if pkIdx < len(e.Table.Columns) {
				primaryKeyColumns = append(primaryKeyColumns, e.Table.Columns[pkIdx].Name)
				if pkIdx < len(row) {
					primaryKey = append(primaryKey, row[pkIdx])
				}
			}
		}

		event := RowChangeEvent{
			Database:          schema,
			Table:             table,
			Type:              "INSERT",
			TimeStamp:         e.Header.Timestamp,
			Position:          fmt.Sprintf("%d", e.Header.LogPos),
			ServerID:          fmt.Sprintf("%d", e.Header.ServerID),
			PrimaryKey:        primaryKey,
			PrimaryKeyColumns: primaryKeyColumns,
			Before:            nil, // No before state for insert
			After:             after,
		}

		events = append(events, event)
	}

	return events, nil
}
