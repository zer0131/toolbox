package header

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"math/rand"
	"strings"
	"time"
)

func createLogId() string {
	t := time.Now().UnixNano() / 1000000
	r := rand.Intn(10000)
	return fmt.Sprintf("%d%d", t, r)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		defer toolbox.CatchPanic(ctx, req)

		// ctx在grpc app中一层层传递，可能还要传递给其他grpc app
		// 所以之类直接将ctx初始化好
		m := make(map[string]string)

		// logId
		var logId string
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			// 优先看是否有Log-Id，外部访问
			v, ok := md[strings.ToLower(log.LogIDKey)]
			if ok {
				logId = v[0]
			}

			// 然后再看是否是内部grpc之间调用
			if logId == "" {
				v, ok = md[log.LogIDKey]
				if ok {
					logId = v[0]
				}
			}
		}
		if logId == "" {
			logId = createLogId()
		}
		m[log.LogIDKey] = logId

		// remoteAddr
		var remoteAddr string
		p, ok := peer.FromContext(ctx)
		if ok {
			m[log.RemoteAddrName] = p.Addr.String()
			remoteAddr = p.Addr.String()
		}
		outgoingMd := metadata.New(m)

		newCtx := metadata.NewOutgoingContext(ctx, outgoingMd)

		// 放到ctxValue中的原因是，util中的log库需要这个logid
		newCtx = context.WithValue(newCtx, log.LogIDKey, logId)
		newCtx = context.WithValue(newCtx, log.RemoteAddrName, remoteAddr)

		return handler(newCtx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		// ctx在grpc app中一层层传递，可能还要传递给其他grpc app
		// 所以之类直接将ctx初始化好
		m := make(map[string]string)

		// logId
		var logId string
		md, ok := metadata.FromIncomingContext(stream.Context())
		if ok {
			v, ok := md[log.LogIDKey]
			if ok {
				logId = v[0]
			}
		}
		if logId == "" {
			logId = createLogId()
		}
		m[log.LogIDKey] = logId

		// remoteAddr
		p, ok := peer.FromContext(stream.Context())
		if ok {
			m[log.RemoteAddrName] = p.Addr.String()
		}
		outgoingMd := metadata.New(m)

		newCtx := metadata.NewOutgoingContext(stream.Context(), outgoingMd)

		// 放到ctxValue中的原因是，util中的log库需要这个logid
		newCtx = context.WithValue(newCtx, log.LogIDKey, logId)

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		return handler(srv, wrapped)
	}
}
