package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// 日志处理建议 https://learnku.com/articles/42231
var log *zap.Logger

func init() {
	consoleEncoder := zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:       "ts",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			FunctionKey:   zapcore.OmitKey,
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    "\r\n",
			//EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			//	enc.AppendString("[" + level.CapitalString() + "]")
			//},
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stdout),
		zap.DebugLevel,
	)
	log = zap.New(zapcore.NewTee(consoleCore), zap.AddCaller(), zap.AddCallerSkip(1))
}

func Info(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	log.Info(msg[:len(msg)-1])
}

func Error(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	log.Error(msg[:len(msg)-1])
}
func Warn(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	log.Warn(msg[:len(msg)-1])
}
func Debug(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	log.Debug(msg[:len(msg)-1])
}
