package main

import (
	"context"
	"os"
	"time"

	"github.com/city-mobil/gobuns/config"
	"github.com/city-mobil/gobuns/redis"
	"github.com/city-mobil/gobuns/redis/redisconfig"
	"github.com/city-mobil/gobuns/zlog"
)

func main() {
	redisConfigFn := redisconfig.NewClusterConfig()
	logger := zlog.New(os.Stdout)

	err := config.InitOnce()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	redisConfig, err := redisConfigFn()
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	cluster := redis.NewDefaultCluster(logger, redisConfig)
	ctx := context.Background()

	_, _ = cluster.Set(ctx, "hello", "world", time.Hour)

	result, err := cluster.Get(ctx, "hello")
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	logger.Info().Msgf("Cached value: %s", result)
}
