package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

const (
	RequestLogName = "request"
	bizDefaultLogName = "default"
)

type LogConf struct {
	ServerName   string        `mapstructure:"serverName"`
	Path         string        `mapstructure:"path"`
	BizMaxAge    time.Duration `mapstructure:"bizMaxAge"`
	LowLevel     string        `mapstructure:"lowLevel"`
	CallMaxAge   time.Duration `mapstructure:"callMaxAge"`
}

type ZapConf struct {
	ServerName   string        // 服务名
	Path         string        // 日志目录地址
	Name         string        // 日志名字
	JsonFormat   bool          // json格式日志
	LogInConsole bool          // 打印到控制台
	RotationTime time.Duration // 日志分割时间
	MaxAge       time.Duration // 日志保留时间，单位:小时
	LowLevel     zapcore.Level // 记录的最小级别
	HighLevel    zapcore.Level // 记录的最高级别
}

func (conf *ZapConf) CheckConf() {
	if conf.Name == "" {
		conf.Name = "default"
	}

	if conf.LowLevel > conf.HighLevel {
		conf.HighLevel = conf.LowLevel
	}

	// 限制只有开发环境日志会打印到控制台
	if gin.Mode() == gin.DebugMode {
		conf.LogInConsole = true
	} else {
		conf.LogInConsole = false
	}
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.DebugLevel
	}
}
