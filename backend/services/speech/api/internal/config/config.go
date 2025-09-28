package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	// 语音服务配置
	Speech SpeechConfig `yaml:"Speech"`

	// 第三方服务配置
	External ExternalConfig `yaml:"External"`
}

// 语音服务配置
type SpeechConfig struct {
	ASR     ASRConfig     `yaml:"ASR"`
	TTS     TTSConfig     `yaml:"TTS"`
	Storage StorageConfig `yaml:"Storage"`
}

// ASR配置
type ASRConfig struct {
	Provider         string        `yaml:"Provider"`
	Timeout          time.Duration `yaml:"Timeout"`
	MaxDuration      time.Duration `yaml:"MaxDuration"`
	SupportedFormats []string      `yaml:"SupportedFormats"`
}

// TTS配置
type TTSConfig struct {
	Provider         string        `yaml:"Provider"`
	Timeout          time.Duration `yaml:"Timeout"`
	MaxTextLength    int           `yaml:"MaxTextLength"`
	DefaultFormat    string        `yaml:"DefaultFormat"`
	SupportedFormats []string      `yaml:"SupportedFormats"`
}

// 存储配置
type StorageConfig struct {
	Type      string `yaml:"Type"`
	BaseURL   string `yaml:"BaseURL"`
	LocalPath string `yaml:"LocalPath"`
}

// 第三方服务配置
type ExternalConfig struct {
	Baidu   BaiduConfig   `yaml:"Baidu"`
	Tencent TencentConfig `yaml:"Tencent"`
	Aliyun  AliyunConfig  `yaml:"Aliyun"`
}

// 百度AI配置
type BaiduConfig struct {
	AppID     string `yaml:"AppID"`
	APIKey    string `yaml:"APIKey"`
	SecretKey string `yaml:"SecretKey"`
}

// 腾讯云配置
type TencentConfig struct {
	SecretID  string `yaml:"SecretID"`
	SecretKey string `yaml:"SecretKey"`
	Region    string `yaml:"Region"`
}

// 阿里云配置
type AliyunConfig struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	Region          string `yaml:"Region"`
}
