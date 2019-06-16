package modules

import (
	"database/sql"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/report"
)

type server struct {
	sql          *sql.DB
	oldSession   *sql.DB
	pubSub       pkg.PubSub
	reportHolder *report.Holder
}

var _server server

func InitServer(app pkg.Application, pajbot1Config *config.Pajbot1Config, reportHolder *report.Holder) error {
	var err error

	_server.sql = app.SQL()
	_server.oldSession, err = sql.Open("mysql", pajbot1Config.SQL.DSN)
	if err != nil {
		return err
	}
	_server.pubSub = app.PubSub()
	_server.reportHolder = reportHolder

	return nil
}
