package mysql

import (
	"context"
	"database/sql"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/mysql/mysqlconfig"

	// register mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Shard describes a single shard from MySQL cluster.
//
// Shard can be used as a standalone MySQL cluster with master and
// slave connections.
// For multi-master clusters, see Cluster.
type Shard interface {
	// GetMasterConn returns master connection.
	//
	// It is recommended not to use this method in production environment
	// and use it for test purpose only.
	GetMasterConn() Adapter

	// GetSlaveConn returns next slave connection according to ReplicaStrategy.
	//
	// It is recommended not to use this method in production environment
	// and use it for test purpose only.
	GetSlaveConn() Adapter

	// ExecMaster execs query on MySQL master.
	//
	// For further information, see Exec for sql.DB
	ExecMaster(context.Context, string, ...interface{}) (sql.Result, error)

	// ExecMasterTolerant execs query on MySQL master with retries.
	ExecMasterTolerant(context.Context, string, ...interface{}) (sql.Result, error)

	// ExecReplica execs query on MySQL replica.
	//
	// Replica is chosen according to ReplicaStrategy from configuration.
	// For further information, see Exec for sql.DB
	ExecReplica(context.Context, string, ...interface{}) (sql.Result, error)

	// ExecReplicaTolerant execs query on MySQL replica with retries.
	ExecReplicaTolerant(context.Context, string, ...interface{}) (sql.Result, error)

	// QueryMaster performs query on MySQL master.
	//
	// For further information, see Query for sql.DB
	QueryMaster(context.Context, string, ...interface{}) (*sql.Rows, error)

	// QueryMasterTolerant performs query on MySQL master with retries.
	QueryMasterTolerant(context.Context, string, ...interface{}) (*sql.Rows, error)

	// QueryReplica performs query on MySQL replica.
	//
	// Replica is chosen according to ReplicaStrategy configuration.
	// For further information, see Query for sql.DB
	QueryReplica(context.Context, string, ...interface{}) (*sql.Rows, error)

	// QueryReplicaTolerant performs query on MySQL replica with retries.
	QueryReplicaTolerant(context.Context, string, ...interface{}) (*sql.Rows, error)

	// QueryRowMaster performs query on MySQL master and returns only one row.
	//
	// For further information, see QueryRow for sql.DB
	QueryRowMaster(context.Context, string, ...interface{}) *sql.Row

	// QueryRowMasterTolerant performs query on MySQL master with retries,
	// gets at most one row and scans the result.
	QueryRowMasterTolerant(ctx context.Context, query string, args, dest []interface{}) error

	// QueryRowReplica performs query on MySQL replica and returns only one row.
	//
	// For further information, see QueryRow for sql.DB
	QueryRowReplica(context.Context, string, ...interface{}) *sql.Row

	// QueryRowReplicaTolerant performs query on MySQL replica with retries,
	// gets at most one row and scans the result.
	QueryRowReplicaTolerant(ctx context.Context, query string, args, dest []interface{}) error

	// BarberStatus collects circuit-breaker stats.
	BarberStats() *BarberStats

	// Close shutdowns all connections in shard.
	Close() error

	// Setup setups Shard.
	//
	// This method MUST be called in order to use Shard in further.
	Setup() error

	// GetSlaveConnections returns slave connections
	GetSlaveConnections() []Adapter
}

// Cluster defines a MySQL cluster.
//
// Cluster can contain more than one master.
// If you have only one master, use 'Shard' instead.
type Cluster interface {
	ChooseShard(StrategyCallback) Shard
	GetShards() []Shard
}

type shard struct {
	config *mysqlconfig.ShardConfig
	master Adapter
	slaves []Adapter

	cirulnik        barber.Barber
	failStats       failStats
	lastUsedReplica uint32
	connectorType   mysqlconfig.ClusterConnectorType
}

// NewShard creates new not initialized Shard.
func NewShard(cfg *mysqlconfig.ShardConfig, t mysqlconfig.ClusterConnectorType) Shard {
	return NewShardWithCB(cfg, t, nil)
}

// NewShardWithCB creates new not initialized Shard with circuit breaker.
func NewShardWithCB(cfg *mysqlconfig.ShardConfig, t mysqlconfig.ClusterConnectorType, cb barber.Barber) Shard {
	if cfg.RetryConfig == nil {
		cfg.RetryConfig = mysqlconfig.NewDefaultRetryConfig()
	}

	return &shard{
		config:        cfg,
		connectorType: t,
		cirulnik:      cb,
	}
}

// GetMasterConn returns master connection.
//
// It is recommended not to use this method in production environment
// and use it for test purpose only.
func (s *shard) GetMasterConn() Adapter {
	return s.master
}

// GetSlaveConn returns next slave connection according to ReplicaStrategy.
//
// It is recommended not to use this method in production environment
// and use it for test purpose only.
func (s *shard) GetSlaveConn() Adapter {
	conn, _ := s.chooseSlave()
	return conn
}

func (s *shard) ExecMaster(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.master.ExecContext(ctx, query, args...)
}

func (s *shard) ExecMasterTolerant(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	retryCfg := s.config.RetryConfig

	for attempt := 0; attempt <= retryCfg.Max; attempt++ {
		if attempt > 0 {
			if retryCfg.ExecOnErr != nil {
				retryCfg.ExecOnErr(s.master.ComponentID(), err)
			}

			time.Sleep(retryCfg.Timeout)
		}

		res, err = s.master.ExecContext(ctx, query, args...)
		if CanRetry(err) {
			continue
		}

		break
	}

	return
}

func (s *shard) QueryMaster(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.master.QueryContext(ctx, query, args...)
}

func (s *shard) QueryMasterTolerant(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	retryCfg := s.config.RetryConfig

	for attempt := 0; attempt <= retryCfg.Max; attempt++ {
		if attempt > 0 {
			if retryCfg.ExecOnErr != nil {
				retryCfg.ExecOnErr(s.master.ComponentID(), err)
			}

			time.Sleep(retryCfg.Timeout)
		}

		rows, err = s.master.QueryContext(ctx, query, args...)
		if CanRetry(err) {
			continue
		}

		break
	}

	return
}

func (s *shard) QueryRowMaster(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.master.QueryRowContext(ctx, query, args...)
}

func (s *shard) QueryRowMasterTolerant(ctx context.Context, query string, args, dest []interface{}) (err error) {
	retryCfg := s.config.RetryConfig

	for attempt := 0; attempt <= retryCfg.Max; attempt++ {
		if attempt > 0 {
			if retryCfg.ExecOnErr != nil {
				retryCfg.ExecOnErr(s.master.ComponentID(), err)
			}

			time.Sleep(retryCfg.Timeout)
		}

		err = s.master.QueryRowContext(ctx, query, args...).Scan(dest...)
		if CanRetry(err) {
			continue
		}

		break
	}

	return
}

func (s *shard) recordFail(serverID int) {
	if serverID > len(s.config.SlaveConfigs) {
		return
	}
	s.failStats.mu.Lock()
	s.failStats.storage[serverID]++
	s.failStats.mu.Unlock()
}

func (s *shard) chooseRandom() (Adapter, int) { //nolint:gocritic
	var idx int
	now := time.Now()
	for i := 0; i < s.config.MaxBarberAttempts; i++ {
		// NOTE(a.petrukhin): just disable and return the first available slave.
		// It is a common strategy when all the replicas are set up under some kind of
		// proxy such as 'haproxy'. It behaves itself as a Circuit breaker.
		if s.cirulnik == nil {
			break
		}
		// TODO(a.petrukhin): implement disabled replicas.
		idx = rand.Intn(len(s.slaves)) //nolint:gosec
		if s.cirulnik.IsAvailable(idx, now) {
			break
		}

		s.recordFail(idx)
	}
	return s.slaves[idx], idx
}

func (s *shard) chooseRoundRobin() (Adapter, int) { //nolint:gocritic
	// NOTE(a.petrukhin): we do not care about the overflow because it moves to 0 and
	// begins to increase again :)
	var next1 int
	now := time.Now()
	for i := 0; i < s.config.MaxBarberAttempts; i++ {
		// NOTE(a.petrukhin): just disable and return the first available slave.
		// It is a common strategy when all the replicas are set up under some kind of
		// proxy such as 'haproxy'. It behaves itself as a Circuit breaker.
		if s.cirulnik == nil {
			break
		}
		next := atomic.AddUint32(&s.lastUsedReplica, 1)
		next1 = (int(next) - 1) % len(s.slaves)
		if s.cirulnik.IsAvailable(next1, now) {
			break
		}

		s.recordFail(next1)
	}

	// NOTE(a.petrukhin): if we havent found any available replicas after some attempts, we just return the last
	// attempted replica.
	return s.slaves[next1], next1
}

func (s *shard) chooseSlave() (Adapter, int) { //nolint:gocritic
	var conn Adapter
	var connID int
	switch s.config.ReplicaStrategy {
	case mysqlconfig.ReplicaStrategyRoundRobin:
		conn, connID = s.chooseRoundRobin()
	default:
		// NOTE(a.petrukhin): if no strategy chosen :)
		conn, connID = s.chooseRandom()
	}

	return conn, connID
}

func (s *shard) ExecReplica(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	conn, connID := s.chooseSlave()

	res, err := conn.ExecContext(ctx, query, args...)
	if err != nil && s.cirulnik != nil {
		// NOTE(a.petrukhin): we record all errors.
		s.cirulnik.AddError(connID, time.Now())
	}
	return res, err
}

func (s *shard) ExecReplicaTolerant(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	retryCfg := s.config.RetryConfig

	for attempt := 0; attempt <= retryCfg.Max; attempt++ {
		conn, connID := s.chooseSlave()

		if attempt > 0 {
			if retryCfg.ExecOnErr != nil {
				retryCfg.ExecOnErr(conn.ComponentID(), err)
			}

			time.Sleep(retryCfg.Timeout)
		}

		res, err = conn.ExecContext(ctx, query, args...)
		if err != nil && s.cirulnik != nil {
			s.cirulnik.AddError(connID, time.Now())
		}

		if CanRetry(err) {
			continue
		}

		break
	}

	return
}

func (s *shard) QueryReplica(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	conn, connID := s.chooseSlave()

	res, err := conn.QueryContext(ctx, query, args...)
	if err != nil && s.cirulnik != nil {
		s.cirulnik.AddError(connID, time.Now())
	}
	return res, err
}

func (s *shard) QueryReplicaTolerant(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error) {
	retryCfg := s.config.RetryConfig

	for attempt := 0; attempt <= retryCfg.Max; attempt++ {
		conn, connID := s.chooseSlave()

		if attempt > 0 {
			if retryCfg.ExecOnErr != nil {
				retryCfg.ExecOnErr(conn.ComponentID(), err)
			}

			time.Sleep(retryCfg.Timeout)
		}

		rows, err = conn.QueryContext(ctx, query, args...)
		if err != nil && s.cirulnik != nil {
			s.cirulnik.AddError(connID, time.Now())
		}

		if CanRetry(err) {
			continue
		}

		break
	}

	return
}

func (s *shard) QueryRowReplica(ctx context.Context, query string, args ...interface{}) *sql.Row {
	conn, _ := s.chooseSlave()

	// NOTE(a.petrukhin): we can not get error from the row even if it has occurred.
	// So it is recommended not to use this method.
	// Use QueryReplica instead.
	res := conn.QueryRowContext(ctx, query, args...)
	return res
}

func (s *shard) QueryRowReplicaTolerant(ctx context.Context, query string, args, dest []interface{}) (err error) {
	retryCfg := s.config.RetryConfig

	for attempt := 0; attempt <= retryCfg.Max; attempt++ {
		conn, connID := s.chooseSlave()

		if attempt > 0 {
			if retryCfg.ExecOnErr != nil {
				retryCfg.ExecOnErr(conn.ComponentID(), err)
			}

			time.Sleep(retryCfg.Timeout)
		}

		err = conn.QueryRowContext(ctx, query, args...).Scan(dest...)
		if err != nil && err != sql.ErrNoRows && s.cirulnik != nil {
			s.cirulnik.AddError(connID, time.Now())
		}

		if CanRetry(err) {
			continue
		}

		break
	}

	return
}

func (s *shard) setupMaster() error {
	var conn Adapter
	var err error

	masterConfig := s.config.MasterConfig.WithDefaults()

	if s.connectorType == mysqlconfig.ClusterConnectorTypeSQLx {
		conn, err = newSQLXAdapter(masterConfig)
	} else {
		conn, err = newSQLAdapter(masterConfig)
	}

	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}

	s.master = conn
	return nil
}

