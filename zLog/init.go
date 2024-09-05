package zLog

import (
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConf struct {
	Compress      bool   `mapstructure:"compress"`
	ConsoleStdout bool   `mapstructure:"consoleStdout"`
	FileStdout    bool   `mapstructure:"fileStdout"`
	Level         string `mapstructure:"level"`
	LocalTime     bool   `mapstructure:"localtime"`
	Path          string `mapstructure:"path"`
	MaxSize       int    `mapstructure:"maxSize"`
	MaxAge        int    `mapstructure:"maxAge"`
	MaxBackups    int    `mapstructure:"maxBackups"`
	CollectorURL  string `mapstructure:"collectorURL"`
	Insecure      bool   `mapstructure:"insecure"`
}

func Init(cfgLog LogConf) {
	rotete := &lumberjack.Logger{
		Filename:   cfgLog.Path,
		MaxSize:    cfgLog.MaxSize,
		MaxAge:     cfgLog.MaxAge,
		MaxBackups: cfgLog.MaxBackups,
		LocalTime:  false,
		Compress:   cfgLog.Compress,
	}

	level, err := zapcore.ParseLevel(cfgLog.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	ws := []zapcore.WriteSyncer{
		zapcore.AddSync(rotete),
	}
	if cfgLog.ConsoleStdout {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}

	SetDefaultLogger(
		MustNew(
			OptLevel(level),
			OptWriteSyncers(ws...),
		),
	)
}
