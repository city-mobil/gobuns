// Package tntcluster contains API to work with Tarantool.
package tntcluster

import (
	"context"
	"errors"
	"hash/crc32"
	"math/rand"
	"sync"
	"time"

	"github.com/city-mobil/gobuns/zlog/glog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/viciious/go-tarantool"

	"github.com/city-mobil/gobuns/retry"
	"github.com/city-mobil/gobuns/tntcluster/helpers"
	"github.com/city-mobil/gobuns/tntcluster/pool"
	"github.com/city-mobil/gobuns/tracing"
)

const (
	// DefaultQueryTimeout describes default timeout for any performed request.
	DefaultQueryTimeout = time.Second

	// DefaultConnectTimeout describes default connection timeout for a single connection initialization.
	DefaultConnectTimeout = time.Second

	// DefaultPoolSize describes default opened connections for a single tarantool host.
	DefaultPoolSize = 5

	// DefaultMaxPoolPacketSize describes default response body size to be
	// put in a sync.Pool for further usage without allocations.
	//
	// Zero value means no limit.
	DefaultMaxPoolPacketSize = 0
)

var (
	// ErrNoAvailableConnections is an error which occurs when
	// invalid shard is chosen during shard selection.
	ErrNoAvailableConnections = errors.New("no connections available")

	// ErrNoAvailableReplica is an error which occurs when no replicas for shard exist.
	ErrNoAvailableReplica = errors.New("no replica found or available")

	// ErrClosed is an error which occurs when all tarantool connections are closed.
	//
	// This error occurs for every attempt to perform a query only after 'Close' method is called.
	ErrClosed = errors.New("connection pool is closed")
)

var tntRetryableErrors = []uint{
	tarantool.ErrNoConnection,
	tarantool.ErrTimeout,
}

var (
	// ExecOnFailedRetry calls every time when a retry was failed.
	ExecOnFailedRetry = func(n uint, addr, reply string) {
		glog.Warn().
			Uint("attempt", n).
			Str("remote_addr", addr).
			Str("reply", reply).
			Msg("failed to exec tarantool query")
	}
)

// ClusterConfig is a configuration for some tarantool cluster.
type ClusterConfig struct {
	Shards []*ShardConfig
}

// ShardConfig is a configuration for a single shard in cluster.
type ShardConfig struct {
	// MasterAddr is a TCP-address for master tarantool database.
	//
	// Only TCP address can be provided. UNIX-sockets do not work.
	MasterAddr string

	// Name is a logical name for some shard.
	//
	// For example 'pechkin-mail-shard1'
	Name string

	// SlaveAddrs contains addresses of slave for master-host tarantool.
	SlaveAddrs []string

	// PoolSize describes amount of simultaneously opened connections for a single
	// host in shard.
	PoolSize int

	// Opts is a some extra tarantool options.
	Opts tarantool.Options

	// RetryConfig is a configuration for retry policy.
	RetryConfig *retry.Config

	// QueryTimeoutConfig is a configuration for query timeout policy.
	QueryTimeoutConfig *retry.WaitConfig
}

// Shard describes a single shard in tarantool cluster.
type Shard interface {
	// Call performs tarantool request without retries for master host.
	Call(context.Context, tarantool.Query) (*tarantool.Result, error)

	// CallReplica performs tarantool request without retries for replica host.
	CallReplica(context.Context, tarantool.Query) (*tarantool.Result, error)

	// CallTolerant performs tarantool request with retries for some master host.
	CallTolerant(context.Context, tarantool.Query) (*tarantool.Result, error)

	// CallReplicaTolerant performs tarantool request with retries for some replica host.
	CallReplicaTolerant(context.Context, tarantool.Query) (*tarantool.Result, error)

	// Close closes all underlying tarantool connections.
	Close()

	// GetMasterConnector returns connection pool associated with master host.
	GetMasterConnector() (pool.ConnectorPool, error)

	// GetSlaveConnectors returns all slave connection pools for shard.
	GetSlaveConnectors() ([]pool.ConnectorPool, error)
}