func (s *shard) stopMaster() {
	_ = s.master.Close()
}

func (s *shard) stopReplicas() {
	for _, v := range s.slaves {
		if v != nil {
			_ = v.Close()
		}
	}
}

func (s *shard) setupReplicas() error {
	for _, cfg := range s.config.SlaveConfigs {
		slaveConfig := cfg.WithDefaults()

		var conn Adapter
		var err error
		if s.connectorType == mysqlconfig.ClusterConnectorTypeSQLx {
			conn, err = newSQLXAdapter(slaveConfig)
		} else {
			conn, err = newSQLAdapter(slaveConfig)
		}

		if err != nil {
			return err
		}

		err = conn.Ping()
		if err != nil {
			return err
		}

		s.slaves = append(s.slaves, conn)
	}

	return nil
}

// Setup setups Shard.
//
// This method MUST be called in order to use Shard in further.
func (s *shard) Setup() error {
	err := s.setupMaster()
	if err != nil {
		return err
	}

	err = s.setupReplicas()
	if err != nil {
		s.stopMaster()
		s.stopReplicas()

		return err
	}

	return nil
}

// Close shutdowns all connections in shard.
func (s *shard) Close() error {
	if err := s.master.Close(); err != nil {
		return err
	}

	for _, v := range s.slaves {
		if err := v.Close(); err != nil {
			return err
		}
	}
	return nil
}

// BarberStatus collects circuit-breaker stats.
func (s *shard) BarberStats() *BarberStats {
	st := s.cirulnik.Stats()

	mp := make(map[string]int)
	for _, v := range st.Hosts {
		h := s.config.SlaveConfigs[v.ServerID].Addr
		mp[h] = v.FailsCount
	}

	stats := &BarberStats{
		CirculStats: mp,
		FailStats:   make(map[string]int),
	}

	s.failStats.mu.RLock()
	defer s.failStats.mu.RUnlock()

	for k, v := range s.failStats.storage {
		h := s.config.SlaveConfigs[k].Addr
		stats.FailStats[h] = v
		s.failStats.storage[k] = 0 // NOTE(a.petrukhin): reset
	}

	return stats
}

func (s *shard) GetSlaveConnections() []Adapter {
	return s.slaves
}
