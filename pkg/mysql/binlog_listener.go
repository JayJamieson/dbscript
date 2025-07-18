package mysql

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
)

type BinlogListener struct {
	canal  *canal.Canal
	Logger *slog.Logger

	myslqPosition mysql.Position
	lastSave      time.Time

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc

	eventCh             chan []RowChangeEvent
	mysqlPositionSaveCh chan mysqlPosition
}

type BinlogListenerOptions struct {
	Host     string
	Port     int
	User     string
	Schema   string
	Tables   []string
	Password string
}

func NewBinlogListener(opt *BinlogListenerOptions) (*BinlogListener, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", opt.Host, opt.Port)
	cfg.User = opt.User
	cfg.Password = opt.Password
	cfg.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	// disable dumping
	// does not work on mysql >8.x
	cfg.Dump.ExecutionPath = ""

	for _, table := range opt.Tables {
		cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, opt.Schema+"\\."+table)
	}

	canal, err := canal.NewCanal(cfg)

	if err != nil {
		return nil, err
	}

	coords, err := canal.GetMasterPos()

	if err != nil {
		canal.Close()
		return nil, err
	}

	listener := &BinlogListener{}
	listener.canal = canal
	listener.Logger = logger
	listener.myslqPosition = coords

	listener.eventCh = make(chan []RowChangeEvent, 4096)
	listener.mysqlPositionSaveCh = make(chan mysqlPosition, 4096)
	listener.ctx, listener.cancel = context.WithCancel(context.Background())

	canal.SetEventHandler(listener)

	return listener, nil
}

func (l *BinlogListener) Listen() error {
	return l.canal.RunFrom(l.myslqPosition)
}

func (l *BinlogListener) Close() {
	l.Logger.Info("Closing dbscript")

	l.cancel()

	l.canal.Close()
}

func (l *BinlogListener) GetEventStream() <-chan []RowChangeEvent {
	return l.eventCh
}

func (l *BinlogListener) GetSavepointStream() <-chan mysqlPosition {
	return l.mysqlPositionSaveCh
}
