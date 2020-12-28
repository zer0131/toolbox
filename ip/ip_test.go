package ip

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func Test_validIp(t *testing.T) {
	var tests = []struct {
		ipStr        string
		expectResult bool
	}{
		{
			ipStr:        "foo",
			expectResult: false,
		},
		{
			ipStr:        "",
			expectResult: false,
		},
		{
			ipStr:        "localhost",
			expectResult: false,
		},

		{
			ipStr:        "10.188.0.24",
			expectResult: true,
		},
		{
			ipStr:        "127.0.0.1",
			expectResult: true,
		},
	}

	for i, tt := range tests {
		if actual := validIp(context.Background(), tt.ipStr); actual != tt.expectResult {
			t.Errorf("Index %d expect %+v actual %+v", i, tt.expectResult, actual)
			t.SkipNow()
		}
	}
}

func Test_Init(t *testing.T) {
	var tests = []struct {
		whiteListStrs []string
		expectResult  []*ipRestrict
	}{
		{
			whiteListStrs: []string{"10.188.0.24", "10.188.0.25", "10.188.0.24-10.188.0.25"},
			expectResult: []*ipRestrict{
				{
					isSegment: false,
					specfic:   "10.188.0.24",
				},
				{
					isSegment: false,
					specfic:   "10.188.0.25",
				},
				{
					isSegment: true,
					start:     "10.188.0.24",
					stop:      "10.188.0.25",
				},
			},
		},
	}

	for idx, tt := range tests {
		// 清空whiteList
		whiteList = make([]*ipRestrict, 0)

		Init(context.Background(), tt.whiteListStrs)
		if !reflect.DeepEqual(whiteList, tt.expectResult) {
			t.Errorf("Index %d expect %+v result %+v", idx, whiteList, tt.expectResult)
			t.SkipNow()
		}
	}
}

func Test_ipRestrict_hit(t *testing.T) {
	var tests = []struct {
		ipTarget     string
		rt           *ipRestrict
		expectResult bool
	}{
		{
			ipTarget: "10.188.0.24",
			rt: &ipRestrict{
				isSegment: false,
				specfic:   "10.188.0.24",
			},
			expectResult: true,
		},
		{
			ipTarget: "10.188.0.24",
			rt: &ipRestrict{
				isSegment: false,
				specfic:   "10.188.0.25",
			},
			expectResult: false,
		},

		{
			ipTarget: "10.188.0.24",
			rt: &ipRestrict{
				isSegment: true,
				start:     "10.188.0.0",
				stop:      "10.188.255.255",
			},
			expectResult: true,
		},
		{
			ipTarget: "10.188.0.24",
			rt: &ipRestrict{
				isSegment: true,
				start:     "10.188.0.25",
				stop:      "10.188.255.255",
			},
			expectResult: false,
		},
	}

	for idx, tt := range tests {
		if actual := tt.rt.hit(context.Background(), tt.ipTarget); actual != tt.expectResult {
			t.Errorf("Index %d expect %+v actual %+v", idx, tt.expectResult, actual)
			t.SkipNow()
		}
	}
}

func Test_GetClientIp(t *testing.T) {
	if GetClientIp(context.TODO(), nil) != "" {
		t.Error("Not empty")
		t.SkipNow()
	}

	r, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.Header.Set("HTTP_X_FORWARDED_FOR", "10.188.0.24:12345,bar")
	if GetClientIp(context.TODO(), r) != "10.188.0.24" {
		t.Error("Not 10.188.0.24")
		t.SkipNow()
	}

	r, _ = http.NewRequest(http.MethodGet, "/test", nil)
	r.Header.Set("HTTP_CLIENT_IP", "10.188.0.25:23456")
	if GetClientIp(context.TODO(), r) != "10.188.0.25" {
		t.Error("Not 10.188.0.25")
		t.SkipNow()
	}

	r, _ = http.NewRequest(http.MethodGet, "/test", nil)
	r.RemoteAddr = "10.188.0.26:34567"
	if GetClientIp(context.TODO(), r) != "10.188.0.26" {
		t.Error("Not 10.188.0.26")
		t.SkipNow()
	}
}
