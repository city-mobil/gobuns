package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/city-mobil/gobuns/zlog"

	goredis "github.com/go-redis/redis/v8"
)

// ZAddType struct for ZAdd* operations
type ZAddType struct {
	Key     string
	TTL     time.Duration
	Members []*goredis.Z
}

// Redis is basic interface for all commands
type Redis interface {
	// Set value under key
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (string, error)

	// SetNX value under key if not exists
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)

	// Get value by key
	Get(ctx context.Context, key string) (string, error)

	// Del delete a key
	Del(ctx context.Context, keys ...string) (int64, error)

	// Exists determine if a key exists
	Exists(ctx context.Context, keys ...string) (int64, error)

	// Incr atomic +1 to value
	Incr(ctx context.Context, key string) (int64, error)

	// IncrBy atomic +inc to value
	IncrBy(ctx context.Context, key string, inc int64) (int64, error)

	// Decr atomic -1 to value
	Decr(ctx context.Context, key string) (int64, error)

	// DecrBy atomic -dec to value
	DecrBy(ctx context.Context, key string, dec int64) (int64, error)

	// Expire set a key's time to live in seconds
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)

	// ExpireAt set the expiration for a key as a UNIX timestamp
	ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)

	// StrLen returns the length of the string value stored at key
	StrLen(ctx context.Context, key string) (int64, error)

	// MGet get the values of all the given keys
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)

	// MSet set multiple keys to multiple values
	MSet(ctx context.Context, values []interface{}, ttl time.Duration) (string, error)

	// HDel delete one or more hash fields
	HDel(ctx context.Context, key string, fields ...string) (int64, error)

	// HLen get the number of fields in a hash
	HLen(ctx context.Context, key string) (int64, error)

	// HKeys get all the fields in a hash
	HKeys(ctx context.Context, key string) ([]string, error)

	// HSet set the string value of a hash field
	HSet(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (int64, error)

	// HSetNX set the value of a hash field, only if the field does not exist
	HSetNX(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (bool, error)

	// HGet get the value of a hash field
	HGet(ctx context.Context, key, field string) (string, error)

	// HMSet set multiple hash fields to multiple values
	HMSet(ctx context.Context, key string, values []interface{}, ttl time.Duration) (bool, error)

	// ZRem remove one or more members from a sorted set
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)

	// ZRemRangeByScore remove all members in a sorted set within the given scores
	ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error)

	// ZRangeByScore return a range of members in a sorted set, by score
	ZRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) ([]string, error)

	// ZCard get the number of members in a sorted set
	ZCard(ctx context.Context, key string) (int64, error)

	// ZAdd add one or more members to a sorted set, or update its score if it already exists
	ZAdd(ctx context.Context, key string, members []*goredis.Z, ttl time.Duration) (int64, error)

	// GeoRadius query a sorted set representing a geospatial index to fetch members matching a given maximum distance from a point
	GeoRadius(ctx context.Context, key string, lat, lng float64, query *goredis.GeoRadiusQuery) ([]goredis.GeoLocation, error)

	// GeoAdd add one or more geospatial items in the geospatial index represented using a sorted set
	GeoAdd(ctx context.Context, key string, locations ...*goredis.GeoLocation) (int64, error)

	// SMembers get all the members in a set
	SMembers(ctx context.Context, key string) ([]string, error)

	// SCard get the number of members in a set
	SCard(ctx context.Context, key string) (int64, error)

	// SIsMember determine if a given value is a member of a set
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)

	// SRem remove one or more members from a set
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)

	// SAdd add one or more members to a set
	SAdd(ctx context.Context, key string, members []interface{}, ttl time.Duration) (int64, error)

	// Eval execute a Lua script server side
	Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error)

	// EvalSha execute a Lua script server side
	EvalSha(ctx context.Context, sha1 string, keys []string, args []interface{}) (interface{}, error)
}

type redis struct {
	logger   zlog.Logger
	instance Redis
	o        func(opt string) string
}

func (r *redis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) (string, error) {
	key = r.o(key)
	return r.instance.Set(ctx, key, value, ttl)
}

func (r *redis) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	key = r.o(key)
	return r.instance.SetNX(ctx, key, value, ttl)
}

func (r *redis) Get(ctx context.Context, key string) (string, error) {
	key = r.o(key)
	return r.instance.Get(ctx, key)
}

func (r *redis) Incr(ctx context.Context, key string) (int64, error) {
	key = r.o(key)
	return r.instance.Incr(ctx, key)
}

func (r *redis) IncrBy(ctx context.Context, key string, inc int64) (int64, error) {
	key = r.o(key)
	return r.instance.IncrBy(ctx, key, inc)
}

func (r *redis) StrLen(ctx context.Context, key string) (int64, error) {
	key = r.o(key)
	return r.instance.StrLen(ctx, key)
}

func (r *redis) Decr(ctx context.Context, key string) (int64, error) {
	key = r.o(key)
	return r.instance.Decr(ctx, key)
}

func (r *redis) DecrBy(ctx context.Context, key string, dec int64) (int64, error) {
	key = r.o(key)
	return r.instance.DecrBy(ctx, key, dec)
}

func (r *redis) SMembers(ctx context.Context, key string) ([]string, error) {
	key = r.o(key)
	return r.instance.SMembers(ctx, key)
}

