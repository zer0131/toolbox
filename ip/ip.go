package ip

import (
	"bytes"
	"net"
	"net/http"
	"strings"

	"context"

	"toolbox/log"
)

// ip库需求来源于两个业务：
// 1. 人财物中心，所有接口只允许white_list限制内的ip访问
// 2. sds，单接口只允许内网访问，所以包含于上面的功能之内

// 用户配置的white_list中的元素有两种形态
// 1. 10.189.240.0-10.189.240.255 ip段
// 2. 10.189.241.35 具体ip
type ipRestrict struct {
	// 标记是否是ip段
	isSegment bool

	// ip段
	start, stop string

	// 具体ip
	specfic string
}

func (r *ipRestrict) hit(ctx context.Context, ipTarget string) bool {
	if r.isSegment {
		if bytes.Compare(net.ParseIP(ipTarget), net.ParseIP(r.start)) >= 0 && bytes.Compare(net.ParseIP(ipTarget), net.ParseIP(r.stop)) <= 0 {
			return true
		}
	} else {
		if r.specfic == ipTarget {
			return true
		}
	}
	return false
}

var (
	whiteList []*ipRestrict
)

func Init(ctx context.Context, whiteListStrs []string) {
	for _, whiteListStr := range whiteListStrs {
		ips := strings.Split(whiteListStr, "-")
		if len(ips) == 1 {
			if !validIp(ctx, ips[0]) {
				continue
			}

			whiteList = append(whiteList, &ipRestrict{specfic: ips[0]})
		} else if len(ips) == 2 {
			if !validIp(ctx, ips[0]) || !validIp(ctx, ips[1]) {
				continue
			}

			whiteList = append(whiteList, &ipRestrict{isSegment: true, start: ips[0], stop: ips[1]})
		} else {
			log.Errorf(ctx, "Illegal ip config %s", whiteListStr)
		}
	}
}

func validIp(ctx context.Context, ipStr string) bool {
	if pip := net.ParseIP(ipStr); pip != nil && pip.To4() != nil {
		return true
	}
	return false
}

//检查IP是否合法
func CheckIp(ctx context.Context, ipStr string) bool {
	if whiteList == nil || len(whiteList) == 0 {
		return true
	}

	if !validIp(ctx, ipStr) {
		log.Warnf(ctx, "Invalid ip %s", ipStr)
		return true
	}

	for _, restrict := range whiteList {
		if restrict.hit(ctx, ipStr) {
			return true
		}
	}
	return false
}

func GetClientIp(ctx context.Context, r *http.Request) (clientIp string) {
	if r == nil {
		return ""
	}

	defer func() {
		if clientIp != "" && strings.Contains(clientIp, ":") {
			clientIp = clientIp[:strings.Index(clientIp, ":")]
		}
	}()

	if r.Header.Get("HTTP_X_FORWARDED_FOR") != "" && r.Header.Get("HTTP_X_FORWARDED_FOR") != "unknown" {
		raw := r.Header.Get("HTTP_X_FORWARDED_FOR")
		ipArray := strings.Split(raw, ",")
		if len(ipArray) > 0 {
			clientIp = ipArray[0]
			return
		}
	}

	if r.Header.Get("HTTP_CLIENT_IP") != "" && r.Header.Get("HTTP_CLIENT_IP") != "unknown" {
		clientIp = r.Header.Get("HTTP_CLIENT_IP")
		return
	}

	clientIp = r.RemoteAddr
	return
}
