package mysql

import (
	"testing"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/schema"
)

func createTestTable() *schema.Table {
	return &schema.Table{
		Schema: "test_db",
		Name:   "test_table",
		Columns: []schema.TableColumn{
			{Name: "id", Type: schema.TYPE_NUMBER},
			{Name: "name", Type: schema.TYPE_STRING},
			{Name: "email", Type: schema.TYPE_STRING},
			{Name: "age", Type: schema.TYPE_NUMBER},
		},
		PKColumns: []int{0}, // id is primary key
	}
}

func createTestHeader() *replication.EventHeader {
	return &replication.EventHeader{
		Timestamp: 1234567890,
		ServerID:  1,
		LogPos:    1000,
	}
}

func TestMakeInsertEvent(t *testing.T) {
	table := createTestTable()
	header := createTestHeader()

	tests := []struct {
		name     string
		rows     [][]any
		expected []RowChangeEvent
	}{
		{
			name: "single insert",
			rows: [][]any{
				{1, "John Doe", "john@example.com", 30},
			},
			expected: []RowChangeEvent{
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "INSERT",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{1},
					PrimaryKeyColumns: []string{"id"},
					Before:            nil,
					After: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@example.com",
						"age":   30,
					},
				},
			},
		},
		{
			name: "multiple inserts",
			rows: [][]any{
				{1, "John Doe", "john@example.com", 30},
				{2, "Jane Smith", "jane@example.com", 25},
			},
			expected: []RowChangeEvent{
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "INSERT",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{1},
					PrimaryKeyColumns: []string{"id"},
					Before:            nil,
					After: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@example.com",
						"age":   30,
					},
				},
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "INSERT",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{2},
					PrimaryKeyColumns: []string{"id"},
					Before:            nil,
					After: map[string]any{
						"id":    2,
						"name":  "Jane Smith",
						"email": "jane@example.com",
						"age":   25,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &canal.RowsEvent{
				Table:  table,
				Header: header,
				Rows:   tt.rows,
			}

			result, err := makeInsertEvent(e)
			if err != nil {
				t.Fatalf("makeInsertEvent() error = %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("makeInsertEvent() returned %d events, expected %d", len(result), len(tt.expected))
			}

			for i, event := range result {
				expected := tt.expected[i]
				if event.Database != expected.Database {
					t.Errorf("event[%d].Database = %v, expected %v", i, event.Database, expected.Database)
				}
				if event.Table != expected.Table {
					t.Errorf("event[%d].Table = %v, expected %v", i, event.Table, expected.Table)
				}
				if event.Type != expected.Type {
					t.Errorf("event[%d].Type = %v, expected %v", i, event.Type, expected.Type)
				}
				if len(event.PrimaryKey) != len(expected.PrimaryKey) {
					t.Errorf("event[%d].PrimaryKey length = %d, expected %d", i, len(event.PrimaryKey), len(expected.PrimaryKey))
				}
				if event.Before != nil {
					t.Errorf("event[%d].Before should be nil for INSERT", i)
				}
				if len(event.After) != len(expected.After) {
					t.Errorf("event[%d].After length = %d, expected %d", i, len(event.After), len(expected.After))
				}
			}
		})
	}
}

