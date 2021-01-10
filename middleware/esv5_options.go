package middleware

import "time"

type esv5Options struct {
	addr string

	// http配置
	keepalive        time.Duration
	timeout          time.Duration
	idleTimeout      time.Duration
	connTimeout      time.Duration
	maxIdleConnCount int
}

var defaultESV5Options = esv5Options{
	keepalive:        30 * time.Second,
	timeout:          3 * time.Second,
	idleTimeout:      30 * time.Second,
	connTimeout:      300 * time.Millisecond,
	maxIdleConnCount: 100,
}

type Esv5OptionsFunc func(*esv5Options)

func ESV5WithAddr(s string) Esv5OptionsFunc {
	return func(o *esv5Options) {
		o.addr = s
	}
}

func ESV5WithKeepalive(d int64) Esv5OptionsFunc {
	return func(o *esv5Options) {
		o.keepalive = time.Duration(d) * time.Millisecond
	}
}

func ESV5WithTimeout(d int64) Esv5OptionsFunc {
	return func(o *esv5Options) {
		o.timeout = time.Duration(d) * time.Millisecond
	}
}

func ESV5WithIdleTimeout(d int64) Esv5OptionsFunc {
	return func(o *esv5Options) {
		o.idleTimeout = time.Duration(d) * time.Millisecond
	}
}

func ESV5WithConnTimeout(d int64) Esv5OptionsFunc {
	return func(o *esv5Options) {
		o.connTimeout = time.Duration(d) * time.Millisecond
	}
}

func ESV5WithMaxIdleConnCount(d int64) Esv5OptionsFunc {
	return func(o *esv5Options) {
		o.maxIdleConnCount = int(d)
	}
}
