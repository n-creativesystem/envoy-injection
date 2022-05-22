package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Infof(msg string, args ...interface{})
	Warningf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
	SetExit(func(code int))
}

type Log struct {
	*zap.Logger
	exit func(code int)
}

func init() {
	var zapLog *zap.Logger
	isProduction := os.Getenv("APP_ENV") == "PRODUCTION"
	var (
		zapEncoder zapcore.EncoderConfig
		zapConfig  zap.Config
	)
	if isProduction {
		zapEncoder = zap.NewProductionEncoderConfig()
		zapConfig = zap.NewProductionConfig()
	} else {
		zapEncoder = zap.NewDevelopmentEncoderConfig()
		zapConfig = zap.NewDevelopmentConfig()
	}
	zapEncoder.EncodeTime = zapcore.RFC3339TimeEncoder
	zapConfig.EncoderConfig = zapEncoder
	zapLog, _ = zapConfig.Build()
	log = &Log{
		Logger: zapLog,
	}
}

var (
	log *Log
)

func (log *Log) Infof(msg string, args ...interface{}) {
	log.Logger.Info(msg, log.argsToFields(args)...)
}

func (log *Log) Warningf(msg string, args ...interface{}) {
	log.Logger.Info(msg, log.argsToFields(args)...)
}

func (log *Log) Errorf(msg string, args ...interface{}) {
	log.Logger.Info(msg, log.argsToFields(args)...)
}

func (log *Log) Debugf(msg string, args ...interface{}) {
	log.Logger.Info(msg, log.argsToFields(args)...)
}

func (log *Log) Fatalf(msg string, args ...interface{}) {
	log.Logger.Fatal(msg, log.argsToFields(args)...)
	if log.exit != nil {
		log.exit(1)
	} else {
		os.Exit(1)
	}
}

func (log *Log) SetExit(fn func(code int)) {
	log.exit = fn
}

func (log *Log) argsToFields(args []interface{}) []zapcore.Field {
	l := len(args)
	fields := make([]zapcore.Field, 0, l)
	for i := 0; i < l; i += 2 {
		end := i + 2
		if l < end {
			end = l
		}
		kv := args[i:end]
		if len(kv) == 2 {
			k, _ := kv[0].(string)
			field := zapcore.Field{
				Key:       k,
				Interface: kv[1],
			}
			fields = append(fields, field)
		}
	}
	return fields
}

func Infof(msg string, args ...interface{})    { log.Infof(msg, args...) }
func Warningf(msg string, args ...interface{}) { log.Warningf(msg, args...) }
func Errorf(msg string, args ...interface{})   { log.Errorf(msg, args...) }
func Debugf(msg string, args ...interface{})   { log.Debugf(msg, args...) }
func Fatalf(msg string, args ...interface{})   { log.Fatalf(msg, args...) }

func Info(msg string, args ...interface{})    { log.Infof(msg, args...) }
func Warning(msg string, args ...interface{}) { log.Warningf(msg, args...) }
func Error(msg string, args ...interface{})   { log.Errorf(msg, args...) }
func Debug(msg string, args ...interface{})   { log.Debugf(msg, args...) }
func Fatal(msg string, args ...interface{})   { log.Fatalf(msg, args...) }

func GetZapLogger() *zap.Logger {
	return log.Logger
}

func GetLogger() Logger {
	return log
}
