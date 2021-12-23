package logging

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-zentao-task/pkg/config"
	"go-zentao-task/pkg/gredis"
	"io"
	"os"
	"time"
)

type Level uint32
type Output uint32

var (
	topic  string
	env    string
	output Output
)

const (
	FatalLevel Level = iota + 1
	ErrorLevel
	WarnLevel
	InfoLevel

	Stdout Output = iota + 1
	File
)

type LogItem struct {
	Topic      string `json:"topic"`
	Level      Level  `json:"level"`
	V1         string `json:"v1"`
	V2         string `json:"v2"`
	V3         string `json:"v3"`
	Message    string `json:"message"`
	Request    string `json:"request"`
	Response   string `json:"response"`
	CreateTime string `json:"create_time"`
}

func newLogItem(v1, v2, v3, message string, request, response interface{}, level Level) *LogItem {
	var req string
	var resp string

	if v, ok := request.(string); ok {
		req = v
	} else {
		b, _ := json.Marshal(request)
		req = string(b)
	}

	if v, ok := response.(string); ok {
		resp = v
	} else {
		b, _ := json.Marshal(response)
		resp = string(b)
	}

	item := &LogItem{
		Topic:      topic,
		Level:      level,
		V1:         v1,
		V2:         v2,
		V3:         v3,
		Message:    message,
		Request:    req,
		Response:   resp,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	return item
}

func Setup(environment string, o Output) {
	env = environment
	topic = config.Get("app.log.topic")
	output = o
}

type Logging struct {
	V1 string
	V2 string
	V3 string
}

func (l *Logging) Info(message string, request, response interface{}) {
	Info(l.V1, l.V2, l.V3, message, request, response)
}

func (l *Logging) Warn(message string, request, response interface{}) {
	Warn(l.V1, l.V2, l.V3, message, request, response)
}

func (l *Logging) ErrorL(message string, request, response interface{}) {
	Error(l.V1, l.V2, l.V3, message, request, response)
}

func (l *Logging) Fatal(message string, request, response interface{}) {
	Fatal(l.V1, l.V2, l.V3, message, request, response)
}

func Info(v1, v2, v3, message string, request, response interface{}) {
	handle(newLogItem(v1, v2, v3, message, request, response, InfoLevel))
}

func Warn(v1, v2, v3, message string, request, response interface{}) {
	handle(newLogItem(v1, v2, v3, message, request, response, WarnLevel))
}

func Error(v1, v2, v3, message string, request, response interface{}) {
	handle(newLogItem(v1, v2, v3, message, request, response, ErrorLevel))
}

func Fatal(v1, v2, v3, message string, request, response interface{}) {
	handle(newLogItem(v1, v2, v3, message, request, response, FatalLevel))
}

func handle(item *LogItem) {
	if env == "development" {
		log2local(item)
	} else {
		pushLogList(item)
	}
}

func pushLogList(item *LogItem) {
	conn := gredis.RedisLogPool.Get()
	defer conn.Close()

	info, _ := json.Marshal(item)
	conn.Do("lpush", "sys:log:list", string(info)) //nolint
}

func log2local(item *LogItem) {
	var out io.Writer
	out = os.Stdout
	if output == File {
		file, err := os.OpenFile("runtime/logs/app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			out = file
		}
	}
	log.SetOutput(out)

	logger := log.WithFields(log.Fields{
		"topic":       item.Topic,
		"level":       item.Level,
		"v1":          item.V1,
		"v2":          item.V2,
		"v3":          item.V3,
		"request":     item.Request,
		"response":    item.Response,
		"create_time": item.CreateTime,
	})

	switch item.Level {
	case FatalLevel:
		fallthrough
	case ErrorLevel:
		logger.Error(item.Message)
	case WarnLevel:
		logger.Warn(item.Message)
	case InfoLevel:
		logger.Info(item.Message)
	}
}
