package middleware

import (
	"github.com/gomodule/redigo/redis"
	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/stat"
	"time"
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

func InitRedigo(opt ...RedigoOptionsFunc) (DpRedigoPool, error) {
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
