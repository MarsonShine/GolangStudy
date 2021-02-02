package dao

import (
	"context"
	"sync"

	"kratos-demo/internal/model"

	"github.com/MSLibs/glogger"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/database/sql"
)

var logStd *glogger.GLogger

var oncer sync.Once

func CreateLogger() *glogger.GLogger {
	oncer.Do(func() {
		log := glogger.CreateLog(glogger.GLoggerConfig{})
		logStd = &log
	})
	return logStd
}

func NewDB() (db *sql.DB, cf func(), err error) {
	var (
		cfg sql.Config
		ct  paladin.TOML
	)
	if err = paladin.Get("db.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = sql.NewMySQL(&cfg)
	cf = func() { db.Close() }
	return
}

func (d *dao) RawArticle(ctx context.Context, id int64) (art *model.Article, err error) {
	// get data from db
	return
}
