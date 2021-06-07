package redisconfig

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/config"

	goredis "github.com/go-redis/redis/v8"
)

var (
	ErrAddressNotFilled = errors.New("address not filled")
)

const (
	defaultDialTimeout     = 300 * time.Millisecond
	defaultReadTimeout     = 300 * time.Millisecond
	defaultWriteTimeout    = 300 * time.Millisecond
	defaultMaxRedirects    = 8
	defaultMaxRetries      = 3
	defaultMinRetryBackoff = 8 * time.Millisecond
	defaultMaxRetryBackoff = 512 * time.Millisecond
	defaultReadOnly        = true
	defaultRouteByLatency  = false
	defaultRouteRandomly   = true

	defaultPoolSize           = 250
	defaultMinIdleConns       = 0
	defaultPoolTimeout        = 100 * time.Millisecond
	defaultIdleTimeout        = time.Duration(0)
	defaultIdleCheckFrequency = time.Duration(0)
	defaultMaxConnAge         = time.Duration(0)

	defaultDatabase = 0

	defaultTracerWithHook = true
)

type Tracer struct {
	WithHook bool
}

type Retry struct {
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration
}

type redisConfig struct {
	Addr               []string
	Username           string
	Password           string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	MinIdleConns       int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	MaxConnAge         time.Duration

	Retry          *Retry
	CircuitBreaker *barber.Config
	Tracer         *Tracer
}

const configPrefix = "redis"

func newRedisConfig() func() *redisConfig {
	prefixBuilder := strings.Builder{}
	prefixBuilder.WriteString(configPrefix)
	o := func(opt string) string {
		optBuilder := strings.Builder{}
		optBuilder.WriteString(prefixBuilder.String())
		optBuilder.WriteByte('.')
		optBuilder.WriteString(opt)
		return optBuilder.String()
	}

	var (
		addr     = config.StringSlice(o("addr"), nil, "Redis hostname/hostnames")
		username = config.String(o("username"), "", "Redis username")
		password = config.String(o("password"), "", "Redis user password")

		dialTimeout  = config.Duration(o("dial-timeout"), defaultDialTimeout, "Redis dial timeout")
		readTimeout  = config.Duration(o("read-timeout"), defaultReadTimeout, "Redis read timeout")
		writeTimeout = config.Duration(o("write-timeout"), defaultWriteTimeout, "Redis write timeout")

		maxRetries      = config.Int(o("max-retries"), defaultMaxRetries, "Redis max retries")
		minRetryBackoff = config.Duration(o("min-retry-backoff"), defaultMinRetryBackoff, "Redis min retry backoff")
		maxRetryBackoff = config.Duration(o("max-retry-backoff"), defaultMaxRetryBackoff, "Redis max retry backoff")

		poolSize           = config.Int(o("pool-size"), defaultPoolSize, "Redis pool size")
		minIdleConns       = config.Int(o("min-idle-conns"), defaultMinIdleConns, "Redis min idle conns")
		poolTimeout        = config.Duration(o("pool-timeout"), defaultPoolTimeout, "Redis pool timeout")
		idleTimeout        = config.Duration(o("idle-timeout"), defaultIdleTimeout, "Redis idle timeout")
		idleCheckFrequency = config.Duration(o("idle-check-frequency"), defaultIdleCheckFrequency, "Redis idle check frequency")
		maxConnAge         = config.Duration(o("max-conn-age"), defaultMaxConnAge, "Redis max conn age")

		tracerWithHook = config.Bool(o("tracer.with-hook"), defaultTracerWithHook, "Redis tracer enabler")
	)

	cbPrefixBuilder := strings.Builder{}
	cbPrefixBuilder.WriteString(prefixBuilder.String())
	cbPrefixBuilder.WriteString(".barber")
	cbCfgFn := barber.NewConfig(cbPrefixBuilder.String())

	return func() *redisConfig {
		return &redisConfig{
			Addr:               *addr,
			Username:           *username,
			Password:           *password,
			DialTimeout:        *dialTimeout,
			ReadTimeout:        *readTimeout,
			WriteTimeout:       *writeTimeout,
			PoolSize:           *poolSize,
			MinIdleConns:       *minIdleConns,
			PoolTimeout:        *poolTimeout,
			IdleTimeout:        *idleTimeout,
			IdleCheckFrequency: *idleCheckFrequency,
			MaxConnAge:         *maxConnAge,

			Retry: &Retry{
				MaxRetries:      *maxRetries,
				MinRetryBackoff: *minRetryBackoff,
				MaxRetryBackoff: *maxRetryBackoff,
			},
			CircuitBreaker: cbCfgFn(),
			Tracer: &Tracer{
				WithHook: *tracerWithHook,
			},
		}
	}
}

// RedisClusterConfig sd
type RedisClusterConfig struct {
	Options        *goredis.ClusterOptions
	CircuitBreaker *barber.Config
	Tracer         *Tracer
}

