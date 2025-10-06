package log

import (
	"fmt"
	"os"
	"path/filepath"
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

type FileFormat string

const FileFormatConsole FileFormat = "CONSOLE"
const FileFormatJson FileFormat = "JSON"

type Config struct {
	WriteFile  bool
	FileFormat FileFormat
	Filename   string
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
	if c.WriteFile {
		file, _ := os.OpenFile(c.Filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		encoder := zapcore.NewConsoleEncoder(consoleEncoder)
		if c.FileFormat == FileFormatJson {
			encoder = zapcore.NewJSONEncoder(jsonEncoder)
		}
		fileCore := zapcore.NewCore(encoder, zapcore.AddSync(file), zap.DebugLevel)
		log = zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		log = zap.New(zapcore.NewTee(consoleCore), zap.AddCaller(), zap.AddCallerSkip(1))
	}
}
func defaultConfig() Config {
	dir, _ := os.Executable()
	return Config{
		WriteFile:  false,
		FileFormat: FileFormatConsole,
		Filename:   filepath.Dir(dir) + "\\log.txt",
	}
}
func init() {
	_init(defaultConfig())
}

// 默认初始化，重新初始化可以更改参数
func Init(c Config) {
	cf := defaultConfig()
	if c.WriteFile {
		cf.WriteFile = c.WriteFile
	}
	if c.FileFormat != "" {
		cf.FileFormat = c.FileFormat
	}
	if c.Filename != "" {
		cf.Filename = c.Filename
	}
	_init(cf)
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
func Fatal(args ...any) {
	log.Fatal(msg(args...))
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
func FatalSkip(skip int, args ...any) {
	log.WithOptions(zap.AddCallerSkip(skip)).Fatal(msg(args...))
}
func msg(args ...any) string {
	m := fmt.Sprintln(args...)
	return m[:len(m)-1]
}
