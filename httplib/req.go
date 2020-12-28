// http client对象自动生成
// 对外公开一些参数，暂时认为不需要app关心
package httplib

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"toolbox"
	"toolbox/log"
	"toolbox/stat"
)

const (
	defaultRetry = 1
	maxRetry     = 3
)

type HttpClient struct {
	hc *http.Client

	addr string

	retry int
}
type DpHttpClient interface {
	Get(ctx context.Context, api string) ([]byte, error)
	Post(ctx context.Context, api string, data []byte) ([]byte, error)
	PostForm(ctx context.Context, api string, values url.Values) ([]byte, error)

	// 策略网关http client和具体backend无关
	RawPost(ctx context.Context, rawUrl string, data []byte) ([]byte, error)

	// url的host变化，支持大网对接kafka的需求
	RawPostForm(ctx context.Context, rawUrl string, values url.Values) ([]byte, error)

	// 对接图片识别服务
	// https://studygolang.com/articles/5171
	PostFile(ctx context.Context, api string, header map[string]string, data *bytes.Buffer) ([]byte, error)

	//可设置header的post
	PostWithHeader(ctx context.Context, api string, header map[string]string, data []byte) ([]byte, error)

	Delete(ctx context.Context, api string) ([]byte, error)
}

// FIXME proxy需要解析功能
type Addr string

func (a Addr) String() string {
	return string(a)
}

func (a Addr) Parse() (string, error) {
	as := string(a)

	if strings.HasPrefix(as, "list://") {
		addrList := strings.Split(strings.TrimPrefix(as, "list://"), ",")
		return fmt.Sprintf("http://%s", addrList[rand.Intn(len(addrList))]), nil
	}

	return as, nil
}

func (c *HttpClient) Get(ctx context.Context, api string) ([]byte, error) {
	if c.addr == "" {
		return nil, errors.New("Empty addr")
	}

	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(fmt.Sprintf("%s.%s", stat.Http, stat.GetRawPath(api))), startTime)

	ap, err := Addr(c.addr).Parse()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s", ap, api)

	var (
		resp     *http.Response
		reqErr   error
		costTime float32
	)

	logId, _ := log.LogIdFromContext(ctx)
	startReqTime := time.Now()
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodGet, urlStr, nil)
		reqErr = err
		if err != nil {
			continue
		}
		req.Header.Add("log-id", logId)
		req.Header.Add("remote-addr", toolbox.LocalIP())

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	costTime = float32(time.Now().UnixNano()-startReqTime.UnixNano()) / 1e9
	defer printReqLog(ctx, resp, reqErr, urlStr, costTime)
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status not 200, resp=%+v", *resp)
	}
	return r, nil
}

func (c *HttpClient) Post(ctx context.Context, api string, data []byte) ([]byte, error) {
	if c.addr == "" {
		return nil, errors.New("Empty addr")
	}

	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(fmt.Sprintf("%s.%s", stat.Http, stat.GetRawPath(api))), startTime)

	ap, err := Addr(c.addr).Parse()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s", ap, api)

	var (
		resp     *http.Response
		reqErr   error
		costTime float32
	)

	logId, _ := log.LogIdFromContext(ctx)
	startReqTime := time.Now()
	// bug: 每次retry都要用新的req，否则出现：Post http://10.188.40.13:8988/multi: http: ContentLength=249 with Body length 0
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(data))
		reqErr = err
		if err != nil {
			continue
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("log-id", logId)
		req.Header.Add("remote-addr", toolbox.LocalIP())

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	costTime = float32(time.Now().UnixNano()-startReqTime.UnixNano()) / 1e9
	defer printReqLog(ctx, resp, reqErr, urlStr, costTime)
	defer resp.Body.Close()
	// bug: 每个resp.Body都要被读取空，否则不能建立长连接
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status not 200, resp=%+v", *resp)
	}
	return r, nil
}