// ChooseFunc is a callback function for choosing a correct shard for request.
type ChooseFunc func([]Shard) (Shard, error)

// Cluster describes a logical cluster.
type Cluster interface {
	// ChooseShard chooses correct shard for provided choosing callback.
	ChooseShard(ChooseFunc) (Shard, error)

	// GetShards returns all shard associated with cluster.
	GetShards() []Shard

	// Close closes all underlying cluster connections.
	Close() error
}

type shard struct {
	options         tarantool.Options
	masterConnector pool.ConnectorPool
	slaveConnectors []pool.ConnectorPool
	onceCloser      sync.Once
	waiter          *retry.Waiter
	retrier         *retry.Retrier
	stop            chan struct{}
}

// NewShard creates new shard for cluster.
func NewShard(conf *ShardConfig) (Shard, error) {
	conf = shardConfigWithDefaults(conf)
	masterConnector := pool.New(conf.MasterAddr, conf.Name, &conf.Opts, conf.PoolSize)
	slaveConnectors := make([]pool.ConnectorPool, 0, len(conf.SlaveAddrs))
	for _, host := range conf.SlaveAddrs {
		connection := pool.New(host, conf.Name, &conf.Opts, conf.PoolSize)
		slaveConnectors = append(slaveConnectors, connection)
	}

	return &shard{
		options:         conf.Opts,
		masterConnector: masterConnector,
		slaveConnectors: slaveConnectors,
		waiter:          retry.NewWaiter(conf.QueryTimeoutConfig),
		retrier:         retry.New(conf.RetryConfig),
		stop:            make(chan struct{}),
	}, nil
}

// GetMasterConnector returns connection pool associated with master host.
func (s *shard) GetMasterConnector() (pool.ConnectorPool, error) {
	if s.isClosed() {
		return nil, ErrClosed
	}

	return s.masterConnector, nil
}

// GetSlaveConnectors returns all slave connection pools for shard.
func (s *shard) GetSlaveConnectors() ([]pool.ConnectorPool, error) {
	if s.isClosed() {
		return nil, ErrClosed
	}

	return s.slaveConnectors, nil
}

func (s *shard) isClosed() bool {
	select {
	case <-s.stop:
		return true
	default:
	}
	return false
}

func (s *shard) getRandomSlaveConnector() (pool.ConnectorPool, error) {
	if s.isClosed() {
		return nil, ErrClosed
	}

	if len(s.slaveConnectors) == 0 {
		return nil, ErrNoAvailableReplica
	}
	r := rand.Intn(len(s.slaveConnectors)) //nolint:gosec
	return s.slaveConnectors[r], nil
}

func (s *shard) followDBSpan(span opentracing.Span, query tarantool.Query) opentracing.Span {
	operation, statement := helpers.TarantoolCommandAndStatement(query)
	dbspan := tracing.DBSpan{
		Type:      tracing.Tarantool,
		Statament: statement,
	}
	siblingSpan, _ := tracing.StartDBSpanFromContextWithTracer(context.Background(), span.Tracer(), string(operation), dbspan, opentracing.FollowsFrom(span.Context()))
	return siblingSpan
}

func (s *shard) initSpan(ctx context.Context, query tarantool.Query) opentracing.Span {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		return nil
	}
	operation, statement := helpers.TarantoolCommandAndStatement(query)
	dbspan := tracing.DBSpan{
		Type:      tracing.Tarantool,
		Statament: statement,
	}
	span, _ := tracing.StartDBSpanFromContextWithTracer(ctx, parentSpan.Tracer(), string(operation), dbspan)
	return span
}

// Call performs tarantool request without retries for master host.
func (s *shard) Call(ctx context.Context, query tarantool.Query) (*tarantool.Result, error) {
	return s.call(ctx, query, s.GetMasterConnector)
}

// CallReplica performs tarantool request without retries for replica host.
func (s *shard) CallReplica(ctx context.Context, query tarantool.Query) (*tarantool.Result, error) {
	return s.call(ctx, query, s.getRandomSlaveConnector)
}

