package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"cnores-skeleton-golang-app/app/shared/utils/stack_trace_error"
	"os"
)

type Logger struct {
	log *zap.Logger
}

type Fields map[string]string

var instance *Logger

func getLogLevel(logLevel string) zapcore.Level {
	zapLogLevel := zap.DebugLevel
	switch logLevel {
	case "INFO":
		zapLogLevel = zap.InfoLevel
	case "DEBUG":
		zapLogLevel = zap.DebugLevel
	case "WARN":
		zapLogLevel = zap.WarnLevel
	case "ERROR":
		zapLogLevel = zap.ErrorLevel
	default:
		zapLogLevel = zap.DebugLevel
	}
	return zapLogLevel
}

func init() {
	log := initLoggerZap()
	log.Info(fmt.Sprintf("Logger loaded successfully with level: %s", getLogLevel(os.Getenv("LOG_LEVEL"))))

	instance = &Logger{log: log}
}

func initLoggerZap() *zap.Logger {
	cfg := zap.Config{
		Encoding:         "json",
		DisableCaller:    false,
		Level:            zap.NewAtomicLevelAt(getLogLevel(os.Getenv("LOG_LEVEL"))),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	log, _ := cfg.Build()
	log = log.With(zap.String("logtopic", "_test.logs"))
	zap.AddCallerSkip(1)
	zap.ReplaceGlobals(log)
	return log
}

func Field(keyInput string, valueInput string) Fields { return instance.Field(keyInput, valueInput) }
func (logger *Logger) Field(keyInput string, valueInput string) Fields {
	return Fields{
		keyInput: valueInput,
	}
}

func NewLogger(ctx context.Context, fields Fields) *Logger {
	logCustom := &Logger{
		log: initLoggerZap(),
	}
	logTraceParent(ctx, logCustom)
	logAdditionalFieldsFomContext(ctx, logCustom)
	for k, v := range fields {
		logCustom.log = logCustom.log.With(zap.String(k, v))
	}
	return logCustom
}

func WithFields(fields Fields) *Logger { return instance.WithFields(fields) }

func (logger *Logger) WithFields(fields Fields) *Logger {
	logCustom := &Logger{
		log: initLoggerZap(),
	}
	for k, v := range fields {
		logCustom.log = logCustom.log.With(zap.String(k, v))
	}
	return logCustom
}

func logTraceParent(ctx context.Context, logger *Logger) {
	traceVersion, _ := ctx.Value("TraceVersion").(string)
	logger.log = logger.log.With(zap.String("TraceVersion", traceVersion))

	traceID, _ := ctx.Value("TraceID").(string)
	logger.log = logger.log.With(zap.String("TraceID", traceID))

	parentID, _ := ctx.Value("ParentID").(string)
	logger.log = logger.log.With(zap.String("ParentID", parentID))

	traceFlags, _ := ctx.Value("TraceFlags").(string)
	logger.log = logger.log.With(zap.String("TraceFlags", traceFlags))

	spanID, _ := ctx.Value("SpanID").(string)
	logger.log = logger.log.With(zap.String("SpanID", spanID))
}

func logAdditionalFieldsFomContext(ctx context.Context, logger *Logger) {

	if fieldsToLog := ctx.Value("fieldsToLog"); fieldsToLog != nil {
		fieldsToLogMap := fieldsToLog.(map[string]string)
		for name, value := range fieldsToLogMap {
			logger.log = logger.log.With(zap.String(name, value))
		}
	}
}

func WithError(err error) *Logger { return instance.WithError(err) }
func (logger *Logger) WithError(err error) *Logger {

	logCustom := &Logger{
		log: initLoggerZap(),
	}
	logCustom.log = logCustom.log.With(zap.String("error", err.Error()))
	logCustom.log = logCustom.log.With(zap.Strings("stackTrace", stack_trace_error.GetStackTraceFromError(err)))
	return logCustom
}

func Info(message string, args ...interface{}) { instance.Info(message, args...) }
func (logger *Logger) Info(message string, args ...interface{}) {
	logger.log.Info(fmt.Sprintf(message, args...))
}

func Error(message string, args ...interface{}) { instance.Error(message, args...) }
func (logger *Logger) Error(message string, args ...interface{}) {
	logger.log.Error(fmt.Sprintf(message, args...))
}

func Debug(message string, args ...interface{}) { instance.Debug(message, args...) }
func (logger *Logger) Debug(message string, args ...interface{}) {
	logger.log.Debug(fmt.Sprintf(message, args...))
}

func Warn(message string, args ...interface{}) { instance.Warn(message, args...) }
func (logger *Logger) Warn(message string, args ...interface{}) {
	logger.log.Warn(fmt.Sprintf(message, args...))
}

func Fatal(message string, args ...interface{}) { instance.Fatal(message, args...) }
func (logger *Logger) Fatal(message string, args ...interface{}) {
	logger.log.Fatal(fmt.Sprintf(message, args...))
}
