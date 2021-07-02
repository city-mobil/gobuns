package mysqlconfig

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/config"
)

const (
	defaultAddr              = "127.0.0.1:3306"
	defaultUser              = "guest"
	defaultPassword          = "guest"
	defaultDBName            = "db"
	defaultReadTimeout       = time.Second
	defaultWriteTimeout      = time.Second
	defaultTimeout           = time.Second
	defaultCharset           = "utf8mb4"
	defaultInterpolateParams = true
	defaultSQLDriver         = "mysql"
	defaultBarberAttempts    = 5
	defaultMaxOpenConns      = 200
	defaultMaxIdleConns      = 200
	defaultConnMaxLifetime   = 40 * time.Second
	defaultTimeZone          = ""
)

// Default retry strategy configuration.
const (
	defaultMaxRetries   = 2
	defaultRetryTimeout = 10 * time.Millisecond
)

// ReplicaStrategy describes algorithm for choosing available replica.
type ReplicaStrategy int

const (
	// ReplicaStrategyRandom is a random slave choosing algorithm.
	ReplicaStrategyRandom ReplicaStrategy = iota

	// ReplicaStrategyRoundRobin is a round-robin slave choosing algorithm.
	ReplicaStrategyRoundRobin
)

// ClusterConnectorType describes underlying MySQL driver.
type ClusterConnectorType string

const (
	// ClusterConnectorTypeSQL is a 'sql' driver.
	ClusterConnectorTypeSQL ClusterConnectorType = "sql"

	// ClusterConnectorTypeSQLx is a 'sqlx' driver.
	ClusterConnectorTypeSQLx ClusterConnectorType = "sqlx"
)

// TraceConfig holds configuration for auto tracing.
// By default all options are set to false intentionally.
type TraceConfig struct {
	// MetricQueryTable, if set to true,
	// will enable recording of sql query tables in metrics.
	MetricQueryTable bool

	// MetricQueryOperation, if set to true,
	// will enable recording of sql query type in metrics.
	MetricQueryOperation bool
}

// ClusterConfig is a configuration for MySQL cluster.
type ClusterConfig struct {
	Shards []*ShardConfig
}

// ShardConfig is a configuration of one MySQL shard.
type ShardConfig struct {
	// MasterConfig is a cluster master database configuration.
	MasterConfig *DatabaseConfig

	// SlaveConfigs is a cluster slaves configuration.
	SlaveConfigs []*DatabaseConfig

	// RetryConfig is a query retry configuration.
	RetryConfig *RetryConfig

	// ReplicaStrategy is a replica choosing strategy.
	ReplicaStrategy ReplicaStrategy

	// CircuitBreakerConfig is a configuration for slaves circuit-breaker.
	CircuitBreakerConfig *barber.Config

	// MaxBarberAttempts is a number of retries for choosing available replica.
	MaxBarberAttempts int

	// TraceConfig is an auto tracing configuration.
	TraceConfig *TraceConfig
}

func (s *ShardConfig) withDefaults() *ShardConfig { //nolint:unused
	var newS *ShardConfig

	if s == nil {
		newS = &ShardConfig{}
	} else {
		newS = s
	}

	if newS.MaxBarberAttempts == 0 {
		newS.MaxBarberAttempts = defaultBarberAttempts
	}

	return newS
}

// DatabaseConfig is a single MySQL instance configuration.
type DatabaseConfig struct {
	Addr              string
	Username          string
	Password          string
	DatabaseName      string
	Charset           string
	Collation         string
	VitessReplicaType string
	Driver            string

	// Name describes some business-logic database name
	//
	// For example, 'coupon' database.
	Name              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	Timeout           time.Duration
	ParseTime         bool
	InterpolateParams bool
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	Timezone          string
}

func NewDefaultDatabaseConfig() *DatabaseConfig {
	d := &DatabaseConfig{}

	d.Driver = defaultSQLDriver

	d.Timeout = defaultTimeout
	d.WriteTimeout = defaultWriteTimeout
	d.ReadTimeout = defaultReadTimeout

	d.Charset = defaultCharset
	d.InterpolateParams = defaultInterpolateParams

	d.MaxOpenConns = defaultMaxOpenConns
	d.MaxIdleConns = defaultMaxIdleConns
	d.ConnMaxLifetime = defaultConnMaxLifetime

	d.Timezone = defaultTimeZone

	return d
}

