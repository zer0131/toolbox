package accesslog

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/zer0131/logfox"
	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/log"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

var LogObj *logfox.Logger

func PrintLog(ctx context.Context, request *http.Request) {
	logId, _ := log.LogIdFromContext(ctx)
	remoteAddr := request.RemoteAddr
	requestUrl := request.RequestURI
	host := request.Host
	method := request.Method
	cookie := request.Header.Get("Cookie")
	userAgent := request.Header.Get("User-Agent")
	serverIP := toolbox.GetLocalIP()
	newctx := log.NewContextWithLogID(ctx)
	noticef(newctx, "logId[%s] remoteAddr=%s requestUrl=%s %s host=%s http_cookie=%s http_user_agent=%s server_addr=%s",
		logId, remoteAddr, method, requestUrl, host, cookie, userAgent, serverIP)
}

func selfInit() error {
	if LogObj != nil {
		return nil
	}
	var err error
	if LogObj, err = logfox.NewLogger(
		"./log",
		"access_log",
		logfox.DEFAULT_FILEWRITER_MAX_EXPIRE_DAY,
		logfox.DEFAULT_FILEWRITER_MSG_SUFFIX_TIME_STRING,
	); err != nil {
		return err
	}
	return nil
}

func noticef(ctx context.Context, format string, v ...interface{}) {
	if LogObj == nil {
		err := selfInit()
		if err != nil {
			log.Errorf(ctx, "err=%s", err.Error())
			return
		}
	}
	LogObj.Output(fmt.Sprintf(format, v...), logfox.NoticeLevel)
}

//http项目的accesslog日志
func WriteHTTPLog(writer io.Writer, params handlers.LogFormatterParams) {
	buf := buildCommonLogLine(params.Request, params.URL, params.TimeStamp, params.StatusCode, params.Size)
	//buf = append(buf, '\n')
	logId := params.Request.Header.Get(log.LogIDKey)
	ctx := log.NewContextWithSpecifyLogID(context.Background(), logId)
	if LogObj == nil {
		var err error
		LogObj, err = logfox.NewLogger(
			"./log",
			"access",
			logfox.DEFAULT_FILEWRITER_MAX_EXPIRE_DAY,
			logfox.DEFAULT_FILEWRITER_FILE_SUFFIX_TIME_STRING)
		if err != nil {
			log.Errorf(ctx, "err=%s", err.Error())
			return
		}
	}
	LogObj.Output(buf, logfox.NoticeLevel)
}
func buildCommonLogLine(req *http.Request, url url.URL, ts time.Time, status int, size int) string {
	username := "-"
	if url.User != nil {
		if name := url.User.Username(); name != "" {
			username = name
		}
	}

	host, port, err := net.SplitHostPort(req.RemoteAddr)

	if err != nil {
		host = req.RemoteAddr
	}

	uri := req.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}
	if uri == "" {
		uri = url.RequestURI()
	}
	format := fmt.Sprintf("log_id=[%s] remote_addr=[%s] remote_port=[%s]"+
		" remote_user=[%s] request_time=[%s] "+
		"request=[%s] status=[%v] "+
		"response_size=[%d] http_cookie=[%s] "+
		"http_user_agent=[%s] server_addr=[%s] "+
		"host=[%s] cost=[%dms]",
		req.Header.Get(log.LogIDKey), host, port, username, ts.Format("02/Jan/2006:15:04:05 -0700"),
		fmt.Sprintf("%s %v %s", req.Method, uri, req.Proto), status,
		size, req.Header.Get("Cookie"), req.Header.Get("User-Agent"),
		toolbox.LocalIP(), req.Host, time.Since(ts).Nanoseconds()/(1000*1000))
	return format
}
