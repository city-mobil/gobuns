package redis

import (
	"context"
	"errors"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

var (
	ErrorRedisUnavailable = errors.New("redis unavailable")
)

type fallback struct{}

func newFallback() goredis.Cmdable {
	return &fallback{}
}

func (fallback) GetEx(ctx context.Context, key string, expiration time.Duration) *goredis.StringCmd {
	return nil
}

func (fallback) GetDel(ctx context.Context, key string) *goredis.StringCmd {
	return nil
}

func (fallback) SetArgs(ctx context.Context, key string, value interface{}, a goredis.SetArgs) *goredis.StatusCmd {
	return nil
}

func (fallback) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.StatusCmd {
	return nil
}

func (fallback) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *goredis.ScanCmd {
	return nil
}

func (fallback) HRandField(ctx context.Context, key string, count int, withValues bool) *goredis.StringSliceCmd {
	return nil
}

func (fallback) LPopCount(ctx context.Context, key string, count int) *goredis.StringSliceCmd {
	return nil
}

func (fallback) LPos(ctx context.Context, key, value string, args goredis.LPosArgs) *goredis.IntCmd {
	return nil
}

func (fallback) LPosCount(ctx context.Context, key, value string, count int64, args goredis.LPosArgs) *goredis.IntSliceCmd {
	return nil
}

func (fallback) LMove(ctx context.Context, source, destination, srcpos, destpos string) *goredis.StringCmd {
	return nil
}

func (fallback) SMIsMember(ctx context.Context, key string, members ...interface{}) *goredis.BoolSliceCmd {
	return nil
}

func (fallback) XInfoStream(ctx context.Context, key string) *goredis.XInfoStreamCmd {
	return nil
}

func (fallback) XInfoConsumers(ctx context.Context, key, group string) *goredis.XInfoConsumersCmd {
	return nil
}

func (fallback) ZMScore(ctx context.Context, key string, members ...string) *goredis.FloatSliceCmd {
	return nil
}

func (fallback) ZRandMember(ctx context.Context, key string, count int, withScores bool) *goredis.StringSliceCmd {
	return nil
}

func (fallback) ZDiff(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	return nil
}

func (fallback) ZDiffWithScores(ctx context.Context, keys ...string) *goredis.ZSliceCmd {
	return nil
}

func (fallback) Pipeline() goredis.Pipeliner {
	return nil
}

func (fallback) Pipelined(ctx context.Context, fn func(goredis.Pipeliner) error) ([]goredis.Cmder, error) {
	return nil, ErrorRedisUnavailable
}

func (fallback) TxPipelined(ctx context.Context, fn func(goredis.Pipeliner) error) ([]goredis.Cmder, error) {
	return nil, ErrorRedisUnavailable
}

func (fallback) TxPipeline() goredis.Pipeliner {
	return nil
}

func (fallback) Command(ctx context.Context) *goredis.CommandsInfoCmd {
	return nil
}

func (fallback) ClientGetName(ctx context.Context) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) Echo(ctx context.Context, message interface{}) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) Ping(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Quit(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Del(ctx context.Context, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) Unlink(ctx context.Context, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) Dump(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) Exists(ctx context.Context, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) Expire(ctx context.Context, key string, expiration time.Duration) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) ExpireAt(ctx context.Context, key string, tm time.Time) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) Keys(ctx context.Context, pattern string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Move(ctx context.Context, key string, db int) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) ObjectRefCount(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ObjectEncoding(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) ObjectIdleTime(ctx context.Context, key string) *goredis.DurationCmd {
	return goredis.NewDurationResult(time.Duration(0), ErrorRedisUnavailable)
}

func (fallback) Persist(ctx context.Context, key string) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) PExpire(ctx context.Context, key string, expiration time.Duration) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) PExpireAt(ctx context.Context, key string, tm time.Time) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) PTTL(ctx context.Context, key string) *goredis.DurationCmd {
	return goredis.NewDurationResult(time.Duration(0), ErrorRedisUnavailable)
}

