package redis

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/redis/redisconfig"
	bunstracing "github.com/city-mobil/gobuns/redis/tracing"
	"github.com/city-mobil/gobuns/tracing"
	"github.com/city-mobil/gobuns/zlog"

	goredis "github.com/go-redis/redis/v8"
)

var (
	cbMasterConnID = 0
	cbFallbackID   = -1
)

var ErrorNoAvailableSlaves = errors.New("no available items")

// node will contain connection and index
// we need to use index for circuit breaker without blocking each command execution with mutex
type node struct {
	item  goredis.Cmdable
	index int
}

// replicas is a concurrent structure which help use replicas
// it controls replicas usage with circuit breaker and round robin algo
type replicas struct {
	items           []*node
	itemsCnt        int
	lastUsedReplica uint32
	cb              barber.Barber
}

// getNode control nodes usage by round robin algo, will return an error if all nodes are not available
func (s *replicas) getNode() (*node, error) {
	var next1 int
	now := time.Now()
	for i := 0; i < s.itemsCnt; i++ {
		next := atomic.AddUint32(&s.lastUsedReplica, 1)
		next1 = (int(next) - 1) % s.itemsCnt
		if s.cb.IsAvailable(next1, now) {
			return s.items[next1], nil
		}
	}
	return nil, ErrorNoAvailableSlaves
}

// standaloneConn contains all connections master, replicas and fallback (nullable connection)
// will control usage of connections with circuit breaker
// can replace nodes if one is not available
type standaloneConn struct {
	fallback *node
	master   *node
	replicas *replicas
	cb       barber.Barber
}

func newStandaloneConn(
	master *node,
	repls []*node,
	fallback *node,
	cb barber.Barber,
) *standaloneConn {
	return &standaloneConn{
		master:   master,
		fallback: fallback,
		replicas: &replicas{
			items:    repls,
			itemsCnt: len(repls),
			cb:       cb,
		},
		cb: cb,
	}
}

func (s *standaloneConn) handleError(index int, err error) {
	if index == cbFallbackID {
		return
	}

	if err != nil && err != ErrorRedisUnavailable {
		s.cb.AddError(index, time.Now())
	}
}

func (s *standaloneConn) getNode() *node {
	node, err := s.replicas.getNode()
	if err != nil {
		return s.getMaster()
	}
	return node
}

func (s *standaloneConn) getMaster() *node {
	if s.cb.IsAvailable(cbMasterConnID, time.Now()) {
		return s.master
	}
	return s.fallback
}

// standalone is struct for standalone connection
type standalone struct { //nolint:golint
	client *standaloneConn
}

// NewStandalone will return standalone instance with custom circuit breaker
func NewStandalone(
	logger zlog.Logger,
	config *redisconfig.RedisStandaloneConfig,
	circuitBreaker barber.Barber,
	prefix string,
) Redis {
	// Configure master connection with tracing hook
	master := goredis.NewClient(config.Master)
	if config.Tracer.WithHook {
		master.AddHook(bunstracing.NewTracingHook(tracing.Redis, config.Master.Addr))
	}
	masterNode := &node{
		item:  master,
		index: cbMasterConnID,
	}

	// Configure items clients with hooks
	replNodes := make([]*node, 0, len(config.Slaves))
	for i, slaveOpts := range config.Slaves {
		slave := goredis.NewClient(slaveOpts)
		if config.Tracer.WithHook {
			slave.AddHook(bunstracing.NewTracingHook(tracing.Redis, slaveOpts.Addr))
		}

		replNodes = append(replNodes, &node{
			item:  slave,
			index: i + 1,
		})
	}

	// Create fallback redis (null object pattern)
	fallback := &node{
		item:  newFallback(),
		index: cbFallbackID,
	}

	// Conn contains master, slave and circuit breaker and can switch them for better performance
	conn := newStandaloneConn(masterNode, replNodes, fallback, circuitBreaker)

	// Instance is default standalone wrapper for redis connection
	instance := newDefaultStandalone(conn)
	return newRedis(logger, instance, prefix)
}

// NewDefaultStandaloneWithPrefix will create default standalone instance and custom key's prefix
// WARNING: Don't use this if you have to use specific circuit breaker handler
// Will create default circuit breaker from redisconfig.RedisStandaloneConfig parameters
func NewDefaultStandaloneWithPrefix(
	logger zlog.Logger,
	config *redisconfig.RedisStandaloneConfig,
	prefix string,
) Redis {
	// Circuit breaker hosts, will contain hosts ids for checking availability
	cbHosts := make([]int, 1, len(config.Slaves)+1)

	// Configure items clients with hooks
	for i := range config.Slaves {
		cbHosts = append(cbHosts, i+1)
	}

	// Create circuit breaker for redis standalone client
	cb := barber.NewBarber(cbHosts, config.CircuitBreaker)
	return NewStandalone(logger, config, cb, prefix)
}

// NewDefaultStandalone creates default standalone instance
// WARNING: Don't use this if you have to use specific circuit breaker handler
// Will create default circuit breaker from redisconfig.RedisStandaloneConfig parameters
func NewDefaultStandalone(logger zlog.Logger, config *redisconfig.RedisStandaloneConfig) Redis {
	return NewDefaultStandaloneWithPrefix(logger, config, "")
}

func newDefaultStandalone(client *standaloneConn) Redis {
	return &standalone{
		client: client,
	}
}

