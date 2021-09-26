package middleware

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/stat"
)

var (
	errMysqlParam string = "%s illegal when init mysql"
)

type Mysql struct {
	*sql.DB
}

type DpMysql interface {
	// 自定义方法，提供给mysqlx使用，也会暴露给rd
	GetDB() *sql.DB

	// PingContext verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	PingContext(ctx context.Context) error
	// Ping verifies a connection to the database is still alive,
	// establishing a connection if necessary.
	Ping() error
	// Close closes the database, releasing any open resources.
	//
	// It is rare to Close a DB, as the DB handle is meant to be
	// long-lived and shared between many goroutines.
	Close() error
	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
	//
	// If n <= 0, no idle connections are retained.
	SetMaxIdleConns(n int)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	SetMaxOpenConns(n int)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are reused forever.
	SetConnMaxLifetime(d time.Duration)
	// Stats returns database statistics.
	Stats() sql.DBStats
	// PrepareContext creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	//
	// The provided context is used for the preparation of the statement, not for the
	// execution of the statement.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	Prepare(query string) (*sql.Stmt, error)
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	Exec(query string, args ...interface{}) (sql.Result, error)
	// QueryContext executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	Query(query string, args ...interface{}) (*sql.Rows, error)
	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRow(query string, args ...interface{}) *sql.Row
	// BeginTx starts a transaction.
	//
	// The provided context is used until the transaction is committed or rolled back.
	// If the context is canceled, the sql package will roll back
	// the transaction. Tx.Commit will return an error if the context provided to
	// BeginTx is canceled.
	//
	// The provided TxOptions is optional and may be nil if defaults should be used.
	// If a non-default isolation level is used that the driver doesn't support,
	// an error will be returned.
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	// Begin starts a transaction. The default isolation level is dependent on
	// the driver.
	Begin() (*sql.Tx, error)
	// Driver returns the database's underlying driver.
	Driver() driver.Driver
	// Conn returns a single connection by either opening a new connection
	// or returning an existing connection from the connection pool. Conn will
	// block until either a connection is returned or ctx is canceled.
	// Queries run on the same Conn will be run in the same database session.
	//
	// Every Conn must be returned to the database pool after use by
	// calling Conn.Close.
	Conn(ctx context.Context) (*sql.Conn, error)
}

func (mysql *Mysql) GetDB() *sql.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB
}

func (mysql *Mysql) PingContext(ctx context.Context) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.PingContext(ctx)
}

func (mysql *Mysql) Ping() error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Ping()
}

func (mysql *Mysql) Close() error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Close()
}

func (mysql *Mysql) SetMaxIdleConns(n int) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	mysql.DB.SetMaxIdleConns(n)
}

func (mysql *Mysql) SetMaxOpenConns(n int) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	mysql.DB.SetMaxOpenConns(n)
}

func (mysql *Mysql) SetConnMaxLifetime(d time.Duration) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	mysql.DB.SetConnMaxLifetime(d)
}

func (mysql *Mysql) Stats() sql.DBStats {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Stats()
}

func (mysql *Mysql) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.PrepareContext(ctx, query)
}

func (mysql *Mysql) Prepare(query string) (*sql.Stmt, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Prepare(query)
}

func (mysql *Mysql) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.ExecContext(ctx, query, args...)
}

func (mysql *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Exec(query, args...)
}

func (mysql *Mysql) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.QueryContext(ctx, query, args...)
}

func (mysql *Mysql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Query(query, args...)
}

func (mysql *Mysql) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.QueryRowContext(ctx, query, args...)
}

func (mysql *Mysql) QueryRow(query string, args ...interface{}) *sql.Row {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.QueryRow(query, args...)
}

func (mysql *Mysql) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.BeginTx(ctx, opts)
}

func (mysql *Mysql) Begin() (*sql.Tx, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Begin()
}

func (mysql *Mysql) Driver() driver.Driver {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Driver()
}

func (mysql *Mysql) Conn(ctx context.Context) (*sql.Conn, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Mysql), startTime)
	return mysql.DB.Conn(ctx)
}

func assemblyConfigAndRegisterDial(opt ...MysqlOptionsFunc) (*mysql.Config, mysqlOptions, error) {
	opts := defaultMysqlOptions
	for _, o := range opt {
		o(&opts)
	}

	// 调用RegisterDial，是在go-sql-driver/mysql全局的map中注册
	// 所以要通过名字区分不同连接生成的目的端地址集
	netName := fmt.Sprintf("tcp_%s", opts.addr)

	dbCfg := mysql.NewConfig()
	dbCfg.Net = netName
	dbCfg.Addr = opts.addr
	dbCfg.User = opts.user
	dbCfg.Passwd = opts.passwd
	dbCfg.DBName = opts.dbname
	dbCfg.Timeout = opts.connTimeout
	dbCfg.ReadTimeout = opts.readTimeout
	dbCfg.WriteTimeout = opts.writeTimeout
	dbCfg.ParseTime = opts.parseTime

	//foxns, err := ns.New(ns.WithService(opts.addr), ns.WithConnTimeout(opts.connTimeout))
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//mysql.RegisterDial(netName, foxns.DialForMysql)

	return dbCfg, opts, nil
}

func InitMysql(opt ...MysqlOptionsFunc) (DpMysql, error) {
	dbCfg, opts, err := assemblyConfigAndRegisterDial(opt...)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(int(opts.maxOpenConnCount))
	db.SetMaxIdleConns(int(opts.maxIdleConnCount))
	db.SetConnMaxLifetime(opts.connMaxLifetime)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	err = mysql.SetLogger(mysql.Logger(log.New(os.Stdout, "[mysql] ", log.Ldate|log.Ltime|log.Lshortfile)))
	if err != nil {
		return nil, err
	}
	return &Mysql{db}, nil
}
