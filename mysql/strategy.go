package mysql

import "hash/crc32"

type StrategyCallback func(shards []Shard) Shard

func StrategyCRC32(key string, shards []Shard) Shard {
	sum := crc32.ChecksumIEEE([]byte(key))

	l := uint32(len(shards))
	return shards[sum%l]
}
