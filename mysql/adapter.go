package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/city-mobil/gobuns/mysql/mysqlconfig"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

const (
	componentTypeMySQL  = "mysql"
	componentTypeMySQLx = "mysql_x"
)

// Adapter is a wrapper for connectivity between sqlx and sql connectors.
type Adapter interface {
	Begin() (*sql.Tx, error)
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	Close() error
	Conn(context.Context) (*sql.Conn, error)
	Database() *sql.DB
	Driver() driver.Driver
	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	Ping() error
	PingContext(context.Context) error
	Prepare(string) (*sql.Stmt, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	SetConnMaxLifetime(time.Duration)
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
	Stats() sql.DBStats

	// NOTE(a.petrukhin): Some SQLx specific API.

	// BindNamed ...
	BindNamed(string, interface{}) (string, []interface{}, error)
	DriverName() string
	Beginx() (*sqlx.Tx, error)
	GetContext(context.Context, interface{}, string, ...interface{}) error
	NamedExec(string, interface{}) (sql.Result, error)
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
	Queryx(string, ...interface{}) (*sqlx.Rows, error)
	QueryRowx(string, ...interface{}) *sqlx.Row
	Select(interface{}, string, ...interface{}) error
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	Rebind(string) string

	ComponentID() string
	ComponentType() string
	Name() string
}

type sqlAdapter struct {
	cfg *mysqlconfig.DatabaseConfig
	*sql.DB
}

func (s *sqlAdapter) BindNamed(s2 string, i interface{}) (a string, b []interface{}, err error) {
	return "", nil, ErrNotImplemented
}

func (s *sqlAdapter) DriverName() string {
	return "sql"
}

func (s *sqlAdapter) Queryx(s2 string, i ...interface{}) (*sqlx.Rows, error) {
	return nil, ErrNotImplemented
}

func (s *sqlAdapter) QueryRowx(s2 string, i ...interface{}) *sqlx.Row {
	panic(ErrNotImplemented.Error())
}

func (s *sqlAdapter) Select(i interface{}, s2 string, i2 ...interface{}) error {
	return ErrNotImplemented
}

func (s *sqlAdapter) Rebind(s2 string) string {
	panic(ErrNotImplemented.Error())
}

func (s *sqlAdapter) Database() *sql.DB {
	return s.DB
}

func (s *sqlAdapter) GetContext(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return ErrNotImplemented
}

func (s *sqlAdapter) SelectContext(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return ErrNotImplemented
}

func (s *sqlAdapter) NamedExec(_ string, _ interface{}) (sql.Result, error) {
	return nil, ErrNotImplemented
}

func (s *sqlAdapter) NamedExecContext(_ context.Context, _ string, _ interface{}) (sql.Result, error) {
	return nil, ErrNotImplemented
}

func (s *sqlAdapter) Beginx() (*sqlx.Tx, error) {
	return nil, ErrNotImplemented
}

func (s *sqlAdapter) ComponentType() string {
	return componentTypeMySQL
}

func (s *sqlAdapter) ComponentID() string {
	return s.cfg.Addr
}

func (s *sqlAdapter) Name() string {
	return s.cfg.Name
}

func newSQLAdapter(cfg *mysqlconfig.DatabaseConfig) (Adapter, error) { //nolint:gocritic
	dsn := cfg.DSN()
	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return &sqlAdapter{
		cfg: cfg,
		DB:  db,
	}, nil
}

type sqlxAdapter struct {
	cfg *mysqlconfig.DatabaseConfig
	*sqlx.DB
}

func (s *sqlxAdapter) Database() *sql.DB {
	return s.DB.DB
}

func (s *sqlxAdapter) ComponentID() string {
	// NOTE(a.petrukhin): better name can be chosen :)
	return s.cfg.Addr
}

func (s *sqlxAdapter) ComponentType() string {
	return componentTypeMySQLx
}

func (s *sqlxAdapter) Name() string {
	return s.cfg.Name
}

func newSQLXAdapter(cfg *mysqlconfig.DatabaseConfig) (Adapter, error) { //nolint:gocritic
	dsn := cfg.DSN()
	db, err := sqlx.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return &sqlxAdapter{
		cfg: cfg,
		DB:  db,
	}, nil
}
