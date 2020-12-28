package nameresolution

import (
	"errors"
	"math/rand"
	"strings"
)

var _ Server = new(LIST)

type LIST struct {
	Addr []string
}

func NewLIST(name string) (*LIST, error) {
	list := new(LIST)
	list.Addr = strings.Split(strings.TrimLeft(name, `list://`), ",")
	return list, nil
}

func (L LIST) GetAllIpPort() ([]string, error) {
	return L.Addr, nil
}

func (L LIST) Get() (string, error) {
	if len(L.Addr) <= 0 {
		return "", errors.New("server list is empty")
	}
	return L.Addr[rand.Int()%len(L.Addr)], nil
}

func (L LIST) GetAllIpPortByCurrentIdc() ([]string, error) {
	return L.Addr, nil
}

func (L LIST) GetByTag(tag string) (string, error) {
	if len(L.Addr) <= 0 {
		return "", errors.New("server list is empty")
	}
	return L.Addr[rand.Int()%len(L.Addr)], nil
}

func (L LIST) GetAllIpPortByTag(tag string) ([]string, error) {
	return L.Addr, nil
}

func (L LIST) Include(addr string) bool {
	for _, v := range L.Addr {
		if v == addr {
			return true
		}
	}
	return false
}

func (L LIST) Close() {
}

func (L LIST) Block(addr string) {
}
