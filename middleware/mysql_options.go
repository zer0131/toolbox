package middleware

import "time"

type mysqlOptions struct {
	addr, user, passwd, dbname string
	connTimeout                time.Duration
	readTimeout                time.Duration
	writeTimeout               time.Duration
	connMaxLifetime            time.Duration
	maxIdleConnCount           int64
	maxOpenConnCount           int64
	parseTime                  bool
	ormLogMode                 bool
	slowThreshold              time.Duration
}

type MysqlOptionsFunc func(*mysqlOptions)

var defaultMysqlOptions = mysqlOptions{
	connTimeout:      300 * time.Millisecond,
	readTimeout:      1 * time.Second,
	writeTimeout:     1 * time.Second,
	connMaxLifetime:  500 * time.Second,
	maxIdleConnCount: 50,
	maxOpenConnCount: 100,
	parseTime:        true,
	ormLogMode:       false,
	slowThreshold:    time.Second * 10,
}

func MysqlAddr(s string) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.addr = s
	}
}

func MysqlUser(s string) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.user = s
	}
}

func MysqlPasswd(s string) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.passwd = s
	}
}

func MysqlDbname(s string) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.dbname = s
	}
}

func MysqlConnTimeout(s int64) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.connTimeout = time.Duration(s) * time.Millisecond
	}
}

func MysqlReadTimeout(s int64) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.readTimeout = time.Duration(s) * time.Millisecond
	}
}

func MysqlWriteTimeout(s int64) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.writeTimeout = time.Duration(s) * time.Millisecond
	}
}
func MysqlConnMaxLifetime(s int64) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.connMaxLifetime = time.Duration(s) * time.Millisecond
	}
}

func MysqlMaxIdleConnCount(s int64) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.maxIdleConnCount = s
	}
}

func MysqlMaxOpenConnCount(s int64) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.maxOpenConnCount = s
	}
}

func MysqlParseTime(s bool) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.parseTime = s
	}
}

func MysqlOrmLogMode(s bool) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.ormLogMode = s
	}
}

func MysqlSlowThreshold(s int64) MysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.slowThreshold = time.Duration(s) * time.Second
	}
}
