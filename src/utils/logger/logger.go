package logger

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CSugaredLogger is a superset of zap.SugaredLogger
type CSugaredLogger struct {
	zap.SugaredLogger
}

// CLogger is a superset of zap.Logger
type CLogger struct {
	zap.Logger
}

var logger *CLogger
var correlationIdContextKey string
var correlationIdFieldKey string

func (l *CLogger) WithContextCorrelationId(ctx context.Context) *CLogger {
	correlationId := ctx.Value(correlationIdContextKey)
	return l.WithCorrelationId(correlationId)
}


func (l *CSugaredLogger) WithContextCorrelationId(ctx context.Context) *CSugaredLogger {
	correlationId := ctx.Value(correlationIdContextKey)
	return l.WithCorrelationId(correlationId)
}

func (l *CLogger) WithCorrelationId(correlationId interface{}) *CLogger {
	if s, ok := correlationId.(string); ok {
		return &CLogger{*l.With(zap.Stringp(correlationIdFieldKey, &s))}
	}
	return l
}

func (l *CSugaredLogger) WithCorrelationId(correlationId interface{}) *CSugaredLogger {
	if s, ok := correlationId.(string); ok {
		return &CSugaredLogger{*l.With(zap.Stringp(correlationIdFieldKey, &s))}
	}
	return l
}

func SugaredLogger() *CSugaredLogger {
	if logger == nil {
		panic("logger not initialized. Call Init(ctx)")
	}
	l := logger.Sugar()
	return &CSugaredLogger{*l}
}

func Logger() *CLogger {
	if logger == nil {
		panic("logger not initialized. Call Init(ctx)")
	}
	return logger
}

func SetCorrelationIdFieldKey(key string) {
	if key == "" {
		return
	}
	correlationIdFieldKey = key
}

func SetCorrelationIdContextKey(key string) {
	if key == "" {
		return
	}
	correlationIdContextKey = key
}

func Init(ctx context.Context, developmentMode bool) {
	if logger != nil {
		return
	}
	var (
		zapConfig     zap.Config
		encoderConfig zapcore.EncoderConfig
		atom          zap.AtomicLevel
		loggerMode    []string
	)

	correlationIdContextKey = "correlation_id"
	correlationIdFieldKey = "correlation_id"

	if developmentMode {
		loggerMode = append(loggerMode, "dev")
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    "func",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		}
		zapConfig = zap.Config{
			Level:             atom,
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderConfig,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stdout"},
			InitialFields:     nil,
		}
	} else {
		loggerMode = append(loggerMode, "prod")
		atom = zap.NewAtomicLevelAt(zap.InfoLevel)
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		zapConfig = zap.Config{
			Level:             atom,
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          &zap.SamplingConfig{Initial: 100, Thereafter: 100},
			Encoding:          "json",
			EncoderConfig:     encoderConfig,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stdout"},
			InitialFields:     nil,
		}
	}

	l, err := zapConfig.Build()
	if err != nil {
		panic(fmt.Sprintf("logger initalization error: %s", err.Error()))
	}

	l.Info("Logger initialized successfully", zap.Strings("logger_modes", loggerMode))

	logger = &CLogger{*l}
}
