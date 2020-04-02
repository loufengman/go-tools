package log

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	path   string `json: "path"`   //log日志目录
	prefix string `json: "prefix"` //log日志文件名前缀
	maxAge int64    `json: "maxAge"` //log保存最大天数
}

var (
	Logger      *zap.Logger
	SugarLogger *zap.SugaredLogger
	isDebug     bool
)

var baseLogConfig = LogConfig{
	path:   "./log/",
	prefix: "test",
	maxAge: 7,
}

func SetDebugMod() {
	isDebug = true
}

func initConfig(conf LogConfig) {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeLevel:  zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.WarnLevel
	})

	lowLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(conf, "_log.")
	warnWriter := getWriter(conf, "_error_log.")
	debugInfoWriter := zapcore.Lock(os.Stdout)
	debugErrorWriter := zapcore.Lock(os.Stdout)
	var core zapcore.Core
	// 最后创建具体的Logger
	if isDebug {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), errorLevel),
			zapcore.NewCore(encoder, debugInfoWriter, lowLevel),
			zapcore.NewCore(encoder, debugErrorWriter, errorLevel),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), errorLevel),
		)
	}
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	SugarLogger = Logger.Sugar()
}

func getWriter(conf LogConfig, filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	if conf.path == "" {
		conf.path = baseLogConfig.path
	}
	if conf.prefix == "" {
		conf.prefix = baseLogConfig.prefix
	}

	if conf.maxAge == 0 {
		conf.maxAge = baseLogConfig.maxAge
	}
	maxAge := conf.maxAge
	hook, err := rotatelogs.New(
		conf.path+conf.prefix+filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Duration(maxAge) * time.Hour * 24),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}