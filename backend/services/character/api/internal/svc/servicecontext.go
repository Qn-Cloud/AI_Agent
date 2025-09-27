package svc

import (
	common "ai-roleplay/common/utils"
	"ai-roleplay/services/character/api/internal/config"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Db:     common.GetDB(c.Mysql),
		Redis:  common.GetRedis(c.Redis),
	}
}
