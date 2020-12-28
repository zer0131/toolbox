package middleware

import (
	"time"

	"github.com/gomodule/redigo/redis"

	"toolbox"
	"toolbox/stat"
)

// 包装redis.Pool不让上层用户直接使用第三方库中的redis，
// 因为会造成用户code中也import上面的github路径，这样就
// 不能控制app的使用方式。
type RedigoPool struct {
	rp *redis.Pool
}
type DpRedigoPool interface {
	Get() redis.Conn
}

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

type redigoOptionsFunc func(*redigoOptions)

func RedigoAddr(s string) redigoOptionsFunc {
	return func(o *redigoOptions) {
		o.addr = s
	}
}

func RedigoMaxOpenConnCount(s int64) redigoOptionsFunc {
	return func(o *redigoOptions) {
		o.maxOpenConnCount = s
	}
}
func RedigoMaxIdleConnCount(s int64) redigoOptionsFunc {
	return func(o *redigoOptions) {
		o.maxIdleConnCount = s
	}
}

func RedigoIdleTimeout(s int64) redigoOptionsFunc {
	return func(o *redigoOptions) {
		o.idleTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedigoConnTimeout(s int64) redigoOptionsFunc {
	return func(o *redigoOptions) {
		o.connTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedigoReadTimeout(s int64) redigoOptionsFunc {
	return func(o *redigoOptions) {
		o.readTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedigoWriteTimeout(s int64) redigoOptionsFunc {
	return func(o *redigoOptions) {
		o.writeTimeout = time.Duration(s) * time.Millisecond
	}
}

func InitRedigo(opt ...redigoOptionsFunc) (DpRedigoPool, error) {
	opts := defaultRedigoOptions
	for _, o := range opt {
		o(&opts)
	}
	//foxns, err := ns.New(ns.WithService(opts.addr), ns.WithConnTimeout(opts.connTimeout))
	//if err != nil {
	//	return nil, err
	//}
	pool := redis.Pool{
		IdleTimeout: opts.idleTimeout,
		// PoolSize是连接池大小，和最大空闲连接数是不同的。maxIdle是保留idle conn的时间，为的是
		// 防止系统峰值性能。最大连接池是保护性质的，但可以临时使用这个。
		// 当系统有maxIdle个conn时，证明系统负载低，否则就会不断申请连接直到maxIdle
		MaxIdle:         int(opts.maxIdleConnCount),
		MaxActive:       int(opts.maxOpenConnCount),
		MaxConnLifetime: 500 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				opts.addr,
				redis.DialConnectTimeout(opts.connTimeout),
				redis.DialReadTimeout(opts.readTimeout),
				redis.DialWriteTimeout(opts.writeTimeout),
				//redis.DialNetDial(foxns.DialForRedigo),
			)
		},
	}
	return &RedigoPool{rp: &pool}, nil
}

func (redigo *RedigoPool) Get() redis.Conn {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redigo), startTime)
	return redigo.rp.Get()
}
