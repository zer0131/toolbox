package httplib

import (
	"context"
	"net/http"

	"toolbox/ip"
	"toolbox/log"
)

// 补全log-id，如果不存在
func CheckLogIdMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logId := r.Header.Get("log-id")
		if logId == "" {
			logId = log.GenLogId()
			r.Header.Set("log-id", logId)
		}

		next.ServeHTTP(w, r)
	})
}

func CheckIp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context
		if logId := r.Header.Get("log-id"); logId != "" {
			ctx = log.NewContextWithSpecifyLogID(ctx, logId)
		}

		if !ip.CheckIp(ctx, r.RemoteAddr) {
			log.Warnf(ctx, "ip[%s] not allow", r.RemoteAddr)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("IP is not allow!"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AddStat(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//name := fmt.Sprintf("%s.%s", r.URL.Path, toolbox.LocalIP())

		// qps统计：60s上传一次，可以观察75% 95% 99%三个指标
		//pfc.Meter(name, 1)

		// 耗时统计
		//start := time.Now()
		//defer func() {
		//pfc.Histogram(name, time.Since(start).Nanoseconds()/(1000*1000))
		//}()

		next.ServeHTTP(w, r)
	})
}
