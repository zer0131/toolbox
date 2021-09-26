package logrus_wrap

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const (
	defaultFileWriterMaxLength           = 8192
	defaultFileWriterRotationTime        = time.Hour
	defaultFileWriterExpireDay           = 7 // 日志保留 7 天后被删除
	defaultFileWriterMsgSuffixTimeString = "06-01-02 15:04:05.999"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

var levelMapperRev = map[string]logrus.Level{
	LevelDebug: logrus.DebugLevel,
	LevelInfo:  logrus.InfoLevel,
	LevelWarn:  logrus.WarnLevel,
	LevelError: logrus.ErrorLevel,
}

var logObj *logger = nil

type logger struct {
	iWriter *logrus.Logger
	fWriter *logrus.Logger
}

// 这个 logger 用于包外访问，方便大家自定义日志路径与文件名等信息
type Logger struct {
	logger *logrus.Logger
}

type logOptions struct {
	path      string
	app       string
	maxLength int
	expireDay int
	level     string
}

var defaultLogOptions = logOptions{
	path:      "./log",
	app:       "app",
	level:     LevelDebug,
	maxLength: defaultFileWriterMaxLength,
	expireDay: defaultFileWriterExpireDay,
}

type LogOptionsFunc func(*logOptions)

func WithPath(v string) LogOptionsFunc {
	return func(o *logOptions) {
		o.path = v
	}
}

func WithApp(v string) LogOptionsFunc {
	return func(o *logOptions) {
		o.app = v
	}
}

func WithExpireDay(v int) LogOptionsFunc {
	return func(o *logOptions) {
		o.expireDay = v
	}
}

func WithMaxLength(v int) LogOptionsFunc {
	return func(o *logOptions) {
		o.maxLength = v
	}
}

func WithLogLevel(l string) LogOptionsFunc {
	return func(o *logOptions) {
		o.level = l
	}
}

func NewLogger(opt ...LogOptionsFunc) error {
	if logObj != nil {
		fmt.Printf("[logrus] logObj is already initialized\n")
		return nil
	}
	opts := defaultLogOptions
	for _, o := range opt {
		o(&opts)
	}

	if err := os.MkdirAll(opts.path, os.FileMode(0755)); err != nil {
		fmt.Printf("[logrus] mkdir failed, path: %s\n", opts.path)
		return err
	}

	iWriter, err := newWriter(opts.path, opts.app+".log", opts.expireDay, opts.maxLength)
	if err != nil {
		return err
	}

	fWriter, err := newWriter(opts.path, opts.app+".log.wf", opts.expireDay, opts.maxLength)
	if err != nil {
		return err
	}

	logObj = &logger{
		iWriter: iWriter,
		fWriter: fWriter,
	}

	SetLevel(opts.level)
	return nil
}

func NewCustomLogger(opt ...LogOptionsFunc) (*Logger, error) {
	opts := defaultLogOptions
	for _, o := range opt {
		o(&opts)
	}

	if err := os.MkdirAll(opts.path, os.FileMode(0755)); err != nil {
		fmt.Printf("[logrus] mkdir failed, path: %s\n", opts.path)
		return nil, err
	}

	writer, err := newWriter(opts.path, opts.app+".log", opts.expireDay, opts.maxLength)
	if err != nil {
		return nil, err
	}

	if v, ok := levelMapperRev[opts.level]; ok {
		writer.SetLevel(v)
	} else {
		fmt.Printf("[logrus] unknown log level %s, now using DEBUG\n", opts.level)
		writer.SetLevel(logrus.DebugLevel)
	}

	return &Logger{
		writer,
	}, nil
}

func SetLevel(level string) {
	if logObj == nil {
		fmt.Printf("[logrus] logObj is uninitialized, set level error\n")
		return
	}
	if v, ok := levelMapperRev[level]; ok {
		logObj.iWriter.SetLevel(v)
		logObj.fWriter.SetLevel(v)
	} else {
		fmt.Printf("[logrus] unknown log level %s, now using DEBUG\n", level)
		logObj.iWriter.SetLevel(logrus.DebugLevel)
		logObj.fWriter.SetLevel(logrus.DebugLevel)
	}
}

func newWriter(filepath, fileName string, expireDay, maxLength int) (logger *logrus.Logger, err error) {
	// 转为绝对路径处理
	var fileWithFullPath string
	if strings.HasPrefix(filepath, "/") {
		fileWithFullPath = path.Join(filepath, fileName)
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("[logrus] os.getwd err, err: %s\n", err)
			return nil, err
		}
		fileWithFullPath = path.Join(pwd, filepath, fileName)
	}
	logger = logrus.New()
	writer, err := rotatelogs.New(
		fileWithFullPath+".%Y%m%d%H",
		rotatelogs.WithLinkName(fileWithFullPath),
		rotatelogs.WithMaxAge(time.Duration(expireDay)*24*time.Hour),
		rotatelogs.WithRotationTime(defaultFileWriterRotationTime),
	)
	if err != nil {
		fmt.Printf("[logrus] failed to create rotatelogs: %s\n", err)
		return nil, err
	}

	logger.SetOutput(writer)
	formatter := &textFormatter{maxLength: maxLength}
	logger.Formatter = formatter
	return
}

