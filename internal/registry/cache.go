package registry

import (
	"github.com/mandarine-io/Backend/internal/config"
	"github.com/mandarine-io/Backend/pkg/storage/cache/memory"
	redis2 "github.com/mandarine-io/Backend/pkg/storage/cache/redis"
	"github.com/rs/zerolog/log"
	"time"
)

func setupCacheManager(c *Container) {
	log.Debug().Msg("setup cache manager")
	switch c.Config.Cache.Type {
	case config.MemoryCacheType:
		c.Cache.Manager = memory.NewManager(time.Duration(c.Config.Cache.TTL) * time.Second)
	case config.RedisCacheType:
		if c.Config.Cache.Redis == nil {
			log.Fatal().Msg("redis config is nil")
		}
		redisConfig := mapAppCacheRedisConfigToRedisConfig(&c.Config.Cache)
		c.Cache.RDB = redis2.MustNewClient(redisConfig)
		c.Cache.Manager = redis2.NewManager(c.Cache.RDB, time.Duration(c.Config.Cache.TTL)*time.Second)
	default:
		log.Fatal().Msgf("unknown cache type: %s", c.Config.Cache.Type)
	}
}

func mapAppCacheRedisConfigToRedisConfig(cfg *config.CacheConfig) *redis2.Config {
	return &redis2.Config{
		Address:  cfg.Redis.Address,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DBIndex:  cfg.Redis.DBIndex,
	}
}
