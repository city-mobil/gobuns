package tntclustercfg

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/viciious/go-tarantool"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/retry"
	"github.com/city-mobil/gobuns/tntcluster"
)

const (
	DefaultQueryTimeout      = time.Second
	DefaultConnectTimeout    = time.Second
	DefaultSpace             = ""
	DefaultAddrs             = ""
	DefaultUser              = ""
	DefaultPassword          = ""
	DefaultPoolSize          = 5
	DefaultMaxPoolPacketSize = 0
)

type userShardConfig struct {
	addr              *string
	slaves            *string
	queryTimeout      func() *retry.WaitConfig
	connectTimeout    *time.Duration
	space             *string
	user              *string
	password          *string
	retryCfgFn        func() *retry.Config
	poolSize          *int
	maxPoolPacketSize *int
	name              *string
}

func defineUserShardConfig(prefix string) *userShardConfig {
	prefix = config.SanitizePrefix(prefix)

	var (
		name              = config.String(prefix+"name", "", "shard name")
		addr              = config.String(prefix+"addr", DefaultAddrs, "shard master address")
		slaves            = config.String(prefix+"slaves", DefaultAddrs, "comma separated shard slaves addresses")
		queryTimeoutFn    = retry.GetWaitConfig(prefix+"query_timeout", DefaultQueryTimeout)
		connectTimeout    = config.Duration(prefix+"connect_timeout", DefaultConnectTimeout, "shard connect timeout")
		space             = config.String(prefix+"default_space", DefaultSpace, "default space to connect")
		user              = config.String(prefix+"user", DefaultUser, "username to connect")
		password          = config.String(prefix+"password", DefaultPassword, "user password to connect")
		poolSize          = config.Int(prefix+"pool_size", DefaultPoolSize, "connection pool size")
		maxPoolPacketSize = config.Int(prefix+"max_pool_packet_size", DefaultMaxPoolPacketSize, "max pool packet size in bytes")
		retryCfgFn        = retry.GetRetryConfig(prefix + "retries")
	)

	return &userShardConfig{
		addr:              addr,
		slaves:            slaves,
		queryTimeout:      queryTimeoutFn,
		connectTimeout:    connectTimeout,
		space:             space,
		user:              user,
		password:          password,
		poolSize:          poolSize,
		retryCfgFn:        retryCfgFn,
		maxPoolPacketSize: maxPoolPacketSize,
		name:              name,
	}
}

func newShardConfig(userCfg *userShardConfig) *tntcluster.ShardConfig {
	slavesCfg := *userCfg.slaves
	var slaves []string
	if slavesCfg != "" {
		slaves = strings.Split(slavesCfg, ",")
	}

	cfg := &tntcluster.ShardConfig{
		Opts: tarantool.Options{
			ConnectTimeout:    *userCfg.connectTimeout,
			User:              *userCfg.user,
			Password:          *userCfg.password,
			DefaultSpace:      *userCfg.space,
			PoolMaxPacketSize: *userCfg.maxPoolPacketSize,
		},
		MasterAddr:         *userCfg.addr,
		SlaveAddrs:         slaves,
		QueryTimeoutConfig: userCfg.queryTimeout(),
		RetryConfig:        userCfg.retryCfgFn(),
		PoolSize:           *userCfg.poolSize,
		Name:               *userCfg.name,
	}

	shardConfigWithDefaults(cfg)

	return cfg
}

func NewShardConfig(prefix string) func() *tntcluster.ShardConfig {
	userCfg := defineUserShardConfig(prefix)

	return func() *tntcluster.ShardConfig {
		return newShardConfig(userCfg)
	}
}

func NewClusterConfig(prefix string) func() (*tntcluster.ClusterConfig, error) {
	prefix = config.SanitizePrefix(prefix) + "tntcluster"

	return func() (*tntcluster.ClusterConfig, error) {
		defined := config.SubConfigSuffixes(prefix)
		sort.Strings(defined)

		userShardConfigs := make([]*userShardConfig, 0, len(defined))
		for _, key := range defined {
			shardPrefix := fmt.Sprintf("%s.%s", prefix, key)
			userShardConfigs = append(userShardConfigs, defineUserShardConfig(shardPrefix))
		}

		config.ReInit(prefix)

		shards := make([]*tntcluster.ShardConfig, 0, len(userShardConfigs))
		for _, cfg := range userShardConfigs {
			shards = append(shards, newShardConfig(cfg))
		}

		return &tntcluster.ClusterConfig{
			Shards: shards,
		}, nil
	}
}

// OldClusterConfig is a callback for registering cluster config.
//
// It is a deprecated stuff for backporting for legacy stuff users.
func OldClusterConfig(prefix string) func() (*tntcluster.ClusterConfig, error) {
	prefix = config.SanitizePrefix(prefix)

	var (
		addrs             = config.String(prefix+"tntcluster.addr", DefaultAddrs, "tnt cluster addrs")
		slaves            = config.String(prefix+"tntcluster.slave_addrs", "", "tnt cluster slave addrs")
		queryTimeoutFn    = retry.GetWaitConfig(prefix+"tntcluster.query_timeout", DefaultQueryTimeout)
		connectTimeout    = config.Duration(prefix+"tntcluster.connect_timeout", DefaultConnectTimeout, "tnt cluster connect timeout")
		space             = config.String(prefix+"tntcluster.default_space", DefaultSpace, "tnt cluster default space")
		user              = config.String(prefix+"tntcluster.user", DefaultUser, "tnt cluster user")
		password          = config.String(prefix+"tntcluster.password", DefaultPassword, "tnt cluster user password")
		maxPoolPacketSize = config.Int(prefix+"tntcluster.max_pool_packet_size", DefaultMaxPoolPacketSize, "Tnt cluster max pool packet size in bytes")
		poolSize          = config.Int(prefix+"tntcluster.pool_size", 42, "Tnt cluster connection pool size")
		retryCfgFn        = retry.GetRetryConfig(prefix + "tntcluster.retries")
	)

	return func() (*tntcluster.ClusterConfig, error) {
		var sl [][]string
		if *slaves != "" {
			err := json.Unmarshal([]byte(*slaves), &sl)
			if err != nil {
				return nil, err
			}
		}
		masters := strings.Split(*addrs, ",")
		res := &tntcluster.ClusterConfig{}
		for i, v := range masters {
			cfg := &tntcluster.ShardConfig{
				MasterAddr:         v,
				RetryConfig:        retryCfgFn(),
				QueryTimeoutConfig: queryTimeoutFn(),
				Opts: tarantool.Options{
					ConnectTimeout:    *connectTimeout,
					User:              *user,
					Password:          *password,
					DefaultSpace:      *space,
					PoolMaxPacketSize: *maxPoolPacketSize,
				},
				PoolSize: *poolSize,
			}
			if i < len(sl) {
				cfg.SlaveAddrs = sl[i]
			}
			shardConfigWithDefaults(cfg)
			res.Shards = append(res.Shards, cfg)
		}
		return res, nil
	}
}

func shardConfigWithDefaults(cfg *tntcluster.ShardConfig) {
	if cfg == nil {
		cfg = &tntcluster.ShardConfig{}
	}
	if cfg.Opts.QueryTimeout == 0 {
		cfg.Opts.QueryTimeout = DefaultQueryTimeout
	}
	if cfg.Opts.ConnectTimeout == 0 {
		cfg.Opts.ConnectTimeout = DefaultConnectTimeout
	}
}
