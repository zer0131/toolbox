package nameresolution

import (
	"fmt"
	"testing"
)

func TestNewResolver(t *testing.T) {
	urlList := []string{
		"list://127.0.0.1:8888,10.188.188.12:2222",
		"list://127.0.0.1:8888",
		"foxns://pusher.inf.svc",
		"domain://www.baidu.com",
		"dns://www.shit.com",
	}
	for _, v := range urlList {
		server, err := NewResolver(v)
		if err != nil {
			t.Errorf("cannot resolver, url is %s", v)
		}
		fmt.Println(server.GetAllIpPort())
	}
}