func (fallback) RandomKey(ctx context.Context) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) Rename(ctx context.Context, key, newkey string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) RenameNX(ctx context.Context, key, newkey string) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) Restore(ctx context.Context, key string, ttl time.Duration, value string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Sort(ctx context.Context, key string, sort *goredis.Sort) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) SortStore(ctx context.Context, key, store string, sort *goredis.Sort) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SortInterfaces(ctx context.Context, key string, sort *goredis.Sort) *goredis.SliceCmd {
	return goredis.NewSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) Touch(ctx context.Context, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) TTL(ctx context.Context, key string) *goredis.DurationCmd {
	return goredis.NewDurationResult(time.Duration(0), ErrorRedisUnavailable)
}

func (fallback) Type(ctx context.Context, key string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Scan(ctx context.Context, cursor uint64, match string, count int64) *goredis.ScanCmd {
	return nil
}

func (fallback) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *goredis.ScanCmd {
	return nil
}

func (fallback) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *goredis.ScanCmd {
	return nil
}

func (fallback) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *goredis.ScanCmd {
	return nil
}

func (fallback) Append(ctx context.Context, key, value string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) BitCount(ctx context.Context, key string, bitCount *goredis.BitCount) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) BitOpAnd(ctx context.Context, destKey string, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) BitOpOr(ctx context.Context, destKey string, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) BitOpXor(ctx context.Context, destKey string, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) BitOpNot(ctx context.Context, destKey, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) BitField(ctx context.Context, key string, args ...interface{}) *goredis.IntSliceCmd {
	return nil
}

func (fallback) Decr(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) DecrBy(ctx context.Context, key string, decrement int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) Get(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) GetBit(ctx context.Context, key string, offset int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) GetRange(ctx context.Context, key string, start, end int64) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) GetSet(ctx context.Context, key string, value interface{}) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) Incr(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) IncrBy(ctx context.Context, key string, value int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) IncrByFloat(ctx context.Context, key string, value float64) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) MGet(ctx context.Context, keys ...string) *goredis.SliceCmd {
	return goredis.NewSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) MSet(ctx context.Context, values ...interface{}) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) MSetNX(ctx context.Context, values ...interface{}) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) SetBit(ctx context.Context, key string, offset int64, value int) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) SetRange(ctx context.Context, key string, offset int64, value string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) StrLen(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) HDel(ctx context.Context, key string, fields ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) HExists(ctx context.Context, key, field string) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) HGet(ctx context.Context, key, field string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) HGetAll(ctx context.Context, key string) *goredis.StringStringMapCmd {
	return goredis.NewStringStringMapResult(nil, ErrorRedisUnavailable)
}

