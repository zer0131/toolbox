package httplib

import "time"

type httpClientOptions struct {
	// addr是否传递会改变httplib的行为
	// 传递：当前传递server参数的方法都不工作
	// 不传递：与上面相反
	addr string

	keepalive        time.Duration
	timeout          time.Duration
	idleTimeout      time.Duration
	connTimeout      time.Duration
	maxIdleConnCount int // net/http支持一个client对接的目的端是无限的，这是整体的限制

	// 包装http.Client实现重试
	retry int
	proxy string
}

var defaultHttpClientOptions = httpClientOptions{
	keepalive:        30 * time.Second,
	timeout:          3 * time.Second,
	idleTimeout:      30 * time.Second,
	connTimeout:      300 * time.Millisecond,
	maxIdleConnCount: 100,
	retry:            1,
}

type HttpClientOptionsFunc func(*httpClientOptions)

func HttpWithAddr(d string) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.addr = d
	}
}

func HttpWithKeepalive(d int64) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.keepalive = time.Duration(d) * time.Millisecond
	}
}

func HttpWithTimeout(d int64) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.timeout = time.Duration(d) * time.Millisecond
	}
}

func HttpWithIdleTimeout(d int64) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.idleTimeout = time.Duration(d) * time.Millisecond
	}
}

func HttpWithConnTimeout(d int64) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.connTimeout = time.Duration(d) * time.Millisecond
	}
}

func HttpWithMaxIdleConnCount(d int64) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.maxIdleConnCount = int(d)
	}
}

func HttpWithRetry(d int64) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		if d > maxRetry {
			d = maxRetry
		}
		if d < 0 {
			d = 0
		}
		//加上默认次数一次
		d += defaultRetry

		o.retry = int(d)
	}
}

func HttpWithProxy(d string) HttpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.proxy = d
	}
}
