package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/luna-duclos/instrumentedsql"
	"github.com/luna-duclos/instrumentedsql/opentracing"
)

func init() {
	sql.Register("traceable-mysql",
		instrumentedsql.WrapDriver(mysql.MySQLDriver{},
			instrumentedsql.WithTracer(opentracing.NewTracer(false)),
			instrumentedsql.WithOmitArgs(),
			instrumentedsql.WithOpsExcluded(instrumentedsql.OpSQLRowsNext),
		),
	)
}
