package middleware

import "time"

// redis库内部存在Options，这里根据最佳实践，筛选出需要
// app关注的，二次封装，虽然屏蔽细节，但经过打磨后应该会
// 节省app开发时间。
type redigoOptions struct {
	addr             string
	maxOpenConnCount int64
	maxIdleConnCount int64
	idleTimeout      time.Duration

	// 下面3个选项，在初始化时时放在Dial方法中用的，这里app不需要关注
	connTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

var defaultRedigoOptions = redigoOptions{
	maxOpenConnCount: 50,
	maxIdleConnCount: 50,
	idleTimeout:      240 * time.Second,

	connTimeout:  300 * time.Millisecond,
	readTimeout:  1 * time.Second,
	writeTimeout: 1 * time.Second,
}

type RedigoOptionsFunc func(*redigoOptions)

func RedigoAddr(s string) RedigoOptionsFunc {
	return func(o *redigoOptions) {
		o.addr = s
	}
}

func RedigoMaxOpenConnCount(s int64) RedigoOptionsFunc {
	return func(o *redigoOptions) {
		o.maxOpenConnCount = s
	}
}
func RedigoMaxIdleConnCount(s int64) RedigoOptionsFunc {
	return func(o *redigoOptions) {
		o.maxIdleConnCount = s
	}
}

func RedigoIdleTimeout(s int64) RedigoOptionsFunc {
	return func(o *redigoOptions) {
		o.idleTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedigoConnTimeout(s int64) RedigoOptionsFunc {
	return func(o *redigoOptions) {
		o.connTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedigoReadTimeout(s int64) RedigoOptionsFunc {
	return func(o *redigoOptions) {
		o.readTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedigoWriteTimeout(s int64) RedigoOptionsFunc {
	return func(o *redigoOptions) {
		o.writeTimeout = time.Duration(s) * time.Millisecond
	}
}
