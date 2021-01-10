package middleware

import "time"

// redis库内部存在Options，这里根据最佳实践，筛选出需要
// app关注的，二次封装，虽然屏蔽细节，但经过打磨后应该会
// 节省app开发时间。
type redisOptions struct {
	masterName       string
	addr             string
	maxOpenConnCount int64
	idleTimeout      time.Duration
	maxConnAge       time.Duration
	maxRetries       int64
	password         string
	method           string

	// 下面3个选项，在初始化时时放在Dial方法中用的，这里app不需要关注
	connTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

var defaultRedisOptions = redisOptions{
	masterName:       "mymaster",
	maxOpenConnCount: 50,
	idleTimeout:      240 * time.Second,
	maxRetries:       1,
	maxConnAge:       500 * time.Second,
	password:         "",
	method:           TYPE_PROXY,

	connTimeout:  300 * time.Millisecond,
	readTimeout:  1 * time.Second,
	writeTimeout: 1 * time.Second,
}

type RedisOptionsFunc func(*redisOptions)

func RedisAddr(s string) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.addr = s
	}
}
func RedisMethod(s string) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.method = s
	}
}

func RedisMaxOpenConnCount(s int64) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.maxOpenConnCount = s
	}
}

func RedisMasterName(s string) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.masterName = s
	}
}

func RedisPassword(s string) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.password = s
	}
}

func RedisIdleTimeout(s int64) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.idleTimeout = time.Duration(s) * time.Millisecond
	}
}
func RedisMaxConnAge(s int64) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.maxConnAge = time.Duration(s) * time.Millisecond
	}
}
func RedisMaxRetries(s int64) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.maxRetries = s
	}
}

func RedisConnTimeout(s int64) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.connTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedisReadTimeout(s int64) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.readTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedisWriteTimeout(s int64) RedisOptionsFunc {
	return func(o *redisOptions) {
		o.writeTimeout = time.Duration(s) * time.Millisecond
	}
}
