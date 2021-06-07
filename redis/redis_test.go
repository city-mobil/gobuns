package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
)

type localStorage struct {
	storage map[string]interface{}
}

func (localStorage) GetEx(ctx context.Context, key string, expiration time.Duration) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) GetDel(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) SetArgs(ctx context.Context, key string, value interface{}, a goredis.SetArgs) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *goredis.ScanCmd {
	panic("implement me")
}

func (localStorage) HRandField(ctx context.Context, key string, count int, withValues bool) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) LPopCount(ctx context.Context, key string, count int) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) LPos(ctx context.Context, key string, value string, args goredis.LPosArgs) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LPosCount(ctx context.Context, key string, value string, count int64, args goredis.LPosArgs) *goredis.IntSliceCmd {
	panic("implement me")
}

func (localStorage) LMove(ctx context.Context, source, destination, srcpos, destpos string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) SMIsMember(ctx context.Context, key string, members ...interface{}) *goredis.BoolSliceCmd {
	panic("implement me")
}

func (localStorage) XInfoStream(ctx context.Context, key string) *goredis.XInfoStreamCmd {
	panic("implement me")
}

func (localStorage) XInfoConsumers(ctx context.Context, key string, group string) *goredis.XInfoConsumersCmd {
	panic("implement me")
}

func (localStorage) ZMScore(ctx context.Context, key string, members ...string) *goredis.FloatSliceCmd {
	panic("implement me")
}

func (localStorage) ZRandMember(ctx context.Context, key string, count int, withScores bool) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZDiff(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZDiffWithScores(ctx context.Context, keys ...string) *goredis.ZSliceCmd {
	panic("implement me")
}

func (localStorage) Pipeline() goredis.Pipeliner {
	panic("implement me")
}

func (localStorage) Pipelined(ctx context.Context, fn func(goredis.Pipeliner) error) ([]goredis.Cmder, error) {
	panic("implement me")
}

func (localStorage) TxPipelined(ctx context.Context, fn func(goredis.Pipeliner) error) ([]goredis.Cmder, error) {
	panic("implement me")
}

func (localStorage) TxPipeline() goredis.Pipeliner {
	panic("implement me")
}

func (localStorage) Command(ctx context.Context) *goredis.CommandsInfoCmd {
	panic("implement me")
}

func (localStorage) ClientGetName(ctx context.Context) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) Echo(ctx context.Context, message interface{}) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) Ping(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Quit(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Del(ctx context.Context, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) Unlink(ctx context.Context, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) Dump(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) Exists(ctx context.Context, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) Expire(ctx context.Context, key string, expiration time.Duration) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) ExpireAt(ctx context.Context, key string, tm time.Time) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) Keys(ctx context.Context, pattern string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) Migrate(ctx context.Context, host, port, key string, db int, timeout time.Duration) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Move(ctx context.Context, key string, db int) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) ObjectRefCount(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ObjectEncoding(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) ObjectIdleTime(ctx context.Context, key string) *goredis.DurationCmd {
	panic("implement me")
}

func (localStorage) Persist(ctx context.Context, key string) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) PExpire(ctx context.Context, key string, expiration time.Duration) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) PExpireAt(ctx context.Context, key string, tm time.Time) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) PTTL(ctx context.Context, key string) *goredis.DurationCmd {
	panic("implement me")
}

func (localStorage) RandomKey(ctx context.Context) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) Rename(ctx context.Context, key, newkey string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) RenameNX(ctx context.Context, key, newkey string) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) Restore(ctx context.Context, key string, ttl time.Duration, value string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Sort(ctx context.Context, key string, sort *goredis.Sort) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) SortStore(ctx context.Context, key, store string, sort *goredis.Sort) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SortInterfaces(ctx context.Context, key string, sort *goredis.Sort) *goredis.SliceCmd {
	panic("implement me")
}

func (localStorage) Touch(ctx context.Context, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) TTL(ctx context.Context, key string) *goredis.DurationCmd {
	panic("implement me")
}

func (localStorage) Type(ctx context.Context, key string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Append(ctx context.Context, key, value string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) Decr(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) DecrBy(ctx context.Context, key string, decrement int64) *goredis.IntCmd {
	panic("implement me")
}

func (l *localStorage) Get(ctx context.Context, key string) *goredis.StringCmd {
	if v, ok := l.storage[key]; ok {
		return goredis.NewStringResult(fmt.Sprintf("%v", v), nil)
	}
	return goredis.NewStringResult("", errors.New("no data"))
}

func (localStorage) GetRange(ctx context.Context, key string, start, end int64) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) GetSet(ctx context.Context, key string, value interface{}) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) Incr(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) IncrBy(ctx context.Context, key string, value int64) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) IncrByFloat(ctx context.Context, key string, value float64) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) MGet(ctx context.Context, keys ...string) *goredis.SliceCmd {
	panic("implement me")
}

