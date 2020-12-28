package log

import (
	"context"
	"fmt"
	"github.com/zer0131/toolbox/log/logrus_wrap"
	"google.golang.org/grpc/metadata"
	"math/rand"
	"net/http"
	"time"
)

const (
	LogIDKey                      = "log-id"
	RemoteAddrName                = "remote-addr"
	DefaultFileWriterMaxBackupDay = 7
	DefaultFileWriterMaxLength    = 8192
)

// 这个 logger 用于包外访问，方便大家自定义日志路径与文件名等信息
type Logger struct {
	logger *logrus_wrap.Logger
}

func GenLogId() string {
	var t int64 = time.Now().UnixNano() / 1000000
	var r int = rand.Intn(10000)
	return fmt.Sprintf("%d%d", t, r)
}

func NewContextWithLogID(ctx context.Context) context.Context {
	md := metadata.New(make(map[string]string))
	md.Set(LogIDKey, GenLogId())
	return metadata.NewIncomingContext(ctx, md)
}

func NewGrpcContextWithLogID(ctx context.Context) context.Context {
	logid := GenLogId()

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md.Set(LogIDKey, logid)
	}

	return metadata.AppendToOutgoingContext(ctx, LogIDKey, logid)
}

func NewContextWithSpecifyLogID(ctx context.Context, logId string) context.Context {
	md := metadata.New(make(map[string]string))
	md.Set(LogIDKey, logId)
	return metadata.NewIncomingContext(ctx, md)
}

func NewContextWithHttpReq(ctx context.Context, r *http.Request) context.Context {
	logId := r.Header.Get(LogIDKey)
	return NewContextWithSpecifyLogID(ctx, logId)
}

func Init(project string) error {
	if err := InitV4(WithProject(project)); err != nil {
		return err
	}
	return nil
}

func InitV2(project string) error {
	if err := InitV4(WithProject(project)); err != nil {
		return err
	}
	return nil
}

func InitV3(project string, logSize int) error {
	if err := InitV4(WithProject(project), WithMaxLength(int64(logSize))); err != nil {
		return err
	}
	return nil
}

func InitV4(opt ...logOptionsFunc) error {
	opts := defaultLogOptions
	for _, o := range opt {
		o(&opts)
	}

	if err := logrus_wrap.NewLogger(logrus_wrap.WithPath(opts.path), logrus_wrap.WithApp(opts.app), logrus_wrap.WithMaxLength(opts.maxLength), logrus_wrap.WithExpireDay(opts.expireDay)); err != nil {
		return err
	}
	return nil
}

func NewCustomLogger(opt ...logOptionsFunc) (logger *Logger, err error) {
	opts := defaultLogOptions
	for _, o := range opt {
		o(&opts)
	}

	logObj, err := logrus_wrap.NewCustomLogger(logrus_wrap.WithPath(opts.path), logrus_wrap.WithApp(opts.app), logrus_wrap.WithMaxLength(opts.maxLength), logrus_wrap.WithExpireDay(opts.expireDay), logrus_wrap.WithLogLevel(opts.level))
	if err != nil {
		return nil, err
	}

	return &Logger{logObj}, nil
}

type logOptions struct {
	path      string
	app       string
	level     string
	expireDay int
	maxLength int
}

var defaultLogOptions = logOptions{
	path:      "./log",
	app:       "app",
	level:     logrus_wrap.LevelDebug,
	expireDay: DefaultFileWriterMaxBackupDay,
	maxLength: DefaultFileWriterMaxLength,
}

type logOptionsFunc func(*logOptions)

func WithPath(v string) logOptionsFunc {
	return func(o *logOptions) {
		o.path = v
	}
}

func WithProject(v string) logOptionsFunc {
	return func(o *logOptions) {
		o.app = v
	}
}

func WithExpireDay(v int64) logOptionsFunc {
	return func(o *logOptions) {
		o.expireDay = int(v)
	}
}

func WithMaxLength(v int64) logOptionsFunc {
	return func(o *logOptions) {
		o.maxLength = int(v)
	}
}

func WithLogLevel(v string) logOptionsFunc {
	return func(o *logOptions) {
		o.level = v
	}
}

