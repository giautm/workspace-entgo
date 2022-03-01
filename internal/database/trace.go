package database

import (
	"context"
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/go-sql-driver/mysql"
)

// TraceConnector is a wrapper for driver.Driver
type TraceConnector struct {
	DSN string
}

// Connect implements the driver.Connector interface.
func (c TraceConnector) Connect(context.Context) (driver.Conn, error) {
	return c.Driver().Open(c.DSN)
}

// Driver implements the driver.Connector interface.
func (TraceConnector) Driver() driver.Driver {
	return ocsql.Wrap(
		mysql.MySQLDriver{},
		ocsql.WithAllTraceOptions(),
		ocsql.WithRowsClose(false),
		ocsql.WithRowsNext(false),
		ocsql.WithDisableErrSkip(true),
	)
}