func (fallback) HIncrBy(ctx context.Context, key, field string, incr int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) HIncrByFloat(ctx context.Context, key, field string, incr float64) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) HKeys(ctx context.Context, key string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) HLen(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) HMGet(ctx context.Context, key string, fields ...string) *goredis.SliceCmd {
	return goredis.NewSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) HSet(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) HMSet(ctx context.Context, key string, values ...interface{}) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) HSetNX(ctx context.Context, key, field string, value interface{}) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) HVals(ctx context.Context, key string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) LIndex(ctx context.Context, key string, index int64) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) LLen(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) LPop(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) LPush(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) LPushX(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) LRange(ctx context.Context, key string, start, stop int64) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) LRem(ctx context.Context, key string, count int64, value interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) LSet(ctx context.Context, key string, index int64, value interface{}) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) LTrim(ctx context.Context, key string, start, stop int64) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) RPop(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) RPopLPush(ctx context.Context, source, destination string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) RPush(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) RPushX(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SAdd(ctx context.Context, key string, members ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SCard(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SDiff(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) SDiffStore(ctx context.Context, destination string, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SInter(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) SInterStore(ctx context.Context, destination string, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SIsMember(ctx context.Context, key string, member interface{}) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) SMembers(ctx context.Context, key string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) SMembersMap(ctx context.Context, key string) *goredis.StringStructMapCmd {
	return nil
}

func (fallback) SMove(ctx context.Context, source, destination string, member interface{}) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) SPop(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) SPopN(ctx context.Context, key string, count int64) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) SRandMember(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) SRandMemberN(ctx context.Context, key string, count int64) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) SRem(ctx context.Context, key string, members ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) SUnion(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) SUnionStore(ctx context.Context, destination string, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XAdd(ctx context.Context, a *goredis.XAddArgs) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) XDel(ctx context.Context, stream string, ids ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XLen(ctx context.Context, stream string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XRange(ctx context.Context, stream, start, stop string) *goredis.XMessageSliceCmd {
	return goredis.NewXMessageSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) XRangeN(ctx context.Context, stream, start, stop string, count int64) *goredis.XMessageSliceCmd {
	return goredis.NewXMessageSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) XRevRange(ctx context.Context, stream, start, stop string) *goredis.XMessageSliceCmd {
	return goredis.NewXMessageSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) XRevRangeN(ctx context.Context, stream, start, stop string, count int64) *goredis.XMessageSliceCmd {
	return goredis.NewXMessageSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) XRead(ctx context.Context, a *goredis.XReadArgs) *goredis.XStreamSliceCmd {
	return goredis.NewXStreamSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) XReadStreams(ctx context.Context, streams ...string) *goredis.XStreamSliceCmd {
	return goredis.NewXStreamSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) XGroupCreate(ctx context.Context, stream, group, start string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) XGroupSetID(ctx context.Context, stream, group, start string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) XGroupDestroy(ctx context.Context, stream, group string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XReadGroup(ctx context.Context, a *goredis.XReadGroupArgs) *goredis.XStreamSliceCmd {
	return goredis.NewXStreamSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) XAck(ctx context.Context, stream, group string, ids ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XPending(ctx context.Context, stream, group string) *goredis.XPendingCmd {
	return nil
}

func (fallback) XPendingExt(ctx context.Context, a *goredis.XPendingExtArgs) *goredis.XPendingExtCmd {
	return nil
}

func (fallback) XClaim(ctx context.Context, a *goredis.XClaimArgs) *goredis.XMessageSliceCmd {
	return nil
}

func (fallback) XClaimJustID(ctx context.Context, a *goredis.XClaimArgs) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) XTrim(ctx context.Context, key string, maxLen int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XTrimApprox(ctx context.Context, key string, maxLen int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) XInfoGroups(ctx context.Context, key string) *goredis.XInfoGroupsCmd {
	return nil
}

func (fallback) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *goredis.ZWithKeyCmd {
	return nil
}

func (fallback) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *goredis.ZWithKeyCmd {
	return nil
}

func (fallback) ZAdd(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZAddNX(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZAddXX(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZAddCh(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZAddNXCh(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZAddXXCh(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZIncr(ctx context.Context, key string, member *goredis.Z) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) ZIncrNX(ctx context.Context, key string, member *goredis.Z) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) ZIncrXX(ctx context.Context, key string, member *goredis.Z) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) ZCard(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZCount(ctx context.Context, key, min, max string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZLexCount(ctx context.Context, key, min, max string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZIncrBy(ctx context.Context, key string, increment float64, member string) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) ZInterStore(ctx context.Context, destination string, store *goredis.ZStore) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZPopMax(ctx context.Context, key string, count ...int64) *goredis.ZSliceCmd {
	return goredis.NewZSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZPopMin(ctx context.Context, key string, count ...int64) *goredis.ZSliceCmd {
	return goredis.NewZSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRange(ctx context.Context, key string, start, stop int64) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *goredis.ZSliceCmd {
	return goredis.NewZSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRangeByLex(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRangeByScoreWithScores(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.ZSliceCmd {
	return goredis.NewZSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRank(ctx context.Context, key, member string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZRem(ctx context.Context, key string, members ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZRemRangeByScore(ctx context.Context, key, min, max string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZRemRangeByLex(ctx context.Context, key, min, max string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZRevRange(ctx context.Context, key string, start, stop int64) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *goredis.ZSliceCmd {
	return goredis.NewZSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRevRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRevRangeByLex(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.ZSliceCmd {
	return goredis.NewZSliceCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ZRevRank(ctx context.Context, key, member string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ZScore(ctx context.Context, key, member string) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) ZUnionStore(ctx context.Context, dest string, store *goredis.ZStore) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) PFAdd(ctx context.Context, key string, els ...interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) PFCount(ctx context.Context, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) PFMerge(ctx context.Context, dest string, keys ...string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) BgRewriteAOF(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) BgSave(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClientKill(ctx context.Context, ipPort string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClientKillByFilter(ctx context.Context, keys ...string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ClientList(ctx context.Context) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) ClientPause(ctx context.Context, dur time.Duration) *goredis.BoolCmd {
	return goredis.NewBoolResult(false, ErrorRedisUnavailable)
}

func (fallback) ClientID(ctx context.Context) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ConfigGet(ctx context.Context, parameter string) *goredis.SliceCmd {
	return goredis.NewSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ConfigResetStat(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ConfigSet(ctx context.Context, parameter, value string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ConfigRewrite(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) DBSize(ctx context.Context) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) FlushAll(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) FlushAllAsync(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) FlushDB(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) FlushDBAsync(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Info(ctx context.Context, section ...string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) LastSave(ctx context.Context) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) Save(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Shutdown(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ShutdownSave(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ShutdownNoSave(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) SlaveOf(ctx context.Context, host, port string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) Time(ctx context.Context) *goredis.TimeCmd {
	return goredis.NewTimeCmdResult(time.Now(), ErrorRedisUnavailable)
}

func (fallback) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *goredis.Cmd {
	return goredis.NewCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *goredis.Cmd {
	return goredis.NewCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ScriptExists(ctx context.Context, hashes ...string) *goredis.BoolSliceCmd {
	return goredis.NewBoolSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ScriptFlush(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ScriptKill(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ScriptLoad(ctx context.Context, script string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) DebugObject(ctx context.Context, key string) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) Publish(ctx context.Context, channel string, message interface{}) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) PubSubChannels(ctx context.Context, pattern string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) PubSubNumSub(ctx context.Context, channels ...string) *goredis.StringIntMapCmd {
	return goredis.NewStringIntMapCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) PubSubNumPat(ctx context.Context) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ClusterSlots(ctx context.Context) *goredis.ClusterSlotsCmd {
	return goredis.NewClusterSlotsCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) ClusterNodes(ctx context.Context) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterMeet(ctx context.Context, host, port string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterForget(ctx context.Context, nodeID string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterReplicate(ctx context.Context, nodeID string) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterResetSoft(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterResetHard(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterInfo(ctx context.Context) *goredis.StringCmd {
	return goredis.NewStringResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterKeySlot(ctx context.Context, key string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ClusterGetKeysInSlot(ctx context.Context, slot, count int) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ClusterCountFailureReports(ctx context.Context, nodeID string) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ClusterCountKeysInSlot(ctx context.Context, slot int) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) ClusterDelSlots(ctx context.Context, slots ...int) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterDelSlotsRange(ctx context.Context, min, max int) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterSaveConfig(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterSlaves(ctx context.Context, nodeID string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ClusterFailover(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterAddSlots(ctx context.Context, slots ...int) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ClusterAddSlotsRange(ctx context.Context, min, max int) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) GeoAdd(ctx context.Context, key string, geoLocation ...*goredis.GeoLocation) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) GeoPos(ctx context.Context, key string, members ...string) *goredis.GeoPosCmd {
	return goredis.NewGeoPosCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *goredis.GeoRadiusQuery) *goredis.GeoLocationCmd {
	return goredis.NewGeoLocationCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *goredis.GeoRadiusQuery) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) GeoRadiusByMember(ctx context.Context, key, member string, query *goredis.GeoRadiusQuery) *goredis.GeoLocationCmd {
	return goredis.NewGeoLocationCmdResult(nil, ErrorRedisUnavailable)
}

func (fallback) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *goredis.GeoRadiusQuery) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}

func (fallback) GeoDist(ctx context.Context, key, member1, member2, unit string) *goredis.FloatCmd {
	return goredis.NewFloatResult(0, ErrorRedisUnavailable)
}

func (fallback) GeoHash(ctx context.Context, key string, members ...string) *goredis.StringSliceCmd {
	return goredis.NewStringSliceResult(nil, ErrorRedisUnavailable)
}

func (fallback) ReadOnly(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) ReadWrite(ctx context.Context) *goredis.StatusCmd {
	return goredis.NewStatusResult("", ErrorRedisUnavailable)
}

func (fallback) MemoryUsage(ctx context.Context, key string, samples ...int) *goredis.IntCmd {
	return goredis.NewIntResult(0, ErrorRedisUnavailable)
}
