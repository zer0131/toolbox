package nameresolution

import "strings"

var _ Server = new(DOMAIN)

type DOMAIN struct {
	name string
}

func NewDOMAIN(name string) (*DOMAIN, error) {
	domain := new(DOMAIN)
	domain.name = strings.TrimLeft(name, `domain://`)
	return domain, nil
}

func (D DOMAIN) GetAllIpPort() ([]string, error) {
	return []string{D.name}, nil
}

func (D DOMAIN) Get() (string, error) {
	return D.name, nil
}

func (D DOMAIN) GetAllIpPortByCurrentIdc() ([]string, error) {
	return []string{D.name}, nil
}

func (D DOMAIN) GetByTag(addr string) (string, error) {
	return D.name, nil
}

func (D DOMAIN) GetAllIpPortByTag(addr string) ([]string, error) {
	return []string{D.name}, nil
}

func (D DOMAIN) Include(addr string) bool {
	return strings.Contains(D.name, addr)
}

func (D DOMAIN) Close() {
}

func (D DOMAIN) Block(addr string) {
}
