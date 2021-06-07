package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/zlog"

	"github.com/city-mobil/gobuns/redis/redisconfig"
	bunstracing "github.com/city-mobil/gobuns/redis/tracing"
	"github.com/city-mobil/gobuns/tracing"

	goredis "github.com/go-redis/redis/v8"
)

var (
	cbClusterServerID = 0
)

// cluster declares struct of Redis cluster connection
type cluster struct { //nolint:golint
	logger   zlog.Logger
	client   goredis.Cmdable
	fallback goredis.Cmdable
	cb       barber.Barber
}

func (c *cluster) handleError(err error) {
	if err != nil && err != ErrorRedisUnavailable {
		c.cb.AddError(cbClusterServerID, time.Now())
	}
}

func (c *cluster) getClient() goredis.Cmdable {
	if c.cb.IsAvailable(cbClusterServerID, time.Now()) {
		return c.client
	}
	return c.fallback
}

// NewDefaultCluster creates new redis cluster instance with default dependencies
// WARNING: Don't use this if you have to handle specific traffic!
func NewDefaultCluster(logger zlog.Logger, config *redisconfig.RedisClusterConfig) Redis {
	return NewDefaultClusterWithPrefix(logger, config, "")
}

// NewDefaultClusterWithPrefix creates new redis cluster instance with default dependencies and key's prefix
func NewDefaultClusterWithPrefix(logger zlog.Logger, config *redisconfig.RedisClusterConfig, prefix string) Redis {
	// Configure circuit breaker for redis cluster connection
	hosts := []int{cbClusterServerID}
	cb := barber.NewBarber(hosts, config.CircuitBreaker)

	return NewCluster(logger, config, cb, prefix)
}

// NewCluster creates new redis cluster instance
func NewCluster(
	logger zlog.Logger,
	config *redisconfig.RedisClusterConfig,
	circuitBreaker barber.Barber,
	prefix string,
) Redis {
	// Configure cluster connection with tracing hook
	client := goredis.NewClusterClient(config.Options)
	if config.Tracer.WithHook {
		client.AddHook(bunstracing.NewTracingHook(tracing.RedisCluster, config.Options.Addrs[0]))
	}

	// Create fallback redis (null object pattern)
	fallback := newFallback()

	instance := newDefaultCluster(logger, client, fallback, circuitBreaker)
	return newRedis(logger, instance, prefix)
}

// NewClusterFromClient creates new redis cluster instance
func newDefaultCluster(
	logger zlog.Logger,
	client,
	fallback goredis.Cmdable,
	cb barber.Barber,
) Redis {
	return &cluster{
		logger:   logger,
		client:   client,
		fallback: fallback,
		cb:       cb,
	}
}

