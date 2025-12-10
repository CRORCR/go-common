package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

// SetUpLog 配置日志
// 使用自定义的 logrus formatter，格式：时间 [级别][trace_id][文件位置] 消息
func SetUpLog(log logx.LogConf) {
	writer := NewLogrusWriter(func(logger *logrus.Logger) {
		// 设置日志等级
		level := logrus.InfoLevel
		switch log.Level {
		case "debug":
			level = logrus.DebugLevel
		case "info":
			level = logrus.InfoLevel
		case "warn":
			level = logrus.WarnLevel
		case "error":
			level = logrus.ErrorLevel
		case "fatal":
			level = logrus.FatalLevel
		case "panic":
			level = logrus.PanicLevel
		}

		logger.SetLevel(level)

		if log.Mode == "file" {
			// 设置分割的文件
			accessfile := "access"
			errorfile := "error"
			if log.Path == "" {
				accessfile = "logs/" + accessfile + "_%Y_%m_%d.log"
				errorfile = "logs/" + errorfile + "_%Y_%m_%d.log"
			} else {
				accessfile = log.Path + "/" + accessfile + "_%Y_%m_%d.log"
				errorfile = log.Path + "/" + errorfile + "_%Y_%m_%d.log"
			}

			// WithMaxAge 设置最大保存时间（默认 3 天）
			maxAge := 3 * 24 * time.Hour
			if log.KeepDays != 0 {
				maxAge = time.Duration(log.KeepDays) * 24 * time.Hour
			}

			// 设置日志分割的时间，这里设置为一天分割一次
			rotationTime := 24 * time.Hour

			accessRotateLog, err := rotatelogs.New(
				accessfile,
				rotatelogs.WithMaxAge(maxAge),
				rotatelogs.WithRotationTime(rotationTime),
			)

			if err != nil {
				logger.Errorf("config rotetelog failed err: %v", err)
			}

			errorRotateLog, err := rotatelogs.New(
				errorfile,
				rotatelogs.WithMaxAge(maxAge),
				rotatelogs.WithRotationTime(rotationTime),
			)

			if err != nil {
				logger.Errorf("config rotetelog failed err: %v", err)
			}

			// 对不同的等级 log 设置 io.writer
			writerMap := lfshook.WriterMap{
				logrus.DebugLevel: accessRotateLog,
				logrus.InfoLevel:  accessRotateLog,
				logrus.WarnLevel:  accessRotateLog,
				logrus.ErrorLevel: errorRotateLog,
				logrus.FatalLevel: errorRotateLog,
				logrus.PanicLevel: errorRotateLog,
			}

			lfsHook := lfshook.NewHook(writerMap, &Formatter{
				FieldsOrder:     []string{"trace", "caller"},
				FieldsIgnore:    []string{"span", "duration"},
				HideKeys:        true,
				NoFieldsSpace:   true,
				TrimMessages:    true,
				NoColors:        true,
				NoFieldsColors:  true,
				TimestampFormat: "2006-01-02 15:04:05",
			})

			logger.AddHook(lfsHook)
		} else {
			// console 模式
			logger.SetFormatter(&Formatter{
				FieldsOrder:     []string{"trace", "caller"},
				FieldsIgnore:    []string{"span", "duration"},
				HideKeys:        true,
				NoFieldsSpace:   true,
				TrimMessages:    true,
				NoColors:        false,
				NoFieldsColors:  false,
				TimestampFormat: "2006-01-02 15:04:05",
			})
		}
	})

	logx.SetWriter(writer)
}