func Debug(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[DEBUG]"}, v...)...)
		return
	}
	if logObj.iWriter.IsLevelEnabled(logrus.DebugLevel) {
		logObj.iWriter.Debug(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[DEBUG] "+format+"\n", v...)
		return
	}
	if logObj.iWriter.IsLevelEnabled(logrus.DebugLevel) {
		logObj.iWriter.Debugf(format, v...)
	}
}

func Info(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[INFO]"}, v...)...)
		return
	}
	if logObj.iWriter.IsLevelEnabled(logrus.InfoLevel) {
		logObj.iWriter.Info(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[INFO] "+format+"\n", v...)
		return
	}
	if logObj.iWriter.IsLevelEnabled(logrus.InfoLevel) {
		logObj.iWriter.Infof(format, v...)
	}
}

func Warn(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[WARN]"}, v...)...)
		return
	}
	if logObj.fWriter.IsLevelEnabled(logrus.WarnLevel) {
		logObj.fWriter.Warn(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[WARN] "+format+"\n", v...)
		return
	}
	if logObj.fWriter.IsLevelEnabled(logrus.WarnLevel) {
		logObj.fWriter.Warnf(format, v...)
	}
}

func Error(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[ERROR]"}, v...)...)
		return
	}
	if logObj.fWriter.IsLevelEnabled(logrus.ErrorLevel) {
		logObj.fWriter.Error(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[ERROR] "+format+"\n", v...)
		return
	}
	if logObj.fWriter.IsLevelEnabled(logrus.ErrorLevel) {
		logObj.fWriter.Errorf(format, v...)
	}
}

func (logObj *Logger) Debug(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[DEBUG]"}, v...)...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.DebugLevel) {
		logObj.logger.Debug(v...)
	}
}

func (logObj *Logger) Debugf(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[DEBUG] "+format+"\n", v...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.DebugLevel) {
		logObj.logger.Debugf(format, v...)
	}
}

func (logObj *Logger) Info(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[INFO]"}, v...)...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.InfoLevel) {
		logObj.logger.Info(v...)
	}
}

func (logObj *Logger) Infof(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[INFO] "+format+"\n", v...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.InfoLevel) {
		logObj.logger.Infof(format, v...)
	}
}

func (logObj *Logger) Warn(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[WARN]"}, v...)...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.WarnLevel) {
		logObj.logger.Warn(v...)
	}
}

func (logObj *Logger) Warnf(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[WARN] "+format+"\n", v...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.WarnLevel) {
		logObj.logger.Warnf(format, v...)
	}
}

func (logObj *Logger) Error(v ...interface{}) {
	if logObj == nil {
		fmt.Println(append([]interface{}{"[ERROR]"}, v...)...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.ErrorLevel) {
		logObj.logger.Error(v...)
	}
}

func (logObj *Logger) Errorf(format string, v ...interface{}) {
	if logObj == nil {
		fmt.Printf("[ERROR] "+format+"\n", v...)
		return
	}
	if logObj.logger.IsLevelEnabled(logrus.ErrorLevel) {
		logObj.logger.Errorf(format, v...)
	}
}

type textFormatter struct {
	maxLength int
}

func (f *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	if f.maxLength > 0 && len(entry.Message) > f.maxLength {
		entry.Message = entry.Message[:f.maxLength]
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	entry.Caller = getCaller()
	funcVal := entry.Caller.Function
	fileVal := fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)

	f.appendValue(b, strings.ToUpper(entry.Level.String()))
	b.WriteByte(' ')
	f.appendValue(b, entry.Time.Format(defaultFileWriterMsgSuffixTimeString))
	b.WriteByte(' ')

	f.appendValue(b, funcVal)
	b.WriteByte(' ')
	f.appendValue(b, fileVal)
	b.WriteByte(' ')

	f.appendValue(b, entry.Message)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *textFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(stringVal)
}

/**
 * toolbox/log toolbox/log/logrus_wrap github.com/sirupsen/logrus
 * 上述这三个包如果要打印 log，函数调用栈会不准确。
 * 理论上这仨库不需要打日志到业务日志文件中，打印到标准输出即可
 */
func getCaller() *runtime.Frame {
	pcs := make([]uintptr, 25)
	depth := runtime.Callers(1, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		if !strings.Contains(pkg, "toolbox/log") && !strings.Contains(pkg, "sirupsen/logrus") {
			return &f
		}
	}

	return nil
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
