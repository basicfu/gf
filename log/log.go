package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var consoleEncoder = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     "\r\n",
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.EpochTimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}
var fileEncoder = consoleEncoder
var jsonEncoder = consoleEncoder

// 日志处理建议 https://learnku.com/articles/42231
var log *zap.Logger

type Config struct {
	filename string
}

func _init(c Config) {
	jsonEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(t.UnixMilli())
	}
	consoleEncoder.LineEnding = "\r\n"
	consoleEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	}
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoder), zapcore.Lock(os.Stdout), zap.DebugLevel)
	//if config.Prod {
	//	file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	//	fileCore := zapcore.NewCore(zapcore.NewJSONEncoder(jsonEncoder), zapcore.AddSync(file), zap.DebugLevel)
	//	log = zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller(), zap.AddCallerSkip(1))
	//} else {
	log = zap.New(zapcore.NewTee(consoleCore), zap.AddCaller(), zap.AddCallerSkip(1))
	//}
}
func init() {
	_init(Config{
		filename: "log.txt",
	})
}

// 默认初始化，重新初始化可以更改参数
func Init(c Config) {
	_init(c)
}

func Debug(args ...any) {
	log.Debug(msg(args...))
}
func Info(args ...any) {
	log.Info(msg(args...))
}
func Warn(args ...any) {
	log.Warn(msg(args...))
}
func Error(args ...any) {
	log.Error(msg(args...))
}

func DebugSkip(skip int, args ...any) {
	log.WithOptions(zap.AddCallerSkip(skip)).Debug(msg(args...))
}
func InfoSkip(skip int, args ...any) {
	log.WithOptions(zap.AddCallerSkip(skip)).Info(msg(args...))
}
func WarnSkip(skip int, args ...any) {
	log.WithOptions(zap.AddCallerSkip(skip)).Warn(msg(args...))
}
func ErrorSkip(skip int, args ...any) {
	log.WithOptions(zap.AddCallerSkip(skip)).Error(msg(args...))
}
func msg(args ...any) string {
	m := fmt.Sprintln(args...)
	return m[:len(m)-1]
}