func (c *HttpClient) RawPost(ctx context.Context, rawUrl string, data []byte) ([]byte, error) {
	var (
		resp     *http.Response
		reqErr   error
		costTime float32
	)

	logId, _ := log.LogIdFromContext(ctx)
	startReqTime := time.Now()
	// bug: 每次retry都要用新的req，否则出现：Post http://10.188.40.13:8988/multi: http: ContentLength=249 with Body length 0
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodPost, rawUrl, bytes.NewBuffer(data))
		reqErr = err
		if err != nil {
			continue
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("log-id", logId)
		req.Header.Add("remote-addr", toolbox.LocalIP())

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	costTime = float32(time.Now().UnixNano()-startReqTime.UnixNano()) / 1e9
	defer printReqLog(ctx, resp, reqErr, rawUrl, costTime)
	defer resp.Body.Close()
	// bug: 每个resp.Body都要被读取空，否则不能建立长连接
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status not 200, resp=%+v", *resp)
	}
	return r, nil
}

func (c *HttpClient) RawPostForm(ctx context.Context, rawUrl string, values url.Values) ([]byte, error) {
	var (
		resp   *http.Response
		reqErr error
	)

	logId, _ := log.LogIdFromContext(ctx)
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodPost, rawUrl, strings.NewReader(values.Encode()))
		reqErr = err
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("log-id", logId)
		req.Header.Add("remote-addr", toolbox.LocalIP())

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status not 200, resp=%+v", *resp)
	}
	return r, nil
}

func (c *HttpClient) PostForm(ctx context.Context, api string, values url.Values) ([]byte, error) {
	if c.addr == "" {
		return nil, errors.New("Empty addr")
	}

	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(fmt.Sprintf("%s.%s", stat.Http, stat.GetRawPath(api))), startTime)

	ap, err := Addr(c.addr).Parse()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s", ap, api)

	var (
		resp     *http.Response
		reqErr   error
		costTime float32
	)

	logId, _ := log.LogIdFromContext(ctx)
	startReqTime := time.Now()
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(values.Encode()))
		reqErr = err
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("log-id", logId)
		req.Header.Add("remote-addr", toolbox.LocalIP())

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	costTime = float32(time.Now().UnixNano()-startReqTime.UnixNano()) / 1e9
	defer printReqLog(ctx, resp, reqErr, urlStr, costTime)
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status not 200, resp=%+v", *resp)
	}
	return r, nil
}

func (c *HttpClient) PostFile(ctx context.Context, api string, header map[string]string, data *bytes.Buffer) ([]byte, error) {
	if c.addr == "" {
		return nil, errors.New("Empty addr")
	}

	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(fmt.Sprintf("%s.%s", stat.Http, stat.GetRawPath(api))), startTime)

	ap, err := Addr(c.addr).Parse()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s", ap, api)

	var (
		resp     *http.Response
		reqErr   error
		costTime float32
	)

	logId, _ := log.LogIdFromContext(ctx)
	startReqTime := time.Now()
	// bug: 每次retry都要用新的req，否则出现：Post http://10.188.40.13:8988/multi: http: ContentLength=249 with Body length 0
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodPost, urlStr, data)
		reqErr = err
		if err != nil {
			continue
		}
		for k, v := range header {
			req.Header.Set(k, v)
		}
		req.Header.Add("log-id", logId)

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	costTime = float32(time.Now().UnixNano()-startReqTime.UnixNano()) / 1e9
	defer printReqLog(ctx, resp, reqErr, urlStr, costTime)
	defer resp.Body.Close()
	// bug: 每个resp.Body都要被读取空，否则不能建立长连接
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status not 200, resp=%+v", *resp)
	}
	return r, nil
}

func (c *HttpClient) PostWithHeader(ctx context.Context, api string, header map[string]string, data []byte) ([]byte, error) {
	if c.addr == "" {
		return nil, errors.New("Empty addr")
	}

	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(fmt.Sprintf("%s.%s", stat.Http, stat.GetRawPath(api))), startTime)

	ap, err := Addr(c.addr).Parse()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s", ap, api)

	var (
		resp     *http.Response
		reqErr   error
		costTime float32
	)

	logId, _ := log.LogIdFromContext(ctx)
	startReqTime := time.Now()
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(data))
		reqErr = err
		if err != nil {
			continue
		}
		for k, v := range header {
			req.Header.Set(k, v)
		}
		req.Header.Add("log-id", logId)
		req.Header.Add("remote-addr", toolbox.LocalIP())

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	costTime = float32(time.Now().UnixNano()-startReqTime.UnixNano()) / 1e9
	defer printReqLog(ctx, resp, reqErr, urlStr, costTime)
	defer resp.Body.Close()
	// bug: 每个resp.Body都要被读取空，否则不能建立长连接
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return r, fmt.Errorf("status not 200, resp=%+v", *resp)
	}
	return r, nil
}

