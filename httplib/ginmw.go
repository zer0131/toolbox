package httplib

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/zer0131/toolbox/ip"
	"github.com/zer0131/toolbox/log"

	"github.com/gin-gonic/gin"
)

// CheckLogIdMwForGin 补全logId
func CheckLogIdMwForGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		logId := c.Request.Header.Get("log-id")
		if logId == "" {
			logId = log.GenLogId()
			c.Request.Header.Set("log-id", logId)
		}
		c.Next()
	}
}

// CheckIpForGin 检查Ip
func CheckIpForGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context
		if logId := c.Request.Header.Get("log-id"); logId != "" {
			ctx = log.NewContextWithSpecifyLogID(ctx, logId)
		}
		if !ip.CheckIp(ctx, c.Request.RemoteAddr) {
			log.Warnf(ctx, "ip[%s] not allow", c.Request.RemoteAddr)
			c.Writer.WriteHeader(http.StatusUnauthorized)
			_, _ = c.Writer.Write([]byte("IP is not allow!"))
			return
		}
		c.Next()
	}
}

// AddStatForGin 增加指标
func AddStatForGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ToDo(zer):代码逻辑
	}
}

// RequestInfoForGin 记录每次请求的参数
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		logID := c.Request.Header.Get("Log-Id")
		statusCode := c.Writer.Status()
		clientAddr := strings.Split(c.Request.RemoteAddr, ":")
		clientIP, clientPort := "", ""
		if len(clientAddr) >= 2 {
			clientIP = clientAddr[0]
			clientPort = clientAddr[1]
		}
		// 日志输出格式
		log.Infof(context.Background(), "logID = [%s], status_code = [%d], latencyTime = [%v], clientIP = [%s], clientPort = [%s], reqMethod = [%s], reqUri = [%s].",
			logID,
			statusCode,
			latencyTime,
			clientIP,
			clientPort,
			reqMethod,
			reqUri,
		)
	}
}
