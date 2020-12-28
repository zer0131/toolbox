package stat

import (
	"net/url"
	"time"
)

const (
	Redis    = "redis"
	Redigo   = "redigo"
	Mysql    = "mysql"
	MysqlX   = "mysqlx"
	MysqlORM = "mysqlorm"
	ESV5     = "esv5"
	Http     = "http"
)

func ClientStat(name string, start time.Time) {
	//pfc.Meter(name, 1)
	//pfc.Histogram(name, time.Since(start).Nanoseconds()/(1000*1000))
}

func GetRawPath(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	return u.Path
}
