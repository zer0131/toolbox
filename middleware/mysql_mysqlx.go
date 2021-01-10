package middleware

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/stat"
	"time"
)

func InitMysqlX(db DpMysql) DpMysqlX {
	dbx := sqlx.NewDb(db.GetDB(), "mysql")
	return &MysqlX{dbx}
}

type MysqlX struct {
	*sqlx.DB
}
type DpMysqlX interface {
	DriverName() string
	MapperFunc(mf func(string) string)
	Rebind(query string) string
	Unsafe() *sqlx.DB
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
	MustBegin() *sqlx.Tx
	Beginx() (*sqlx.Tx, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	MustExec(query string, args ...interface{}) sql.Result
	Preparex(query string) (*sqlx.Stmt, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

func (mysqlx *MysqlX) DriverName() string {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.DriverName()
}

func (mysqlx *MysqlX) MapperFunc(mf func(string) string) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	mysqlx.DB.MapperFunc(mf)
}

func (mysqlx *MysqlX) Rebind(query string) string {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.Rebind(query)
}

func (mysqlx *MysqlX) Unsafe() *sqlx.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.Unsafe()
}

func (mysqlx *MysqlX) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.BindNamed(query, arg)
}

func (mysqlx *MysqlX) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.NamedQuery(query, arg)
}

func (mysqlx *MysqlX) NamedExec(query string, arg interface{}) (sql.Result, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.NamedExec(query, arg)
}

func (mysqlx *MysqlX) Select(dest interface{}, query string, args ...interface{}) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.Select(dest, query, args...)
}

func (mysqlx *MysqlX) Get(dest interface{}, query string, args ...interface{}) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.Get(dest, query, args...)
}

func (mysqlx *MysqlX) MustBegin() *sqlx.Tx {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.MustBegin()
}

func (mysqlx *MysqlX) Beginx() (*sqlx.Tx, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.Beginx()
}

func (mysqlx *MysqlX) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.Queryx(query, args...)
}

func (mysqlx *MysqlX) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.QueryRowx(query, args...)
}

func (mysqlx *MysqlX) MustExec(query string, args ...interface{}) sql.Result {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.MustExec(query, args...)
}

func (mysqlx *MysqlX) Preparex(query string) (*sqlx.Stmt, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.Preparex(query)
}

func (mysqlx *MysqlX) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlX), startTime)
	return mysqlx.DB.PrepareNamed(query)
}
