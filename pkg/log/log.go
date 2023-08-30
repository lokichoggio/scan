package log

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(logCfg *Config) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   //指定时间格式
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.LevelKey = "lvl"
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.CallerKey = "caller"

	cfg.OutputPaths = []string{
		logCfg.File,
	}
	if logCfg.Stdout {
		cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
	}

	cfg.Encoding = logCfg.Encoding

	var encoder zapcore.Encoder
	if cfg.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	}

	//日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //error级别
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool { //info和debug级别,debug级别是最低的
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logCfg.File,       //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    logCfg.MaxSize,    //文件大小限制,单位MB
		MaxBackups: logCfg.MaxBackups, //最大保留日志文件数量
		MaxAge:     logCfg.MaxAge,     //日志文件保留天数
		Compress:   logCfg.Compress,   //是否压缩处理
	})

	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logCfg.ErrFile,    //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    logCfg.MaxSize,    //文件大小限制,单位MB
		MaxBackups: logCfg.MaxBackups, //最大保留日志文件数量
		MaxAge:     logCfg.MaxAge,     //日志文件保留天数
		Compress:   logCfg.Compress,   //是否压缩处理
	})

	//第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	infoSyncer := []zapcore.WriteSyncer{infoFileWriteSyncer}
	errorSyncer := []zapcore.WriteSyncer{errorFileWriteSyncer}
	if logCfg.Stdout {
		infoSyncer = append(infoSyncer, zapcore.AddSync(os.Stdout))
		errorSyncer = append(errorSyncer, zapcore.AddSync(os.Stdout))
	}
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoSyncer...), lowPriority)
	errFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorSyncer...), highPriority)

	core := zapcore.NewTee(infoFileCore, errFileCore)
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = log
}

func With(fields ...interface{}) *zap.SugaredLogger {
	return logger.Sugar().With(fields...)
}

func WithFields(fields ...zapcore.Field) *zap.SugaredLogger {
	return logger.With(fields...).Sugar()
}

func Info(args ...interface{}) {
	logger.Sugar().Info(args...)
}

// func WithTraceIdInfo(traceId string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Info(args...)
// }

func Infof(template string, args ...interface{}) {
	logger.Sugar().Infof(template, args...)
}

// func WithTraceIdInfof(traceId string, template string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Infof(template, args...)
// }

func Debug(args ...interface{}) {
	logger.Sugar().Debug(args...)
}

// func WithTraceIdDebug(traceId string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Debug(args...)
// }

func Debugf(template string, args ...interface{}) {
	logger.Sugar().Debugf(template, args...)
}

// func WithTraceIdDebugf(traceId string, template string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Debugf(template, args...)
// }

func Warn(args ...interface{}) {
	logger.Sugar().Warn(args...)
}

// func WithTraceIdWarn(traceId string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Warn(args...)
// }

func Warnf(template string, args ...interface{}) {
	logger.Sugar().Warnf(template, args...)
}

// func WithTraceIdWarnf(traceId string, template string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Warnf(template, args...)
// }

func Error(args ...interface{}) {
	logger.Sugar().Error(args...)
}

// func WithTraceIdError(traceId string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Error(args...)
// }

func Errorf(template string, args ...interface{}) {
	logger.Sugar().Errorf(template, args...)
}

// func WithTraceIdErrorf(traceId string, template string, args ...interface{}) {
// 	logger.With(zap.String("traceId", traceId)).Sugar().Errorf(template, args...)
// }
