package nameresolution

import (
	"errors"
	"net/url"
)

var (
	ErrTagNotFound error = errors.New("no tag founded in dns")
	ErrNoInstance  error = errors.New("no instance avaliable founded in server")
)

type Server interface {
	//获取所有ip+port
	GetAllIpPort() ([]string, error)
	//获取当前idc的一个ip+port
	Get() (string, error)
	//获取当前idc的所有ip+port
	GetAllIpPortByCurrentIdc() ([]string, error)
	//获取指定idc的一个ip+port
	GetByTag(string) (string, error)
	//获取指定idc的所有ip+port
	GetAllIpPortByTag(string) ([]string, error)

	Include(string) bool
	Close()
	//use for update ip list in the cache of Prov
	Block(string)
}

func NewResolver(addr string) (Server, error) {
	urlSchema, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	switch urlSchema.Scheme {
	case "domain":
		return NewDOMAIN(urlSchema.Host)
	case "list":
		return NewLIST(urlSchema.Host)
	default:
		return nil, errors.New("wrong service type")
	}
}