func TestMakeDeleteEvent(t *testing.T) {
	table := createTestTable()
	header := createTestHeader()

	tests := []struct {
		name     string
		rows     [][]any
		expected []RowChangeEvent
	}{
		{
			name: "single delete",
			rows: [][]any{
				{1, "John Doe", "john@example.com", 30},
			},
			expected: []RowChangeEvent{
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "DELETE",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{1},
					PrimaryKeyColumns: []string{"id"},
					Before: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@example.com",
						"age":   30,
					},
					After: nil,
				},
			},
		},
		{
			name: "multiple deletes",
			rows: [][]any{
				{1, "John Doe", "john@example.com", 30},
				{2, "Jane Smith", "jane@example.com", 25},
			},
			expected: []RowChangeEvent{
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "DELETE",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{1},
					PrimaryKeyColumns: []string{"id"},
					Before: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@example.com",
						"age":   30,
					},
					After: nil,
				},
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "DELETE",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{2},
					PrimaryKeyColumns: []string{"id"},
					Before: map[string]any{
						"id":    2,
						"name":  "Jane Smith",
						"email": "jane@example.com",
						"age":   25,
					},
					After: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &canal.RowsEvent{
				Table:  table,
				Header: header,
				Rows:   tt.rows,
			}

			result, err := makeDeleteEvent(e)
			if err != nil {
				t.Fatalf("makeDeleteEvent() error = %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("makeDeleteEvent() returned %d events, expected %d", len(result), len(tt.expected))
			}

			for i, event := range result {
				expected := tt.expected[i]
				if event.Database != expected.Database {
					t.Errorf("event[%d].Database = %v, expected %v", i, event.Database, expected.Database)
				}
				if event.Table != expected.Table {
					t.Errorf("event[%d].Table = %v, expected %v", i, event.Table, expected.Table)
				}
				if event.Type != expected.Type {
					t.Errorf("event[%d].Type = %v, expected %v", i, event.Type, expected.Type)
				}
				if len(event.PrimaryKey) != len(expected.PrimaryKey) {
					t.Errorf("event[%d].PrimaryKey length = %d, expected %d", i, len(event.PrimaryKey), len(expected.PrimaryKey))
				}
				if event.After != nil {
					t.Errorf("event[%d].After should be nil for DELETE", i)
				}
				if len(event.Before) != len(expected.Before) {
					t.Errorf("event[%d].Before length = %d, expected %d", i, len(event.Before), len(expected.Before))
				}
			}
		})
	}
}

func TestMakeUpdateEvent(t *testing.T) {
	table := createTestTable()
	header := createTestHeader()

	tests := []struct {
		name     string
		rows     [][]any
		expected []RowChangeEvent
	}{
		{
			name: "single update",
			rows: [][]any{
				{1, "John Doe", "john@example.com", 30},    // before
				{1, "John Doe", "john@newexample.com", 31}, // after
			},
			expected: []RowChangeEvent{
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "UPDATE",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{1},
					PrimaryKeyColumns: []string{"id"},
					Before: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@example.com",
						"age":   30,
					},
					After: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@newexample.com",
						"age":   31,
					},
				},
			},
		},
		{
			name: "multiple updates",
			rows: [][]any{
				{1, "John Doe", "john@example.com", 30},      // before
				{1, "John Doe", "john@newexample.com", 31},   // after
				{2, "Jane Smith", "jane@example.com", 25},    // before
				{2, "Jane Smith", "jane@newexample.com", 26}, // after
			},
			expected: []RowChangeEvent{
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "UPDATE",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{1},
					PrimaryKeyColumns: []string{"id"},
					Before: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@example.com",
						"age":   30,
					},
					After: map[string]any{
						"id":    1,
						"name":  "John Doe",
						"email": "john@newexample.com",
						"age":   31,
					},
				},
				{
					Database:          "test_db",
					Table:             "test_table",
					Type:              "UPDATE",
					TimeStamp:         1234567890,
					Position:          "1000",
					ServerID:          "1",
					PrimaryKey:        []any{2},
					PrimaryKeyColumns: []string{"id"},
					Before: map[string]any{
						"id":    2,
						"name":  "Jane Smith",
						"email": "jane@example.com",
						"age":   25,
					},
					After: map[string]any{
						"id":    2,
						"name":  "Jane Smith",
						"email": "jane@newexample.com",
						"age":   26,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &canal.RowsEvent{
				Table:  table,
				Header: header,
				Rows:   tt.rows,
			}

			result, err := makeUpdateEvent(e)
			if err != nil {
				t.Fatalf("makeUpdateEvent() error = %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Fatalf("makeUpdateEvent() returned %d events, expected %d", len(result), len(tt.expected))
			}

			for i, event := range result {
				expected := tt.expected[i]
				if event.Database != expected.Database {
					t.Errorf("event[%d].Database = %v, expected %v", i, event.Database, expected.Database)
				}
				if event.Table != expected.Table {
					t.Errorf("event[%d].Table = %v, expected %v", i, event.Table, expected.Table)
				}
				if event.Type != expected.Type {
					t.Errorf("event[%d].Type = %v, expected %v", i, event.Type, expected.Type)
				}
				if len(event.PrimaryKey) != len(expected.PrimaryKey) {
					t.Errorf("event[%d].PrimaryKey length = %d, expected %d", i, len(event.PrimaryKey), len(expected.PrimaryKey))
				}
				if len(event.Before) != len(expected.Before) {
					t.Errorf("event[%d].Before length = %d, expected %d", i, len(event.Before), len(expected.Before))
				}
				if len(event.After) != len(expected.After) {
					t.Errorf("event[%d].After length = %d, expected %d", i, len(event.After), len(expected.After))
				}
			}
		})
	}
}

