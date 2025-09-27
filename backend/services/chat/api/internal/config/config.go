package config

import (
	common "ai-roleplay/common/utils"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql common.Config
	Redis common.RedisCfg
}
