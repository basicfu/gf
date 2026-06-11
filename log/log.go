package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func EncodeCallerRelToWD(width int) zapcore.CallerEncoder {
	wd, _ := os.Getwd()
	wd = filepath.Clean(wd)
	return func(c zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		abs := filepath.Clean(c.File)
		rel, err := filepath.Rel(wd, abs)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			enc.AppendString(fmt.Sprintf("%s:%d", filepath.ToSlash(abs), c.Line))
			return
		}
		out := fmt.Sprintf("%s:%d", filepath.ToSlash(rel), c.Line)
		if width > 0 && len(out) < width {
			out += strings.Repeat(" ", width-len(out))
		}
		enc.AppendString(out)
	}
}

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
	EncodeCaller:   EncodeCallerRelToWD(16),
	//EncodeCaller: EncodeCallerOSC8(),//goland官方说明在2025.3修复，目前测试未修复待修复后使用这种方式
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
	LogDir     string
	MaxAge     int
}

func _init(c Config) {
	jsonEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(t.UnixMilli())
	}
	consoleEncoder.LineEnding = "\r\n"
	consoleEncoder.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
	}
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoder), zapcore.Lock(os.Stdout), zap.DebugLevel)
	if c.WriteFile {
		rw := &rotatingWriter{dir: c.LogDir, maxAge: c.MaxAge, nowFunc: time.Now}
		encoder := zapcore.NewConsoleEncoder(consoleEncoder)
		if c.FileFormat == FileFormatJson {
			encoder = zapcore.NewJSONEncoder(jsonEncoder)
		}
		fileCore := zapcore.NewCore(encoder, zapcore.AddSync(rw), zap.DebugLevel)
		log = zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		//log = zap.New(zapcore.NewTee(consoleCore), zap.AddCaller(), zap.AddCallerSkip(1))
		log = zap.New(
			zapcore.NewTee(consoleCore),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return forcedCallerCore{Core: core}
			}),
		)
	}
}
func defaultConfig() Config {
	dir, _ := os.Executable()
	return Config{
		WriteFile:  false,
		FileFormat: FileFormatConsole,
		LogDir:     filepath.Join(filepath.Dir(dir), "logs"),
		MaxAge:     180,
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
	if c.LogDir != "" {
		cf.LogDir = c.LogDir
	}
	if c.MaxAge != 0 {
		cf.MaxAge = c.MaxAge
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

const (
	forcedCallerFileKey = "__forced_caller_file__"
	forcedCallerLineKey = "__forced_caller_line__"
	forcedCallerFuncKey = "__forced_caller_func__"
)

func WithForcedCaller(file string, line int, function string) *zap.Logger {
	return log.With(
		zap.String(forcedCallerFileKey, file),
		zap.Int(forcedCallerLineKey, line),
		zap.String(forcedCallerFuncKey, function),
	)
}

type forcedCallerState struct {
	File string
	Line int
	Func string
	Ok   bool
}

type forcedCallerCore struct {
	zapcore.Core
	fc forcedCallerState
}

func (c forcedCallerCore) With(fields []zapcore.Field) zapcore.Core {
	var (
		next = c.fc
		out  = make([]zapcore.Field, 0, len(fields))
	)
	for i := range fields {
		f := fields[i]
		switch f.Key {
		case forcedCallerFileKey:
			if f.Type == zapcore.StringType {
				next.File = f.String
				next.Ok = true
				continue
			}
		case forcedCallerFuncKey:
			if f.Type == zapcore.StringType {
				next.Func = f.String
				next.Ok = true
				continue
			}
		case forcedCallerLineKey:
			if f.Type == zapcore.Int64Type {
				next.Line = int(f.Integer)
				next.Ok = true
				continue
			}
		}
		out = append(out, f)
	}
	return forcedCallerCore{
		Core: c.Core.With(out),
		fc:   next,
	}
}

func (c forcedCallerCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c forcedCallerCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	if c.fc.Ok && c.fc.File != "" && c.fc.Line != 0 {
		ent.Caller = zapcore.EntryCaller{
			Defined:  true,
			File:     filepath.Clean(c.fc.File),
			Line:     c.fc.Line,
			Function: c.fc.Func,
		}
	}
	if len(fields) > 0 {
		out := make([]zapcore.Field, 0, len(fields))
		for i := range fields {
			switch fields[i].Key {
			case forcedCallerFileKey, forcedCallerLineKey, forcedCallerFuncKey:
				continue
			default:
				out = append(out, fields[i])
			}
		}
		fields = out
	}
	return c.Core.Write(ent, fields)
}
