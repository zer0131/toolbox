package ip

import (
	"context"
	"errors"
	"github.com/zer0131/toolbox/ip"
	"github.com/zer0131/toolbox/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		var remoteAddr string
		p, ok := peer.FromContext(ctx)
		if ok {
			remoteAddr = p.Addr.String()
		}
		if !ip.CheckIp(ctx, remoteAddr) {
			log.Warnf(ctx, "Unauthorized ip[%s]", remoteAddr)
			return nil, errors.New("Ip unauthorized")
		}

		return handler(ctx, req)
	}
}