func (localStorage) MSet(ctx context.Context, values ...interface{}) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) MSetNX(ctx context.Context, values ...interface{}) *goredis.BoolCmd {
	panic("implement me")
}

func (l *localStorage) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.StatusCmd {
	l.storage[key] = value
	return goredis.NewStatusResult("ok", nil)
}

func (l *localStorage) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd {
	l.Set(ctx, key, value, expiration)
	return goredis.NewBoolResult(true, nil)
}

func (l *localStorage) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *goredis.BoolCmd {
	return l.SetNX(ctx, key, value, expiration)
}

func (localStorage) SetRange(ctx context.Context, key string, offset int64, value string) *goredis.IntCmd {
	panic("implement me")
}

func (l *localStorage) StrLen(ctx context.Context, key string) *goredis.IntCmd {
	if v, ok := l.storage[key]; ok {
		return goredis.NewIntResult(int64(len(fmt.Sprintf("%v", v))), nil)
	}
	return goredis.NewIntResult(0, errors.New("no data"))
}

func (localStorage) GetBit(ctx context.Context, key string, offset int64) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SetBit(ctx context.Context, key string, offset int64, value int) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) BitCount(ctx context.Context, key string, bitCount *goredis.BitCount) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) BitOpAnd(ctx context.Context, destKey string, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) BitOpOr(ctx context.Context, destKey string, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) BitOpXor(ctx context.Context, destKey string, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) BitOpNot(ctx context.Context, destKey string, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) BitField(ctx context.Context, key string, args ...interface{}) *goredis.IntSliceCmd {
	panic("implement me")
}

func (localStorage) Scan(ctx context.Context, cursor uint64, match string, count int64) *goredis.ScanCmd {
	panic("implement me")
}

func (localStorage) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *goredis.ScanCmd {
	panic("implement me")
}

func (localStorage) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *goredis.ScanCmd {
	panic("implement me")
}

func (localStorage) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *goredis.ScanCmd {
	panic("implement me")
}

func (localStorage) HDel(ctx context.Context, key string, fields ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) HExists(ctx context.Context, key, field string) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) HGet(ctx context.Context, key, field string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) HGetAll(ctx context.Context, key string) *goredis.StringStringMapCmd {
	panic("implement me")
}

func (localStorage) HIncrBy(ctx context.Context, key, field string, incr int64) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) HIncrByFloat(ctx context.Context, key, field string, incr float64) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) HKeys(ctx context.Context, key string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) HLen(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) HMGet(ctx context.Context, key string, fields ...string) *goredis.SliceCmd {
	panic("implement me")
}