func (s *shard) call(ctx context.Context, query tarantool.Query, connectorFunc func() (pool.ConnectorPool, error)) (*tarantool.Result, error) {
	var (
		connector pool.ConnectorPool
		result    *tarantool.Result
		err       error
	)

	span := s.initSpan(ctx, query)
	if span != nil {
		defer func() {
			if err != nil {
				ext.Error.Set(span, true)
				span.LogFields(log.Error(err))
			}
			if connector != nil {
				span.SetTag(string(ext.DBInstance), connector.RemoteAddr())
			}
			span.Finish()
		}()
	}

	connector, err = connectorFunc()
	if err != nil {
		return nil, err
	}
	conn, err := connector.Connect()
	if err != nil {
		return nil, err
	}
	result = s.exec(ctx, conn, query, 0)
	if result == nil {
		return &tarantool.Result{}, nil
	}
	if result.Error != nil {
		err = result.Error
		result = nil
	}
	return result, err
}

// CallTolerant performs tarantool request with retries for some master host.
func (s *shard) CallTolerant(ctx context.Context, query tarantool.Query) (*tarantool.Result, error) {
	return s.callTolerant(ctx, query, s.GetMasterConnector)
}

// CallReplicaTolerant performs tarantool request with retries for some replica host.
func (s *shard) CallReplicaTolerant(ctx context.Context, query tarantool.Query) (*tarantool.Result, error) {
	return s.callTolerant(ctx, query, s.getRandomSlaveConnector)
}

func (s *shard) isRetryable(code uint) bool {
	for _, rc := range tntRetryableErrors {
		if rc == code {
			return true
		}
	}
	return false
}

func (s *shard) callTolerant(ctx context.Context, query tarantool.Query, connectorFunc func() (pool.ConnectorPool, error)) (*tarantool.Result, error) {
	var (
		connector pool.ConnectorPool
		conn      *tarantool.Connection
		result    *tarantool.Result
		span      opentracing.Span
		err       error
	)

	span = s.initSpan(ctx, query)
	if span != nil {
		defer func() {
			if err != nil {
				ext.Error.Set(span, true)
				span.LogFields(log.Error(err))
			}
			if connector != nil {
				span.SetTag(string(ext.DBInstance), connector.RemoteAddr())
			}
			span.Finish()
		}()
	}

	iter := uint(0)
	action := func() error {
		connector, err = connectorFunc()
		if err != nil {
			return retry.Unrecoverable(err)
		}

		if s.isClosed() {
			err = ErrClosed
			return retry.Unrecoverable(err)
		}
		if ctx.Err() != nil {
			err = ctx.Err()
			return retry.Unrecoverable(err)
		}

		conn, err = connector.Connect()
		if err != nil {
			return retry.Unrecoverable(err)
		}

		result = s.exec(ctx, conn, query, iter)
		if result.Error != nil {
			if s.isRetryable(result.ErrorCode) {
				if ExecOnFailedRetry != nil {
					ExecOnFailedRetry(iter, connector.RemoteAddr(), result.String())
				}

				return result.Error
			}
			err = result.Error
			return retry.Unrecoverable(err)
		}

		return nil
	}

	onRetry := func(n uint, err error) {
		iter = n + 1
		if span != nil {
			span.LogFields(log.Int("attempt", int(n)))
			ext.Error.Set(span, true)
			span.LogFields(log.Error(err))
			span.Finish()
			span = s.followDBSpan(span, query)
		}
	}

	err = s.retrier.Do(ctx, action, onRetry)
	if err != nil {
		result = nil
	}

	return result, err
}

func (s *shard) exec(ctx context.Context, conn *tarantool.Connection, query tarantool.Query, iter uint) *tarantool.Result {
	execCtx := ctx
	var cancel context.CancelFunc = func() {}

	// wait for merging https://github.com/viciious/go-tarantool/pull/49
	timeout := s.waiter.Get(iter)
	s.options.QueryTimeout = timeout
	if s.options.QueryTimeout != 0 {
		execCtx, cancel = context.WithTimeout(ctx, s.options.QueryTimeout)
	}
	defer cancel()

	return conn.Exec(execCtx, query)
}

