package disp

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sort"
	"sync"

	"golang.org/x/sync/errgroup"
)

type resp struct {
	index int
	data  interface{}
}

type respList []*resp

func (l respList) Len() int {
	return len(l)
}

func (l respList) Less(i, j int) bool {
	return l[i].index < l[j].index
}

func (l respList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// rd实现这个方法，区分param的不同类型分开处理
type HandleFunc func(ctx context.Context, req interface{}) (interface{}, error)

func Execute(ctx context.Context, reqL []interface{}, f HandleFunc, errStop bool) ([]interface{}, error) {
	var (
		g    *errgroup.Group
		gctx context.Context

		mu sync.Mutex
		rl respList
	)

	// 遇到错误立即返回，并尝试取消所有线程
	if errStop {
		g, gctx = errgroup.WithContext(ctx)
	} else {
		// 这里直接忽略cancel，如果直接使用ctx，ctx可能存在cancel方法，会干扰上游逻辑
		gctx, _ = context.WithCancel(ctx)

		g = &errgroup.Group{}
	}

	for i, p := range reqL {
		// https://golang.org/doc/faq#closures_and_goroutines
		i, p := i, p
		g.Go(func() (err error) {
			defer func() {
				if perr := recover(); perr != nil {
					var buf [2048]byte
					n := runtime.Stack(buf[:], false)
					err = fmt.Errorf("panic: %v %s", perr, errors.New(string(buf[:n])))
				}
			}()

			var r interface{}

			r, err = f(gctx, p)

			mu.Lock()
			rl = append(rl, &resp{index: i, data: r})
			mu.Unlock()

			return
		})
	}

	sort.Sort(rl)

	// 这里拿到的是第一个报错
	werr := g.Wait()
	var fr []interface{}
	for _, r := range rl {
		fr = append(fr, r.data)
	}
	return fr, werr
}
