package httplib

import (
	"fmt"
	"math/rand"
	"strings"
)

type Addr string

func (a Addr) String() string {
	return string(a)
}

func (a Addr) Parse() (string, error) {
	as := string(a)

	if strings.HasPrefix(as, "list://") {
		addrList := strings.Split(strings.TrimPrefix(as, "list://"), ",")
		return fmt.Sprintf("http://%s", addrList[rand.Intn(len(addrList))]), nil
	}

	return as, nil
}