func NewDatabaseConfig(prefix string) func() *DatabaseConfig {
	n := func(opt string) string {
		return fmt.Sprintf("%s.%s", prefix, opt)
	}

	var (
		addr     = config.String(n("addr"), defaultAddr, "MySQL network address")
		user     = config.String(n("user"), defaultUser, "MySQL username")
		password = config.String(n("password"), defaultPassword, "MySQL user password")
		dbName   = config.String(n("dbname"), defaultDBName, "MySQL database name")
		driver   = config.String(n("driver"), defaultSQLDriver, "MySQL database driver")
		timezone = config.String(n("timezone"), defaultTimeZone, "MySQL Database Timezone")

		timeout      = config.Duration(n("timeout"), defaultTimeout, "MySQL connection timeout")
		readTimeout  = config.Duration(n("read_timeout"), defaultReadTimeout, "MySQL I/O read timeout")
		writeTimeout = config.Duration(n("write_timeout"), defaultWriteTimeout, "MySQL I/O write timeout")

		maxOpenConnections    = config.Int(n("pool.max_open_connections"), defaultMaxOpenConns, "The maximum number of open connections to the database")
		maxIdleConnections    = config.Int(n("pool.max_idle_connections"), defaultMaxIdleConns, "the maximum number of connections in the idle connection pool")
		connectionMaxLifetime = config.Duration(n("pool.max_life_time"), defaultConnMaxLifetime, "The maximum amount of time a connection may be reused")
	)

	return func() *DatabaseConfig {
		cfg := NewDefaultDatabaseConfig()
		cfg.Driver = *driver
		cfg.Addr = *addr
		cfg.Username = *user
		cfg.Password = *password
		cfg.DatabaseName = *dbName
		cfg.Timeout = *timeout
		cfg.ReadTimeout = *readTimeout
		cfg.WriteTimeout = *writeTimeout
		cfg.MaxOpenConns = *maxOpenConnections
		cfg.MaxIdleConns = *maxIdleConnections
		cfg.ConnMaxLifetime = *connectionMaxLifetime
		cfg.Timezone = *timezone

		return cfg
	}
}

func (d *DatabaseConfig) WithDefaults() *DatabaseConfig {
	if d == nil {
		d = NewDefaultDatabaseConfig()
	}
	if d.MaxIdleConns == 0 {
		d.MaxIdleConns = defaultMaxIdleConns
	}
	if d.Timeout == 0 {
		d.Timeout = defaultTimeout
	}
	if d.WriteTimeout == 0 {
		d.WriteTimeout = defaultWriteTimeout
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = defaultReadTimeout
	}
	if d.ConnMaxLifetime == 0 {
		d.ConnMaxLifetime = defaultConnMaxLifetime
	}

	return d
}

func escapeTimezone(tz string) string {
	return url.QueryEscape(tz)
}

func (d *DatabaseConfig) DSN() string {
	d = d.WithDefaults()

	builder := strings.Builder{}
	builder.WriteString(d.Username)
	builder.WriteByte(':')
	builder.WriteString(d.Password)
	builder.WriteString("@tcp(")
	builder.WriteString(d.Addr)
	builder.WriteString(")/")
	builder.WriteString(d.DatabaseName)
	if d.VitessReplicaType != "" {
		builder.WriteByte('@')
		builder.WriteString(d.VitessReplicaType)
	}
	builder.WriteString("?timeout=")
	builder.WriteString(d.Timeout.String())
	builder.WriteString("&readTimeout=")
	builder.WriteString(d.ReadTimeout.String())
	builder.WriteString("&writeTimeout=")
	builder.WriteString(d.WriteTimeout.String())
	builder.WriteString("&interpolateParams=")
	builder.WriteString(strconv.FormatBool(d.InterpolateParams))
	if d.Charset != "" {
		builder.WriteString("&charset=")
		builder.WriteString(d.Charset)
	}
	builder.WriteString("&parseTime=")
	builder.WriteString(strconv.FormatBool(d.ParseTime))
	if d.Collation != "" {
		builder.WriteString("&collation=")
		builder.WriteString(d.Collation)
	}

	if d.Timezone != "" {
		builder.WriteString("&loc=")
		builder.WriteString(escapeTimezone(d.Timezone))
	}

	return builder.String()
}

type RetryConfig struct {
	// Max retries before give up.
	Max int
	// Timeout is a duration to wait before trying another attempt.
	Timeout time.Duration
	// ExecOnErr executes before query retry.
	ExecOnErr func(host string, err error)
}

func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		Max:       defaultMaxRetries,
		Timeout:   defaultRetryTimeout,
		ExecOnErr: nil,
	}
}

func NewRetryConfig(prefix string) func() *RetryConfig {
	n := func(opt string) string {
		return fmt.Sprintf("%s.%s", prefix, opt)
	}

	var (
		retryMax     = config.Int(n("max"), defaultMaxRetries, "The maximum number of attempts before give up")
		retryTimeout = config.Duration(n("timeout"), defaultRetryTimeout, "Wait given period before trying another attempt")
	)

	return func() *RetryConfig {
		return &RetryConfig{
			Max:       *retryMax,
			Timeout:   *retryTimeout,
			ExecOnErr: nil,
		}
	}
}

func NewWithSingleNode(cfg *DatabaseConfig, retryCfg *RetryConfig) *ShardConfig {
	if retryCfg == nil {
		retryCfg = NewDefaultRetryConfig()
	}

	return &ShardConfig{
		MasterConfig:    cfg,
		SlaveConfigs:    []*DatabaseConfig{cfg},
		RetryConfig:     retryCfg,
		ReplicaStrategy: ReplicaStrategyRandom,
	}
}
