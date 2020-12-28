package httplib

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"
)

func Test_Header2IncomingContext(t *testing.T) {

	header := make(http.Header)
	header.Set("a1", "b1")
	header.Set("a2", "b2")

	ctx := Header2IncomingContext(header)

	ctx = context.WithValue(ctx, "a", "b")
	ctx = context.WithValue(ctx, "c", "d")
	ctx = context.WithValue(ctx, "a1", "b3")

	fmt.Println(ctx)

	md, _ := metadata.FromIncomingContext(ctx)

	fmt.Println(md.Get("a1"), md.Get("a2"))
}