func TestMakeUpdateEventError(t *testing.T) {
	table := createTestTable()
	header := createTestHeader()

	t.Run("missing after row", func(t *testing.T) {
		e := &canal.RowsEvent{
			Table:  table,
			Header: header,
			Rows: [][]interface{}{
				{1, "John Doe", "john@example.com", 30}, // before only, missing after
			},
		}

		_, err := makeUpdateEvent(e)
		if err == nil {
			t.Fatal("makeUpdateEvent() should return error for missing after row")
		}
		if err.Error() != "missing after row for update event" {
			t.Errorf("makeUpdateEvent() error = %v, expected 'missing after row for update event'", err)
		}
	})
}

func TestMakeEventWithEdgeCases(t *testing.T) {
	table := createTestTable()
	header := createTestHeader()

	t.Run("empty rows", func(t *testing.T) {
		e := &canal.RowsEvent{
			Table:  table,
			Header: header,
			Rows:   [][]interface{}{},
		}

		insertResult, err := makeInsertEvent(e)
		if err != nil {
			t.Fatalf("makeInsertEvent() error = %v", err)
		}
		if len(insertResult) != 0 {
			t.Errorf("makeInsertEvent() returned %d events, expected 0", len(insertResult))
		}

		deleteResult, err := makeDeleteEvent(e)
		if err != nil {
			t.Fatalf("makeDeleteEvent() error = %v", err)
		}
		if len(deleteResult) != 0 {
			t.Errorf("makeDeleteEvent() returned %d events, expected 0", len(deleteResult))
		}

		updateResult, err := makeUpdateEvent(e)
		if err != nil {
			t.Fatalf("makeUpdateEvent() error = %v", err)
		}
		if len(updateResult) != 0 {
			t.Errorf("makeUpdateEvent() returned %d events, expected 0", len(updateResult))
		}
	})

	t.Run("fewer columns than expected", func(t *testing.T) {
		e := &canal.RowsEvent{
			Table:  table,
			Header: header,
			Rows: [][]interface{}{
				{1, "John Doe"}, // missing email and age columns
			},
		}

		result, err := makeInsertEvent(e)
		if err != nil {
			t.Fatalf("makeInsertEvent() error = %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("makeInsertEvent() returned %d events, expected 1", len(result))
		}

		event := result[0]
		if len(event.After) != 2 {
			t.Errorf("event.After length = %d, expected 2", len(event.After))
		}
		if event.After["id"] != 1 {
			t.Errorf("event.After[id] = %v, expected 1", event.After["id"])
		}
		if event.After["name"] != "John Doe" {
			t.Errorf("event.After[name] = %v, expected 'John Doe'", event.After["name"])
		}
	})
}
