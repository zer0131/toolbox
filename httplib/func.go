package httplib

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func Header2IncomingContext(header http.Header) context.Context {
	m := make(map[string]string, len(header))
	for k, _ := range header {
		m[k] = header.Get(k)
	}
	ctx := context.Background()
	return metadata.NewIncomingContext(ctx, metadata.New(m))
}