// NewClusterConfig returns cluster configuration with default parameters
func NewClusterConfig() func() (*RedisClusterConfig, error) {
	o := func(opt string) string {
		return fmt.Sprintf("%s.%s", "redis.cluster", opt)
	}

	var (
		maxRedirects   = config.Int(o("max-redirects"), defaultMaxRedirects, "Redis cluster max redirects (max moved)")
		readOnly       = config.Bool(o("readonly"), defaultReadOnly, "Redis cluster readonly")
		routeByLatency = config.Bool(o("route-by-latency"), defaultRouteByLatency, "Redis cluster route by latency")
		routeRandomly  = config.Bool(o("route-randomly"), defaultRouteRandomly, "Redis cluster route randomly")
	)

	commonConfig := newRedisConfig()
	return func() (*RedisClusterConfig, error) {
		conf := commonConfig()

		return &RedisClusterConfig{
			Options: &goredis.ClusterOptions{
				Addrs:          conf.Addr,
				MaxRedirects:   *maxRedirects,
				ReadOnly:       *readOnly,
				RouteByLatency: *routeByLatency,
				RouteRandomly:  *routeRandomly,

				Username:           conf.Username,
				Password:           conf.Password,
				DialTimeout:        conf.DialTimeout,
				ReadTimeout:        conf.ReadTimeout,
				WriteTimeout:       conf.WriteTimeout,
				PoolSize:           conf.PoolSize,
				MinIdleConns:       conf.MinIdleConns,
				MaxConnAge:         conf.MaxConnAge,
				PoolTimeout:        conf.PoolTimeout,
				IdleTimeout:        conf.IdleTimeout,
				IdleCheckFrequency: conf.IdleCheckFrequency,
				MaxRetries:         conf.Retry.MaxRetries,
				MinRetryBackoff:    conf.Retry.MinRetryBackoff,
				MaxRetryBackoff:    conf.Retry.MaxRetryBackoff,
			},
			CircuitBreaker: conf.CircuitBreaker,
			Tracer:         conf.Tracer,
		}, nil
	}
}

// RedisStandaloneConfig sd
type RedisStandaloneConfig struct {
	Master         *goredis.Options
	Slaves         []*goredis.Options
	CircuitBreaker *barber.Config
	Tracer         *Tracer
}

// NewStandaloneConfig returns
func NewStandaloneConfig() func() (*RedisStandaloneConfig, error) {
	o := func(opt string) string {
		return fmt.Sprintf("%s.%s", "redis.standalone", opt)
	}

	var (
		replicas = config.StringSlice(o("replicas"), nil, "Redis replicas")
		db       = config.Int(o("db"), defaultDatabase, "Redis database")
	)

	commonConfig := newRedisConfig()
	return func() (*RedisStandaloneConfig, error) {
		conf := commonConfig()

		if len(conf.Addr) == 0 {
			return nil, ErrAddressNotFilled
		}

		var slaves []*goredis.Options
		for _, slave := range *replicas {
			slaves = append(slaves, &goredis.Options{
				Addr: slave,
				DB:   *db,

				Username:           conf.Username,
				Password:           conf.Password,
				DialTimeout:        conf.DialTimeout,
				ReadTimeout:        conf.ReadTimeout,
				WriteTimeout:       conf.WriteTimeout,
				PoolSize:           conf.PoolSize,
				MinIdleConns:       conf.MinIdleConns,
				MaxConnAge:         conf.MaxConnAge,
				PoolTimeout:        conf.PoolTimeout,
				IdleTimeout:        conf.IdleTimeout,
				IdleCheckFrequency: conf.IdleCheckFrequency,
				MaxRetries:         conf.Retry.MaxRetries,
				MinRetryBackoff:    conf.Retry.MinRetryBackoff,
				MaxRetryBackoff:    conf.Retry.MaxRetryBackoff,
			})
		}

		return &RedisStandaloneConfig{
			Master: &goredis.Options{
				Addr: conf.Addr[0],
				DB:   *db,

				Username:           conf.Username,
				Password:           conf.Password,
				DialTimeout:        conf.DialTimeout,
				ReadTimeout:        conf.ReadTimeout,
				WriteTimeout:       conf.WriteTimeout,
				PoolSize:           conf.PoolSize,
				MinIdleConns:       conf.MinIdleConns,
				MaxConnAge:         conf.MaxConnAge,
				PoolTimeout:        conf.PoolTimeout,
				IdleTimeout:        conf.IdleTimeout,
				IdleCheckFrequency: conf.IdleCheckFrequency,
				MaxRetries:         conf.Retry.MaxRetries,
				MinRetryBackoff:    conf.Retry.MinRetryBackoff,
				MaxRetryBackoff:    conf.Retry.MaxRetryBackoff,
			},
			Slaves:         slaves,
			CircuitBreaker: conf.CircuitBreaker,
			Tracer:         conf.Tracer,
		}, nil
	}
}
