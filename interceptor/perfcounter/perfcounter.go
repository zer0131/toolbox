package perfcounter

import (
	"fmt"
	"net"
	"strings"

	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var localIp string

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// /helloworld.Greeter/SayHello
func metrix(method string) string {
	if localIp == "" {
		localIp = getLocalIP()
		if localIp == "" {
			localIp = "unknown"
		}
	}
	return fmt.Sprintf("%s.%s", strings.Replace(method[1:], "/", ".", -1), localIp)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		//name := metrix(info.FullMethod)

		// qps统计：60s上传一次，可以观察75% 95% 99%三个指标
		//pfc.Meter(name, 1)

		// 耗时统计
		//start := time.Now()
		//defer func() {
		//pfc.Histogram(name, time.Since(start).Nanoseconds()/(1000*1000))
		//}()

		return handler(ctx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		//name := metrix(info.FullMethod)

		// qps统计：60s上传一次，可以观察75% 95% 99%三个指标
		//pfc.Meter(name, 1)

		// 耗时统计
		//start := time.Now()
		//defer func() {
		//pfc.Histogram(name, time.Since(start).Nanoseconds()/(1000*1000))
		//}()

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = stream.Context()
		return handler(srv, wrapped)
	}
}
