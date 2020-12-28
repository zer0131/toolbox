package disp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

func Test_respList_sort(t *testing.T) {
	var tests = []struct {
		rl respList
	}{
		{rl: respList{&resp{index: 3}, &resp{index: 2}, &resp{index: 1}}},
		{rl: respList{&resp{index: 2}, &resp{index: 3}, &resp{index: 1}}},
		{rl: respList{&resp{index: 1}, &resp{index: 3}, &resp{index: 2}}},
	}

	expect := respList{&resp{index: 1}, &resp{index: 2}, &resp{index: 3}}

	for i, tt := range tests {
		sort.Sort(tt.rl)
		if !reflect.DeepEqual(tt.rl, expect) {
			t.Errorf("index %d %+v", i, tt.rl)
		}
	}
}

type a struct{ name string }

type b struct{ name string }

// 普通有err的情况
func testHandle_v1(ctx context.Context, req interface{}) (interface{}, error) {
	switch req.(type) {
	case a:
		return "a", nil
	case b:
		return "b", errors.New("berr")
	}
	panic("")
}

func Test_Execute_v1(t *testing.T) {
	ar := a{name: "a"}
	br := b{name: "b"}
	rs, err := Execute(context.Background(), []interface{}{ar, br}, testHandle_v1, true)

	fmt.Println(rs, err)
}

// 有err，需要取消goroutine的情况
func testHandle_v2(ctx context.Context, req interface{}) (interface{}, error) {
	switch req.(type) {
	case a:
		select {
		case <-time.After(3 * time.Second):
			return "a", nil
		case <-ctx.Done():
			fmt.Println("aerr happened")
			return "a", errors.New("aerr")
		}
	case b:
		panic("bpanic")
		return "b", errors.New("berr")
	}
	panic("")
}

func Test_Execute_v2(t *testing.T) {
	ar := a{name: "a"}
	br := b{name: "b"}
	_, err := Execute(context.Background(), []interface{}{ar, br}, testHandle_v2, true)

	fmt.Println(err)
}