func (c *HttpClient) Delete(ctx context.Context, api string) ([]byte, error) {
	if c.addr == "" {
		return nil, errors.New("Empty addr")
	}

	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(fmt.Sprintf("%s.%s", stat.Http, stat.GetRawPath(api))), startTime)

	ap, err := Addr(c.addr).Parse()
	if err != nil {
		return nil, err
	}
	urlStr := fmt.Sprintf("%s%s", ap, api)

	var (
		resp     *http.Response
		reqErr   error
		costTime float32
	)

	logId, _ := log.LogIdFromContext(ctx)
	startReqTime := time.Now()
	for i := 0; i < c.retry; i++ {
		req, err := http.NewRequest(http.MethodDelete, urlStr, nil)
		reqErr = err
		if err != nil {
			continue
		}
		req.Header.Add("log-id", logId)
		req.Header.Add("remote-addr", toolbox.LocalIP())

		resp, err = c.hc.Do(req)
		reqErr = err
		if err != nil {
			continue
		}

		break
	}
	if reqErr != nil {
		return nil, reqErr
	}
	costTime = float32(time.Now().UnixNano()-startReqTime.UnixNano()) / 1e9
	defer printReqLog(ctx, resp, reqErr, urlStr, costTime)
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return r, fmt.Errorf("status not 200 or 204, resp=%+v", *resp)
	}
	return r, nil
}

func printReqLog(ctx context.Context, response *http.Response, err error, urlStr string, cost float32) {
	if response == nil {
		log.Infof(ctx, "err=%s url=%s", err, urlStr)
		return
	}
	request := response.Request
	protocol := request.Proto
	requestUrl := request.URL
	method := request.Method
	statusCode := response.StatusCode
	log.Infof(ctx, "proto=%s requestUrl=%s method=%s err=%+v realUrl=%s statusCode=%d cost=%f",
		protocol, requestUrl, method, err, urlStr, statusCode, cost)
}

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

type httpClientOptionsFunc func(*httpClientOptions)

func HttpWithAddr(d string) httpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.addr = d
	}
}

func HttpWithKeepalive(d int64) httpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.keepalive = time.Duration(d) * time.Millisecond
	}
}

func HttpWithTimeout(d int64) httpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.timeout = time.Duration(d) * time.Millisecond
	}
}

func HttpWithIdleTimeout(d int64) httpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.idleTimeout = time.Duration(d) * time.Millisecond
	}
}

func HttpWithConnTimeout(d int64) httpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.connTimeout = time.Duration(d) * time.Millisecond
	}
}

func HttpWithMaxIdleConnCount(d int64) httpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.maxIdleConnCount = int(d)
	}
}

func HttpWithRetry(d int64) httpClientOptionsFunc {
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

func HttpWithProxy(d string) httpClientOptionsFunc {
	return func(o *httpClientOptions) {
		o.proxy = d
	}
}

func InitHttpClient(opt ...httpClientOptionsFunc) (DpHttpClient, error) {
	opts := defaultHttpClientOptions
	for _, o := range opt {
		o(&opts)
	}

	// if opts.addr == "" {
	// 	return nil, errors.New("addr err")
	// }

	var proxy func(*http.Request) (*url.URL, error)
	if opts.proxy == "" {
		proxy = http.ProxyFromEnvironment
	} else {
		proxy = func(_ *http.Request) (*url.URL, error) {
			addr := Addr(opts.proxy)
			ap, err := addr.Parse()
			if err != nil {
				return nil, err
			}
			return url.Parse(ap)
		}
	}

	httpDialContextFunc := (&net.Dialer{Timeout: opts.connTimeout, KeepAlive: opts.keepalive, DualStack: true}).DialContext

	c := &HttpClient{
		addr: opts.addr,
		hc: &http.Client{
			Transport: &http.Transport{
				Proxy:       proxy,
				DialContext: httpDialContextFunc,

				IdleConnTimeout:       opts.idleTimeout,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 0,

				MaxIdleConns:        opts.maxIdleConnCount,
				MaxIdleConnsPerHost: opts.maxIdleConnCount,
			},
			Timeout: opts.timeout,
		},
		retry: opts.retry,
	}
	return c, nil
}