func (r *standalone) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (string, error) {
	node := r.client.getMaster()
	res, err := node.item.Set(ctx, key, value, ttl).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	node := r.client.getMaster()
	res, err := node.item.SetNX(ctx, key, value, ttl).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) IncrBy(ctx context.Context, key string, inc int64) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.IncrBy(ctx, key, inc).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) DecrBy(ctx context.Context, key string, dec int64) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.DecrBy(ctx, key, dec).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) StrLen(ctx context.Context, key string) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.StrLen(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) Get(ctx context.Context, key string) (string, error) {
	node := r.client.getMaster()
	res, err := node.item.Get(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) SAdd(ctx context.Context, key string, members []interface{}, ttl time.Duration) (int64, error) {
	node := r.client.getMaster()
	pipeline := node.item.Pipeline()
	if pipeline == nil {
		return 0, ErrorRedisUnavailable
	}
	result := pipeline.SAdd(ctx, key, members...)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	r.client.handleError(node.index, err)
	return result.Val(), err
}

func (r *standalone) Incr(ctx context.Context, key string) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.Incr(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) Decr(ctx context.Context, key string) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.Decr(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) SMembers(ctx context.Context, key string) ([]string, error) {
	node := r.client.getNode()
	res, err := node.item.SMembers(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	node := r.client.getNode()
	res, err := node.item.MGet(ctx, keys...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) MSet(ctx context.Context, values []interface{}, ttl time.Duration) (string, error) {
	node := r.client.getMaster()
	pipeline := node.item.Pipeline()
	if pipeline == nil {
		return "", ErrorRedisUnavailable
	}
	result := pipeline.MSet(ctx, values...)
	if ttl > 0 {
		for i := 0; i < len(values); i += 2 {
			key := fmt.Sprintf("%v", values[i])
			pipeline.Expire(ctx, key, ttl)
		}
	}
	_, err := pipeline.Exec(ctx)
	r.client.handleError(node.index, err)
	return result.Val(), err
}

func (r *standalone) GeoAdd(ctx context.Context, key string, locations ...*goredis.GeoLocation) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.GeoAdd(ctx, key, locations...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) GeoRadius(ctx context.Context, key string, lat, lng float64, query *goredis.GeoRadiusQuery) ([]goredis.GeoLocation, error) {
	node := r.client.getNode()
	res, err := node.item.GeoRadius(ctx, key, lng, lat, query).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) SCard(ctx context.Context, key string) (int64, error) {
	node := r.client.getNode()
	res, err := node.item.SCard(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	node := r.client.getNode()
	res, err := node.item.SIsMember(ctx, key, member).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.SRem(ctx, key, members...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) Del(ctx context.Context, keys ...string) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.Del(ctx, keys...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.HDel(ctx, key, fields...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) HSet(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (int64, error) {
	node := r.client.getMaster()
	pipeline := node.item.Pipeline()
	if pipeline == nil {
		return 0, ErrorRedisUnavailable
	}
	result := pipeline.HSet(ctx, key, field, value)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	r.client.handleError(node.index, err)
	return result.Val(), err
}

func (r *standalone) HSetNX(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (bool, error) {
	node := r.client.getMaster()
	pipeline := node.item.Pipeline()
	if pipeline == nil {
		return false, ErrorRedisUnavailable
	}
	result := pipeline.HSetNX(ctx, key, field, value)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	r.client.handleError(node.index, err)
	return result.Val(), err
}

func (r *standalone) HGet(ctx context.Context, key, field string) (string, error) {
	node := r.client.getMaster()
	res, err := node.item.HGet(ctx, key, field).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) HMSet(ctx context.Context, key string, values []interface{}, ttl time.Duration) (bool, error) {
	node := r.client.getMaster()
	pipeline := node.item.Pipeline()
	if pipeline == nil {
		return false, ErrorRedisUnavailable
	}
	result := pipeline.HMSet(ctx, key, values)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	r.client.handleError(node.index, err)
	return result.Val(), err
}

func (r *standalone) HLen(ctx context.Context, key string) (int64, error) {
	node := r.client.getNode()
	res, err := node.item.HLen(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) HKeys(ctx context.Context, key string) ([]string, error) {
	node := r.client.getNode()
	res, err := node.item.HKeys(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) ZAdd(ctx context.Context, key string, members []*goredis.Z, ttl time.Duration) (int64, error) {
	node := r.client.getMaster()
	pipeline := node.item.Pipeline()
	if pipeline == nil {
		return 0, ErrorRedisUnavailable
	}
	result := pipeline.ZAdd(ctx, key, members...)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	r.client.handleError(node.index, err)
	return result.Val(), err
}

func (r *standalone) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.ZRem(ctx, key, members...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	node := r.client.getMaster()
	res, err := node.item.ZRemRangeByScore(ctx, key, min, max).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) ZRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) ([]string, error) {
	node := r.client.getNode()
	res, err := node.item.ZRangeByScore(ctx, key, opt).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) ZCard(ctx context.Context, key string) (int64, error) {
	node := r.client.getNode()
	res, err := node.item.ZCard(ctx, key).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	node := r.client.getMaster()
	res, err := node.item.Expire(ctx, key, expiration).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	node := r.client.getMaster()
	res, err := node.item.ExpireAt(ctx, key, tm).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) Exists(ctx context.Context, keys ...string) (int64, error) {
	node := r.client.getNode()
	res, err := node.item.Exists(ctx, keys...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error) {
	node := r.client.getMaster()
	res, err := node.item.Eval(ctx, script, keys, args...).Result()
	r.client.handleError(node.index, err)
	return res, err
}

func (r *standalone) EvalSha(ctx context.Context, sha1 string, keys []string, args []interface{}) (interface{}, error) {
	node := r.client.getMaster()
	res, err := node.item.EvalSha(ctx, sha1, keys, args...).Result()
	r.client.handleError(node.index, err)
	return res, err
}
