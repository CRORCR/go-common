package logger

import (
	"bytes"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/sirupsen/logrus"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Formatter - 自定义 logrus 格式化器
// 格式：时间 [级别][trace_id][文件位置] 消息
type Formatter struct {
	// FieldsOrder - 字段排序，默认: ["trace", "caller"]
	FieldsOrder []string

	// TimestampFormat - 时间格式，默认: "2006-01-02 15:04:05"
	TimestampFormat string

	// HideKeys - 显示 [fieldValue] 而不是 [fieldKey:fieldValue]
	HideKeys bool

	// NoColors - 禁用颜色
	NoColors bool

	// NoFieldsColors - 仅对级别应用颜色
	NoFieldsColors bool

	// NoFieldsSpace - 字段之间无空格
	NoFieldsSpace bool

	// ShowFullLevel - 显示完整级别 [WARNING] 而不是 [WARN]
	ShowFullLevel bool

	// NoUppercaseLevel - 级别不大写
	NoUppercaseLevel bool

	// TrimMessages - 修剪消息前后空格
	TrimMessages bool

	// CallerFirst - 在消息前打印调用者信息
	CallerFirst bool

	// CustomCallerFormatter - 自定义调用者格式化器
	CustomCallerFormatter func(*runtime.Frame) string

	// FieldsIgnore - 忽略的字段
	FieldsIgnore []string
}

// Format 格式化日志条目
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// 输出缓冲区
	b := &bytes.Buffer{}

	// 写入时间
	b.WriteString(entry.Time.Format(timestampFormat))

	// 写入级别
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	if f.CallerFirst {
		f.writeCaller(b, entry)
	}

	if !f.NoColors {
		fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("]")

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// 删除忽略的字段
	for _, v := range f.FieldsIgnore {
		delete(entry.Data, v)
	}

	// 写入字段
	if f.FieldsOrder == nil {
		f.writeFields(b, entry)
	} else {
		f.writeOrderedFields(b, entry)
	}

	if f.NoFieldsSpace {
		b.WriteString(" ")
	}

	if !f.NoColors && !f.NoFieldsColors {
		b.WriteString("\x1b[0m")
	}

	// 写入消息
	if f.TrimMessages {
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(entry.Message)
	}

	if !f.CallerFirst {
		f.writeCaller(b, entry)
	}

	fstr := b.String()
	// 错误日志如果太长，截断到 512 字符
	if entry.Level == logrus.ErrorLevel && len(fstr) > 512 {
		fstr = strutil.Substring(fstr, 0, 512)
		b = &bytes.Buffer{}
		b.WriteString(fstr)
		b.WriteByte('\n')
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				" (%s:%d %s)",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if slice.Contain(f.FieldsOrder, field) {
		if f.HideKeys {
			fmt.Fprintf(b, "[%v]", entry.Data[field])
		} else {
			fmt.Fprintf(b, "[%s:%v]", field, entry.Data[field])
		}
	} else {
		fmt.Fprintf(b, " %s:%v,", field, entry.Data[field])
	}

	if !f.NoFieldsSpace {
		b.WriteString(" ")
	}
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return colorGray
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorBlue
	}
}