func (localStorage) HSet(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) HMSet(ctx context.Context, key string, values ...interface{}) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) HSetNX(ctx context.Context, key, field string, value interface{}) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) HVals(ctx context.Context, key string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) LIndex(ctx context.Context, key string, index int64) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LInsertBefore(ctx context.Context, key string, pivot, value interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LInsertAfter(ctx context.Context, key string, pivot, value interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LLen(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LPop(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) LPush(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LPushX(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LRange(ctx context.Context, key string, start, stop int64) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) LRem(ctx context.Context, key string, count int64, value interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) LSet(ctx context.Context, key string, index int64, value interface{}) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) LTrim(ctx context.Context, key string, start, stop int64) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) RPop(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) RPopLPush(ctx context.Context, source, destination string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) RPush(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) RPushX(ctx context.Context, key string, values ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SAdd(ctx context.Context, key string, members ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SCard(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SDiff(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) SDiffStore(ctx context.Context, destination string, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SInter(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) SInterStore(ctx context.Context, destination string, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SIsMember(ctx context.Context, key string, member interface{}) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) SMembers(ctx context.Context, key string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) SMembersMap(ctx context.Context, key string) *goredis.StringStructMapCmd {
	panic("implement me")
}

func (localStorage) SMove(ctx context.Context, source, destination string, member interface{}) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) SPop(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) SPopN(ctx context.Context, key string, count int64) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) SRandMember(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) SRandMemberN(ctx context.Context, key string, count int64) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) SRem(ctx context.Context, key string, members ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) SUnion(ctx context.Context, keys ...string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) SUnionStore(ctx context.Context, destination string, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XAdd(ctx context.Context, a *goredis.XAddArgs) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) XDel(ctx context.Context, stream string, ids ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XLen(ctx context.Context, stream string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XRange(ctx context.Context, stream, start, stop string) *goredis.XMessageSliceCmd {
	panic("implement me")
}

func (localStorage) XRangeN(ctx context.Context, stream, start, stop string, count int64) *goredis.XMessageSliceCmd {
	panic("implement me")
}

func (localStorage) XRevRange(ctx context.Context, stream string, start, stop string) *goredis.XMessageSliceCmd {
	panic("implement me")
}

func (localStorage) XRevRangeN(ctx context.Context, stream string, start, stop string, count int64) *goredis.XMessageSliceCmd {
	panic("implement me")
}

func (localStorage) XRead(ctx context.Context, a *goredis.XReadArgs) *goredis.XStreamSliceCmd {
	panic("implement me")
}

func (localStorage) XReadStreams(ctx context.Context, streams ...string) *goredis.XStreamSliceCmd {
	panic("implement me")
}

func (localStorage) XGroupCreate(ctx context.Context, stream, group, start string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) XGroupCreateMkStream(ctx context.Context, stream, group, start string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) XGroupSetID(ctx context.Context, stream, group, start string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) XGroupDestroy(ctx context.Context, stream, group string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XGroupDelConsumer(ctx context.Context, stream, group, consumer string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XReadGroup(ctx context.Context, a *goredis.XReadGroupArgs) *goredis.XStreamSliceCmd {
	panic("implement me")
}

func (localStorage) XAck(ctx context.Context, stream, group string, ids ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XPending(ctx context.Context, stream, group string) *goredis.XPendingCmd {
	panic("implement me")
}

func (localStorage) XPendingExt(ctx context.Context, a *goredis.XPendingExtArgs) *goredis.XPendingExtCmd {
	panic("implement me")
}

func (localStorage) XClaim(ctx context.Context, a *goredis.XClaimArgs) *goredis.XMessageSliceCmd {
	panic("implement me")
}

func (localStorage) XClaimJustID(ctx context.Context, a *goredis.XClaimArgs) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) XTrim(ctx context.Context, key string, maxLen int64) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XTrimApprox(ctx context.Context, key string, maxLen int64) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) XInfoGroups(ctx context.Context, key string) *goredis.XInfoGroupsCmd {
	panic("implement me")
}

func (localStorage) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *goredis.ZWithKeyCmd {
	panic("implement me")
}

func (localStorage) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *goredis.ZWithKeyCmd {
	panic("implement me")
}

func (localStorage) ZAdd(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZAddNX(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZAddXX(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZAddCh(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZAddNXCh(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZAddXXCh(ctx context.Context, key string, members ...*goredis.Z) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZIncr(ctx context.Context, key string, member *goredis.Z) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) ZIncrNX(ctx context.Context, key string, member *goredis.Z) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) ZIncrXX(ctx context.Context, key string, member *goredis.Z) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) ZCard(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZCount(ctx context.Context, key, min, max string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZLexCount(ctx context.Context, key, min, max string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZIncrBy(ctx context.Context, key string, increment float64, member string) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) ZInterStore(ctx context.Context, destination string, store *goredis.ZStore) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZPopMax(ctx context.Context, key string, count ...int64) *goredis.ZSliceCmd {
	panic("implement me")
}

func (localStorage) ZPopMin(ctx context.Context, key string, count ...int64) *goredis.ZSliceCmd {
	panic("implement me")
}

func (localStorage) ZRange(ctx context.Context, key string, start, stop int64) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *goredis.ZSliceCmd {
	panic("implement me")
}

func (localStorage) ZRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZRangeByLex(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZRangeByScoreWithScores(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.ZSliceCmd {
	panic("implement me")
}

func (localStorage) ZRank(ctx context.Context, key, member string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZRem(ctx context.Context, key string, members ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZRemRangeByScore(ctx context.Context, key, min, max string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZRemRangeByLex(ctx context.Context, key, min, max string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZRevRange(ctx context.Context, key string, start, stop int64) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *goredis.ZSliceCmd {
	panic("implement me")
}

func (localStorage) ZRevRangeByScore(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZRevRangeByLex(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *goredis.ZRangeBy) *goredis.ZSliceCmd {
	panic("implement me")
}

func (localStorage) ZRevRank(ctx context.Context, key, member string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ZScore(ctx context.Context, key, member string) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) ZUnionStore(ctx context.Context, dest string, store *goredis.ZStore) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) PFAdd(ctx context.Context, key string, els ...interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) PFCount(ctx context.Context, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) PFMerge(ctx context.Context, dest string, keys ...string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) BgRewriteAOF(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) BgSave(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClientKill(ctx context.Context, ipPort string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClientKillByFilter(ctx context.Context, keys ...string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ClientList(ctx context.Context) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) ClientPause(ctx context.Context, dur time.Duration) *goredis.BoolCmd {
	panic("implement me")
}

func (localStorage) ClientID(ctx context.Context) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ConfigGet(ctx context.Context, parameter string) *goredis.SliceCmd {
	panic("implement me")
}

func (localStorage) ConfigResetStat(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ConfigSet(ctx context.Context, parameter, value string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ConfigRewrite(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) DBSize(ctx context.Context) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) FlushAll(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) FlushAllAsync(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) FlushDB(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) FlushDBAsync(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Info(ctx context.Context, section ...string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) LastSave(ctx context.Context) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) Save(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Shutdown(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ShutdownSave(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ShutdownNoSave(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) SlaveOf(ctx context.Context, host, port string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) Time(ctx context.Context) *goredis.TimeCmd {
	panic("implement me")
}

func (localStorage) DebugObject(ctx context.Context, key string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) ReadOnly(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ReadWrite(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) MemoryUsage(ctx context.Context, key string, samples ...int) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *goredis.Cmd {
	panic("implement me")
}

func (localStorage) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) *goredis.Cmd {
	panic("implement me")
}

func (localStorage) ScriptExists(ctx context.Context, hashes ...string) *goredis.BoolSliceCmd {
	panic("implement me")
}

func (localStorage) ScriptFlush(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ScriptKill(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ScriptLoad(ctx context.Context, script string) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) Publish(ctx context.Context, channel string, message interface{}) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) PubSubChannels(ctx context.Context, pattern string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) PubSubNumSub(ctx context.Context, channels ...string) *goredis.StringIntMapCmd {
	panic("implement me")
}

func (localStorage) PubSubNumPat(ctx context.Context) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ClusterSlots(ctx context.Context) *goredis.ClusterSlotsCmd {
	panic("implement me")
}

func (localStorage) ClusterNodes(ctx context.Context) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) ClusterMeet(ctx context.Context, host, port string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterForget(ctx context.Context, nodeID string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterReplicate(ctx context.Context, nodeID string) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterResetSoft(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterResetHard(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterInfo(ctx context.Context) *goredis.StringCmd {
	panic("implement me")
}

func (localStorage) ClusterKeySlot(ctx context.Context, key string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ClusterGetKeysInSlot(ctx context.Context, slot int, count int) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ClusterCountFailureReports(ctx context.Context, nodeID string) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ClusterCountKeysInSlot(ctx context.Context, slot int) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) ClusterDelSlots(ctx context.Context, slots ...int) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterDelSlotsRange(ctx context.Context, min, max int) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterSaveConfig(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterSlaves(ctx context.Context, nodeID string) *goredis.StringSliceCmd {
	panic("implement me")
}

func (localStorage) ClusterFailover(ctx context.Context) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterAddSlots(ctx context.Context, slots ...int) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) ClusterAddSlotsRange(ctx context.Context, min, max int) *goredis.StatusCmd {
	panic("implement me")
}

func (localStorage) GeoAdd(ctx context.Context, key string, geoLocation ...*goredis.GeoLocation) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) GeoPos(ctx context.Context, key string, members ...string) *goredis.GeoPosCmd {
	panic("implement me")
}

func (localStorage) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *goredis.GeoRadiusQuery) *goredis.GeoLocationCmd {
	panic("implement me")
}

func (localStorage) GeoRadiusStore(ctx context.Context, key string, longitude, latitude float64, query *goredis.GeoRadiusQuery) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) GeoRadiusByMember(ctx context.Context, key, member string, query *goredis.GeoRadiusQuery) *goredis.GeoLocationCmd {
	panic("implement me")
}

func (localStorage) GeoRadiusByMemberStore(ctx context.Context, key, member string, query *goredis.GeoRadiusQuery) *goredis.IntCmd {
	panic("implement me")
}

func (localStorage) GeoDist(ctx context.Context, key string, member1, member2, unit string) *goredis.FloatCmd {
	panic("implement me")
}

func (localStorage) GeoHash(ctx context.Context, key string, members ...string) *goredis.StringSliceCmd {
	panic("implement me")
}

func newLocalStorage() goredis.Cmdable {
	return &localStorage{
		storage: make(map[string]interface{}),
	}
}
