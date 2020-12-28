package log

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	_ = InitV3("app", 1024)
	SetLogLevel("DEBUG")
	Debug(context.Background(), "logrus test", getRandomString(10))
	Debugf(context.Background(), "logrus test %s", getRandomString(10))
	Info(context.Background(), "logrus test", getRandomString(2048))
	Infof(context.Background(), "logrus test %s", getRandomString(2048))
	InfofArray(context.Background(), "logrus test %+v", []string{"aaa", "bbb", "cccc", "ddd", "eee"}, 3)
	Warn(context.Background(), "logrus test", getRandomString(10))
	Warnf(context.Background(), "logrus test %s", getRandomString(10))
	Error(context.Background(), "logrus test", getRandomString(10))
	Errorf(context.Background(), "logrus test %s", getRandomString(10))
	Errorf(context.Background(), "logrus test %s", getRandomString(2048))
}

func BenchmarkTestLog(b *testing.B) {
	_ = InitV3("app", 8196)
	SetLogLevel("DEBUG")
	fill := getRandomString(2048)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(context.Background(), fill)
	}
}

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func TestCustomLog(t *testing.T) {
	logObj, err := NewCustomLogger(WithProject("abc"), WithPath("./log"), WithExpireDay(7), WithMaxLength(100), WithLogLevel("WARN"))
	if err != nil {
		fmt.Printf("new custom logger err: %s", err)
		return
	}
	logObj.Debug(context.Background(), "logrus test", getRandomString(10000))
	logObj.Debugf(context.Background(), "logrus test %s", getRandomString(10))
	logObj.Info(context.Background(), "logrus test", getRandomString(2048))
	logObj.Infof(context.Background(), "logrus test %s", getRandomString(2048))
	logObj.InfofArray(context.Background(), "logrus test %+v", []string{"aaa", "bbb", "cccc", "ddd", "eee"}, 3)
	logObj.Warn(context.Background(), "logrus test", getRandomString(10))
	logObj.Warnf(context.Background(), "logrus test %s", getRandomString(10))
	logObj.Error(context.Background(), "logrus test", getRandomString(10))
	logObj.Errorf(context.Background(), "logrus test %s", getRandomString(10))
	logObj.Errorf(context.Background(), "logrus test %s", getRandomString(2048))
}