func SetLogLevel(l string) {
	logrus_wrap.SetLevel(l)
}

func LogIdFromContext(ctx context.Context) (string, bool) {
	md, _ := metadata.FromIncomingContext(ctx)
	arr := md.Get(LogIDKey)
	if len(arr) == 1 {
		return arr[0], true
	}
	return "", false
}

func RemoteAddrNameFromContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	remoteIp := md.Get(RemoteAddrName)
	if len(remoteIp) == 1 {
		return remoteIp[0]
	}
	return ""
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Debugf(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logrus_wrap.Debugf(format, v...)
	}
}

func Debug(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Debug(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logrus_wrap.Debug(v...)
	}
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Infof(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logrus_wrap.Infof(format, v...)
	}
}

func Info(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Info(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logrus_wrap.Info(v...)
	}
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Warnf(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logrus_wrap.Warnf(format, v...)
	}
}

func Warn(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Warn(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logrus_wrap.Warn(v...)
	}
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Errorf(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logrus_wrap.Errorf(format, v...)
	}
}

func Error(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logrus_wrap.Error(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logrus_wrap.Error(v...)
	}
}

func Close() {
}

func InfofArray(ctx context.Context, format string, array []string, splitSize int) {
	if len(array) == 0 {
		return
	}
	logId, ok := LogIdFromContext(ctx)

	if len(array) < splitSize {
		if ok {
			logrus_wrap.Infof(fmt.Sprintf("[%s] %s", logId, format), array)
		} else {
			logrus_wrap.Infof(format, array)
		}
		return
	}

	var tmp []string
	for _, v := range array {
		tmp = append(tmp, v)
		if len(tmp) == splitSize {
			if ok {
				logrus_wrap.Infof(fmt.Sprintf("[%s] %s", logId, format), tmp)
			} else {
				logrus_wrap.Infof(format, tmp)
			}

			tmp = tmp[:0]
		}
	}

	if len(tmp) > 0 {
		if ok {
			logrus_wrap.Infof(fmt.Sprintf("[%s] %s", logId, format), tmp)
		} else {
			logrus_wrap.Infof(format, tmp)
		}
	}
}

func (logObj *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Debugf(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logObj.logger.Debugf(format, v...)
	}
}

func (logObj *Logger) Debug(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Debug(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logObj.logger.Debug(v...)
	}
}

func (logObj *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Infof(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logObj.logger.Infof(format, v...)
	}
}

func (logObj *Logger) Info(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Info(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logObj.logger.Info(v...)
	}
}

func (logObj *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Warnf(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logObj.logger.Warnf(format, v...)
	}
}

func (logObj *Logger) Warn(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Warn(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logObj.logger.Warn(v...)
	}
}

func (logObj *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Errorf(fmt.Sprintf("[%s] %s", logId, format), v...)
	} else {
		logObj.logger.Errorf(format, v...)
	}
}

func (logObj *Logger) Error(ctx context.Context, v ...interface{}) {
	logId, ok := LogIdFromContext(ctx)
	if ok {
		logObj.logger.Error(fmt.Sprintf("[%s]", logId) + " " + fmt.Sprint(v...))
	} else {
		logObj.logger.Error(v...)
	}
}

func (logObj *Logger) InfofArray(ctx context.Context, format string, array []string, splitSize int) {
	if len(array) == 0 {
		return
	}
	logId, ok := LogIdFromContext(ctx)

	if len(array) < splitSize {
		if ok {
			logObj.logger.Infof(fmt.Sprintf("[%s] %s", logId, format), array)
		} else {
			logObj.logger.Infof(format, array)
		}
		return
	}

	var tmp []string
	for _, v := range array {
		tmp = append(tmp, v)
		if len(tmp) == splitSize {
			if ok {
				logObj.logger.Infof(fmt.Sprintf("[%s] %s", logId, format), tmp)
			} else {
				logObj.logger.Infof(format, tmp)
			}

			tmp = tmp[:0]
		}
	}

	if len(tmp) > 0 {
		if ok {
			logObj.logger.Infof(fmt.Sprintf("[%s] %s", logId, format), tmp)
		} else {
			logObj.logger.Infof(format, tmp)
		}
	}
}
