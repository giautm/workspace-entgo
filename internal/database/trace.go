package database

import (
	"context"
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/go-sql-driver/mysql"
)

type TraceConnector struct {
	DSN string
}

func (c TraceConnector) Connect(context.Context) (driver.Conn, error) {
	return c.Driver().Open(c.DSN)
}

func (TraceConnector) Driver() driver.Driver {
	return ocsql.Wrap(
		mysql.MySQLDriver{},
		ocsql.WithAllTraceOptions(),
		ocsql.WithRowsClose(false),
		ocsql.WithRowsNext(false),
		ocsql.WithDisableErrSkip(true),
	)
}
