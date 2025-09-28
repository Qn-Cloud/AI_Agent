package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 语音服务配置
	Speech SpeechConfig

	// 第三方服务配置
	External ExternalConfig
}

// 语音服务配置
type SpeechConfig struct {
	ASR     ASRConfig
	TTS     TTSConfig
	Storage StorageConfig
}

// ASR配置
type ASRConfig struct {
	Provider         string
	Timeout          time.Duration
	MaxDuration      time.Duration
	SupportedFormats []string
}

// TTS配置
type TTSConfig struct {
	Provider         string
	Timeout          time.Duration
	MaxTextLength    int
	DefaultFormat    string
	SupportedFormats []string
}

// 存储配置
type StorageConfig struct {
	Type      string
	BaseURL   string
	LocalPath string
}

// 第三方服务配置
type ExternalConfig struct {
	Baidu   BaiduConfig
	Tencent TencentConfig
	Aliyun  AliyunConfig
}

// 百度AI配置
type BaiduConfig struct {
	AppID     string
	APIKey    string
	SecretKey string
}

// 腾讯云配置
type TencentConfig struct {
	SecretID  string
	SecretKey string
	Region    string
}

// 阿里云配置
type AliyunConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string
}