func (s *shard) Close() {
	if s.isClosed() {
		return
	}

	s.onceCloser.Do(func() {
		close(s.stop)
	})

	s.masterConnector.Close()
	for _, conn := range s.slaveConnectors {
		conn.Close()
	}
}

type tntCluster struct {
	shards []Shard
}

// NewCluster initializes new cluster for with given configuration.
func NewCluster(conf *ClusterConfig) (Cluster, error) {
	shards := make([]Shard, 0, len(conf.Shards))
	for _, shardConf := range conf.Shards {
		shard, err := NewShard(shardConf)
		if err != nil {
			return nil, err
		}
		shards = append(shards, shard)
	}

	return &tntCluster{
		shards: shards,
	}, nil
}

// ChooseShard chooses correct shard for provided choosing callback.
func (c *tntCluster) ChooseShard(choose ChooseFunc) (Shard, error) {
	return choose(c.shards)
}

// GetShards returns all shard associated with cluster.
func (c *tntCluster) GetShards() []Shard {
	return c.shards
}

// Close closes all underlying cluster connections.
func (c *tntCluster) Close() error {
	for _, s := range c.shards {
		s.Close()
	}
	return nil
}

// ChooseGivenShard is a callback function which returns shard for given shard idx.
//
// This shard is chosen according to idx in configuration array.
// For example, if we have such configuration:
// shards: ["1", "2", "3", "4", "5"], --
// and we provide '0' as idx, we get "1" shard.
//
// Shards order in configuration MUST NOT be changed.
func ChooseGivenShard(idx int) ChooseFunc {
	return func(shards []Shard) (Shard, error) {
		if idx > len(shards) {
			return nil, ErrNoAvailableConnections
		}
		return shards[idx], nil
	}
}

// ChooseRandomShard is a callback function which returns random shard from cluster.
func ChooseRandomShard() ChooseFunc {
	return func(shards []Shard) (Shard, error) {
		if len(shards) == 0 {
			return nil, ErrNoAvailableConnections
		}
		if len(shards) == 1 {
			return shards[0], nil
		}
		r := rand.Intn(len(shards)) //nolint:gosec
		return shards[r], nil
	}
}

// ChooseFirstShard is a callback function which returns first shard from cluster.
//
// For example, if we have such configuration
// shards: ["1", "2", "3", "4", "5"]
// we get "1" shard as a result.
//
// Shards order in configuration MUST NOT be changed.
func ChooseFirstShard() ChooseFunc {
	return func(shards []Shard) (Shard, error) {
		if len(shards) == 0 {
			return nil, ErrNoAvailableConnections
		}
		return shards[0], nil
	}
}

// ChooseShardByCrc32 returns shard based on CRC32 hash result for given key.
func ChooseShardByCrc32(key string) ChooseFunc {
	return func(shards []Shard) (Shard, error) {
		if len(shards) == 0 {
			return nil, ErrNoAvailableConnections
		}
		sum := crc32.ChecksumIEEE([]byte(key))
		l := uint32(len(shards))
		return shards[sum%l], nil
	}
}

func shardConfigWithDefaults(conf *ShardConfig) *ShardConfig { // nolint:interfacer
	if conf == nil {
		conf = &ShardConfig{}
	}

	if conf.SlaveAddrs == nil {
		conf.SlaveAddrs = []string{}
	}

	if conf.PoolSize <= 0 {
		conf.PoolSize = DefaultPoolSize
	}

	if conf.Opts.QueryTimeout <= 0 {
		conf.Opts.QueryTimeout = DefaultQueryTimeout
	}

	if conf.Opts.ConnectTimeout <= 0 {
		conf.Opts.ConnectTimeout = DefaultConnectTimeout
	}

	if conf.Opts.PoolMaxPacketSize < 0 {
		conf.Opts.PoolMaxPacketSize = DefaultMaxPoolPacketSize
	}

	conf.RetryConfig = retry.ConfigWithDefaults(conf.RetryConfig, 0)
	conf.QueryTimeoutConfig = retry.WaitConfigWithDefaults(conf.QueryTimeoutConfig, DefaultQueryTimeout)

	return conf
}