func (r *redis) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	keys = r.appendKeysPrefix(keys)
	return r.instance.MGet(ctx, keys...)
}

func (r *redis) GeoRadius(ctx context.Context, key string, lat, lng float64, query *goredis.GeoRadiusQuery) ([]goredis.GeoLocation, error) {
	key = r.o(key)
	return r.instance.GeoRadius(ctx, key, lng, lat, query)
}

func (r *redis) GeoAdd(ctx context.Context, key string, locations ...*goredis.GeoLocation) (int64, error) {
	key = r.o(key)
	return r.instance.GeoAdd(ctx, key, locations...)
}

func (r *redis) SCard(ctx context.Context, key string) (int64, error) {
	key = r.o(key)
	return r.instance.SCard(ctx, key)
}

func (r *redis) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	key = r.o(key)
	return r.instance.SIsMember(ctx, key, member)
}

func (r *redis) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	key = r.o(key)
	return r.instance.SRem(ctx, key, members...)
}

func (r *redis) Del(ctx context.Context, keys ...string) (int64, error) {
	keys = r.appendKeysPrefix(keys)
	return r.instance.Del(ctx, keys...)
}

func (r *redis) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	key = r.o(key)
	return r.instance.HDel(ctx, key, fields...)
}

func (r *redis) HLen(ctx context.Context, key string) (int64, error) {
	key = r.o(key)
	return r.instance.HLen(ctx, key)
}

func (r *redis) HKeys(ctx context.Context, key string) ([]string, error) {
	key = r.o(key)
	return r.instance.HKeys(ctx, key)
}

func (r *redis) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	key = r.o(key)
	return r.instance.ZRem(ctx, key, members...)
}

func (r *redis) ZRemRangeByScore(ctx context.Context, key, min, max string) (int64, error) {
	key = r.o(key)
	return r.instance.ZRemRangeByScore(ctx, key, min, max)
}

func (r *redis) ZRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) ([]string, error) {
	key = r.o(key)
	return r.instance.ZRangeByScore(ctx, key, opt)
}

func (r *redis) ZCard(ctx context.Context, key string) (int64, error) {
	key = r.o(key)
	return r.instance.ZCard(ctx, key)
}

func (r *redis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	key = r.o(key)
	return r.instance.Expire(ctx, key, expiration)
}

func (r *redis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	key = r.o(key)
	return r.instance.ExpireAt(ctx, key, tm)
}

func (r *redis) Eval(ctx context.Context, script string, keys []string, args []interface{}) (interface{}, error) {
	keys = r.appendKeysPrefix(keys)
	return r.instance.Eval(ctx, script, keys, args)
}

func (r *redis) EvalSha(ctx context.Context, sha1 string, keys []string, args []interface{}) (interface{}, error) {
	keys = r.appendKeysPrefix(keys)
	return r.instance.EvalSha(ctx, sha1, keys, args)
}

func (r *redis) SAdd(ctx context.Context, key string, members []interface{}, ttl time.Duration) (int64, error) {
	key = r.o(key)
	return r.instance.SAdd(ctx, key, members, ttl)
}

func (r *redis) MSet(ctx context.Context, values []interface{}, ttl time.Duration) (string, error) {
	var preparedValues []interface{}
	for i := 0; i < len(values); i += 2 {
		key := r.o(fmt.Sprintf("%v", values[i]))
		value := values[i+1]
		preparedValues = append(preparedValues, key, value)
	}
	return r.instance.MSet(ctx, preparedValues, ttl)
}

func (r *redis) HSet(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (int64, error) {
	key = r.o(key)
	return r.instance.HSet(ctx, r.o(key), field, value, ttl)
}

func (r *redis) HSetNX(ctx context.Context, key, field string, value interface{}, ttl time.Duration) (bool, error) {
	key = r.o(key)
	return r.instance.HSetNX(ctx, key, field, value, ttl)
}

func (r *redis) HGet(ctx context.Context, key, field string) (string, error) {
	key = r.o(key)
	return r.instance.HGet(ctx, key, field)
}

func (r *redis) HMSet(ctx context.Context, key string, values []interface{}, ttl time.Duration) (bool, error) {
	key = r.o(key)
	return r.instance.HMSet(ctx, key, values, ttl)
}

func (r *redis) ZAdd(ctx context.Context, key string, members []*goredis.Z, ttl time.Duration) (int64, error) {
	key = r.o(key)
	return r.instance.ZAdd(ctx, r.o(key), members, ttl)
}

func (r *redis) Exists(ctx context.Context, keys ...string) (int64, error) {
	keys = r.appendKeysPrefix(keys)
	return r.instance.Exists(ctx, keys...)
}

func (r *redis) appendKeysPrefix(keys []string) []string {
	newKeys := make([]string, 0)
	for _, key := range keys {
		newKeys = append(newKeys, r.o(key))
	}
	return newKeys
}

func newRedis(logger zlog.Logger, instance Redis, prefix string) Redis {
	prefixBuilder := strings.Builder{}
	prefixBuilder.WriteString(prefix)
	if prefix != "" {
		prefixBuilder.WriteByte(':')
	}

	return &redis{
		logger:   logger,
		instance: instance,
		o: func(opt string) string {
			builder := strings.Builder{}
			builder.WriteString(prefixBuilder.String())
			builder.WriteString(opt)
			return builder.String()
		},
	}
}