func (c *cluster) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (string, error) {
	res, err := c.getClient().Set(ctx, key, value, ttl).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	res, err := c.getClient().SetNX(ctx, key, value, ttl).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) Get(ctx context.Context, key string) (string, error) {
	res, err := c.getClient().Get(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) Incr(ctx context.Context, key string) (int64, error) {
	res, err := c.getClient().Incr(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) IncrBy(ctx context.Context, key string, inc int64) (int64, error) {
	res, err := c.getClient().IncrBy(ctx, key, inc).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) StrLen(ctx context.Context, key string) (int64, error) {
	res, err := c.getClient().StrLen(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) Decr(ctx context.Context, key string) (int64, error) {
	res, err := c.getClient().Decr(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) DecrBy(ctx context.Context, key string, dec int64) (int64, error) {
	res, err := c.getClient().DecrBy(ctx, key, dec).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) SMembers(ctx context.Context, key string) ([]string, error) {
	res, err := c.getClient().SMembers(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	res, err := c.getClient().MGet(ctx, keys...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) GeoRadius(ctx context.Context, key string, lat, lng float64, query *goredis.GeoRadiusQuery) ([]goredis.GeoLocation, error) {
	res, err := c.getClient().GeoRadius(ctx, key, lng, lat, query).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) GeoAdd(ctx context.Context, key string, locations ...*goredis.GeoLocation) (int64, error) {
	res, err := c.getClient().GeoAdd(ctx, key, locations...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) SCard(ctx context.Context, key string) (int64, error) {
	res, err := c.getClient().SCard(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	res, err := c.getClient().SIsMember(ctx, key, member).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	res, err := c.getClient().SRem(ctx, key, members...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) Del(ctx context.Context, keys ...string) (int64, error) {
	res, err := c.getClient().Del(ctx, keys...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	res, err := c.getClient().HDel(ctx, key, fields...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) HLen(ctx context.Context, key string) (int64, error) {
	res, err := c.getClient().HLen(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) HKeys(ctx context.Context, key string) ([]string, error) {
	res, err := c.getClient().HKeys(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	res, err := c.getClient().ZRem(ctx, key, members...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	res, err := c.getClient().ZRemRangeByScore(ctx, key, min, max).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) ZRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) ([]string, error) {
	res, err := c.getClient().ZRangeByScore(ctx, key, opt).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) ZCard(ctx context.Context, key string) (int64, error) {
	res, err := c.getClient().ZCard(ctx, key).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	res, err := c.getClient().Expire(ctx, key, expiration).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	res, err := c.getClient().ExpireAt(ctx, key, tm).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error) {
	res, err := c.getClient().Eval(ctx, script, keys, args...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) EvalSha(ctx context.Context, sha1 string, keys []string, args []interface{}) (interface{}, error) {
	res, err := c.getClient().EvalSha(ctx, sha1, keys, args...).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) SAdd(ctx context.Context, key string, members []interface{}, ttl time.Duration) (int64, error) {
	pipeline := c.getClient().Pipeline()
	if pipeline == nil {
		return 0, ErrorRedisUnavailable
	}
	defer func(pipeline goredis.Pipeliner) {
		err := pipeline.Close()
		if err != nil {
			c.logger.Warn().Err(err).Msg("Redis got an error while close pipeline")
		}
	}(pipeline)

	result := pipeline.SAdd(ctx, key, members...)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	c.handleError(err)
	return result.Val(), err
}

func (c *cluster) MSet(ctx context.Context, values []interface{}, ttl time.Duration) (string, error) {
	pipeline := c.getClient().Pipeline()
	if pipeline == nil {
		return "", ErrorRedisUnavailable
	}
	defer func(pipeline goredis.Pipeliner) {
		err := pipeline.Close()
		if err != nil {
			c.logger.Warn().Err(err).Msg("Got an error while closing pipeline")
		}
	}(pipeline)

	result := pipeline.MSet(ctx, values)

	if ttl > 0 {
		for i := 0; i < len(values); i += 2 {
			pipeline.Expire(ctx, fmt.Sprintf("%v", values[i]), ttl)
		}
	}

	_, err := pipeline.Exec(ctx)
	c.handleError(err)
	return result.Val(), err
}

func (c *cluster) HSet(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (int64, error) { //nolint:dupl
	pipeline := c.getClient().Pipeline()
	if pipeline == nil {
		return 0, ErrorRedisUnavailable
	}
	defer func(pipeline goredis.Pipeliner) {
		err := pipeline.Close()
		if err != nil {
			c.logger.Warn().Err(err).Msg("Got an error while closing pipeline")
		}
	}(pipeline)

	result := pipeline.HSet(ctx, key, field, value)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	c.handleError(err)
	return result.Val(), err
}

func (c *cluster) HSetNX(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (bool, error) { //nolint:dupl
	pipeline := c.getClient().Pipeline()
	if pipeline == nil {
		return false, ErrorRedisUnavailable
	}
	defer func(pipeline goredis.Pipeliner) {
		err := pipeline.Close()
		if err != nil {
			c.logger.Warn().Err(err).Msg("Got an error while closing pipeline")
		}
	}(pipeline)

	result := pipeline.HSetNX(ctx, key, field, value)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	c.handleError(err)
	return result.Val(), err
}

func (c *cluster) HGet(ctx context.Context, key, field string) (string, error) {
	res, err := c.getClient().HGet(ctx, key, field).Result()
	c.handleError(err)
	return res, err
}

func (c *cluster) HMSet(ctx context.Context, key string, values []interface{}, ttl time.Duration) (bool, error) {
	pipeline := c.getClient().Pipeline()
	if pipeline == nil {
		return false, ErrorRedisUnavailable
	}
	defer func(pipeline goredis.Pipeliner) {
		err := pipeline.Close()
		if err != nil {
			c.logger.Warn().Err(err).Msg("Got an error while closing pipeline")
		}
	}(pipeline)

	result := c.getClient().HMSet(ctx, key, values)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	c.handleError(err)
	return result.Val(), err
}

func (c *cluster) ZAdd(ctx context.Context, key string, members []*goredis.Z, ttl time.Duration) (int64, error) {
	pipeline := c.getClient().Pipeline()
	if pipeline == nil {
		return 0, ErrorRedisUnavailable
	}
	defer func(pipeline goredis.Pipeliner) {
		err := pipeline.Close()
		if err != nil {
			c.logger.Warn().Err(err).Msg("Got an error while closing pipeline")
		}
	}(pipeline)

	result := pipeline.ZAdd(ctx, key, members...)
	if ttl > 0 {
		pipeline.Expire(ctx, key, ttl)
	}
	_, err := pipeline.Exec(ctx)
	c.handleError(err)
	return result.Val(), err
}

func (c *cluster) Exists(ctx context.Context, keys ...string) (int64, error) {
	res, err := c.getClient().Exists(ctx, keys...).Result()
	c.handleError(err)
	return res, err
}
