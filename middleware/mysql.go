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
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	"github.com/zer0131/toolbox"
	foxlog "github.com/zer0131/toolbox/log"
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

// support orm
type MysqlORM struct {
	*gorm.DB
}

// DpMysqlORM ...
type DpMysqlORM interface {
	// New clone a new db connection without search conditions
	New() *gorm.DB
	// Close close current db connection.  If database connection is not an io.Closer, returns an error.
	Close() error

	// 去掉DB()方法，因为MysqlORM导致有field和method冲突的问题
	// DB get `*sql.DB` from current connection
	// If the underlying database connection is not a *sql.DB, returns nil
	// DB() *sql.DB

	// 自定义方法，提供给mysqlx使用，也会暴露给rd，解决上面描述的问题
	GetDB() *sql.DB

	// CommonDB return the underlying `*sql.DB` or `*sql.Tx` instance, mainly intended to allow coexistence with legacy non-GORM code.
	CommonDB() gorm.SQLCommon
	// Dialect get dialect
	Dialect() gorm.Dialect
	// Callback return `Callbacks` container, you could add/change/delete callbacks with it
	//     db.Callback().Create().Register("update_created_at", updateCreated)
	// Refer https://jinzhu.github.io/gorm/development.html#callbacks
	Callback() *gorm.Callback
	// SetLogger replace default logger
	// SetLogger(log logger)
	// LogMode set log mode, `true` for detailed logs, `false` for no log, default, will only print error logs
	LogMode(enable bool) *gorm.DB
	// SetNowFuncOverride set the function to be used when creating a new timestamp
	SetNowFuncOverride(nowFuncOverride func() time.Time) *gorm.DB
	// BlockGlobalUpdate if true, generates an error on update/delete without where clause.
	// This is to prevent eventual error with empty objects updates/deletions
	BlockGlobalUpdate(enable bool) *gorm.DB
	// HasBlockGlobalUpdate return state of block
	HasBlockGlobalUpdate() bool
	// SingularTable use singular table by default
	SingularTable(enable bool)
	// NewScope create a scope for current operation
	NewScope(value interface{}) *gorm.Scope
	// QueryExpr returns the query as expr object
	// QueryExpr() *expr
	// SubQuery returns the query as sub query
	// SubQuery() *expr
	// Where return a new relation, filter records with given conditions, accepts `map`, `struct` or `string` as conditions, refer http://jinzhu.github.io/gorm/crud.html#query
	Where(query interface{}, args ...interface{}) *gorm.DB
	// Or filter records that match before conditions or this one, similar to `Where`
	Or(query interface{}, args ...interface{}) *gorm.DB
	// Not filter records that don't match current conditions, similar to `Where`
	Not(query interface{}, args ...interface{}) *gorm.DB
	// Limit specify the number of records to be retrieved
	Limit(limit interface{}) *gorm.DB
	// Offset specify the number of records to skip before starting to return the records
	Offset(offset interface{}) *gorm.DB
	// Order specify order when retrieve records from database, set reorder to `true` to overwrite defined conditions
	//     db.Order("name DESC")
	//     db.Order("name DESC", true) // reorder
	//     db.Order(gorm.Expr("name = ? DESC", "first")) // sql expression
	Order(value interface{}, reorder ...bool) *gorm.DB
	// Select specify fields that you want to retrieve from database when querying, by default, will select all fields;
	// When creating/updating, specify fields that you want to save to database
	Select(query interface{}, args ...interface{}) *gorm.DB
	// Omit specify fields that you want to ignore when saving to database for creating, updating
	Omit(columns ...string) *gorm.DB
	// Group specify the group method on the find
	Group(query string) *gorm.DB
	// Having specify HAVING conditions for GROUP BY
	Having(query interface{}, values ...interface{}) *gorm.DB
	// Joins specify Joins conditions
	//     db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&user)
	Joins(query string, args ...interface{}) *gorm.DB
	// Scopes pass current database connection to arguments `func(*gorm.DB) *gorm.DB`, which could be used to add conditions dynamically
	//     func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
	//         return db.Where("amount > ?", 1000)
	//     }
	//
	//     func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
	//         return func (db *gorm.DB) *gorm.DB {
	//             return db.Scopes(AmountGreaterThan1000).Where("status in (?)", status)
	//         }
	//     }
	//
	//     db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
	// Refer https://jinzhu.github.io/gorm/crud.html#scopes
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
	// Unscoped return all record including deleted record, refer Soft Delete https://jinzhu.github.io/gorm/crud.html#soft-delete
	Unscoped() *gorm.DB
	// Attrs initialize struct with argument if record not found with `FirstOrInit` https://jinzhu.github.io/gorm/crud.html#firstorinit or `FirstOrCreate` https://jinzhu.github.io/gorm/crud.html#firstorcreate
	Attrs(attrs ...interface{}) *gorm.DB
	// Assign assign result with argument regardless it is found or not with `FirstOrInit` https://jinzhu.github.io/gorm/crud.html#firstorinit or `FirstOrCreate` https://jinzhu.github.io/gorm/crud.html#firstorcreate
	Assign(attrs ...interface{}) *gorm.DB
	// First find first record that match given conditions, order by primary key
	First(out interface{}, where ...interface{}) *gorm.DB
	// Take return a record that match given conditions, the order will depend on the database implementation
	Take(out interface{}, where ...interface{}) *gorm.DB
	// Last find last record that match given conditions, order by primary key
	Last(out interface{}, where ...interface{}) *gorm.DB
	// Find find records that match given conditions
	Find(out interface{}, where ...interface{}) *gorm.DB
	//Preloads preloads relations, don`t touch out
	Preloads(out interface{}) *gorm.DB
	// Scan scan value to a struct
	Scan(dest interface{}) *gorm.DB
	// Row return `*sql.Row` with given conditions
	Row() *sql.Row
	// Rows return `*sql.Rows` with given conditions
	Rows() (*sql.Rows, error)
	// ScanRows scan `*sql.Rows` to give struct
	ScanRows(rows *sql.Rows, result interface{}) error
	// Pluck used to query single column from a model as a map
	//     var ages []int64
	//     db.Find(&users).Pluck("age", &ages)
	Pluck(column string, value interface{}) *gorm.DB
	// Count get how many records for a model
	Count(value interface{}) *gorm.DB
	// Related get related associations
	Related(value interface{}, foreignKeys ...string) *gorm.DB
	// FirstOrInit find first matched record or initialize a new one with given conditions (only works with struct, map conditions)
	// https://jinzhu.github.io/gorm/crud.html#firstorinit
	FirstOrInit(out interface{}, where ...interface{}) *gorm.DB
	// FirstOrCreate find first matched record or create a new one with given conditions (only works with struct, map conditions)
	// https://jinzhu.github.io/gorm/crud.html#firstorcreate
	FirstOrCreate(out interface{}, where ...interface{}) *gorm.DB
	// Update update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
	Update(attrs ...interface{}) *gorm.DB
	// Updates update attributes with callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
	Updates(values interface{}, ignoreProtectedAttrs ...bool) *gorm.DB
	// UpdateColumn update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
	UpdateColumn(attrs ...interface{}) *gorm.DB
	// UpdateColumns update attributes without callbacks, refer: https://jinzhu.github.io/gorm/crud.html#update
	UpdateColumns(values interface{}) *gorm.DB
	// Save update value in database, if the value doesn't have primary key, will insert it
	Save(value interface{}) *gorm.DB
	// Create insert the value into database
	Create(value interface{}) *gorm.DB
	// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
	Delete(value interface{}, where ...interface{}) *gorm.DB
	// Raw use raw sql as conditions, won't run it unless invoked by other methods
	//    db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
	Raw(sql string, values ...interface{}) *gorm.DB
	// Exec execute raw sql
	Exec(sql string, values ...interface{}) *gorm.DB
	// Model specify the model you would like to run db operations
	//    // update all users's name to `hello`
	//    db.Model(&User{}).Update("name", "hello")
	//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
	//    db.Model(&user).Update("name", "hello")
	Model(value interface{}) *gorm.DB
	// Table specify the table you would like to run db operations
	Table(name string) *gorm.DB
	// Debug start debug mode
	Debug() *gorm.DB
	// Transaction start a transaction as a block,
	// return error will rollback, otherwise to commit.
	Transaction(fc func(tx *gorm.DB) error) error
	// Begin begin a transaction
	Begin() *gorm.DB
	// BeginTx begins a transaction with options
	BeginTx(ctx context.Context, opts *sql.TxOptions) *gorm.DB
	// Commit commit a transaction
	Commit() *gorm.DB
	// Rollback rollback a transaction
	Rollback() *gorm.DB
	// RollbackUnlessCommitted rollback a transaction if it has not yet been
	// committed.
	RollbackUnlessCommitted() *gorm.DB
	// NewRecord check if value's primary key is blank
	NewRecord(value interface{}) bool
	// RecordNotFound check if returning ErrRecordNotFound error
	RecordNotFound() bool
	// CreateTable create table for models
	CreateTable(models ...interface{}) *gorm.DB
	// DropTable drop table for models
	DropTable(values ...interface{}) *gorm.DB
	// DropTableIfExists drop table if it is exist
	DropTableIfExists(values ...interface{}) *gorm.DB
	// HasTable check has table or not
	HasTable(value interface{}) bool
	// AutoMigrate run auto migration for given models, will only add missing fields, won't delete/change current data
	AutoMigrate(values ...interface{}) *gorm.DB
	// ModifyColumn modify column to type
	ModifyColumn(column string, typ string) *gorm.DB
	// DropColumn drop a column
	DropColumn(column string) *gorm.DB
	// AddIndex add index for columns with given name
	AddIndex(indexName string, columns ...string) *gorm.DB
	// AddUniqueIndex add unique index for columns with given name
	AddUniqueIndex(indexName string, columns ...string) *gorm.DB
	// RemoveIndex remove index with name
	RemoveIndex(indexName string) *gorm.DB
	// AddForeignKey Add foreign key to the given scope, e.g:
	//     db.Model(&User{}).AddForeignKey("city_id", "cities(id)", "RESTRICT", "RESTRICT")
	AddForeignKey(field string, dest string, onDelete string, onUpdate string) *gorm.DB
	// RemoveForeignKey Remove foreign key from the given scope, e.g:
	//     db.Model(&User{}).RemoveForeignKey("city_id", "cities(id)")
	RemoveForeignKey(field string, dest string) *gorm.DB
	// Association start `Association Mode` to handler relations things easir in that mode, refer: https://jinzhu.github.io/gorm/associations.html#association-mode
	Association(column string) *gorm.Association
	// Preload preload associations with given conditions
	//    db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
	Preload(column string, conditions ...interface{}) *gorm.DB
	// Set set setting by name, which could be used in callbacks, will clone a new db, and update its setting
	Set(name string, value interface{}) *gorm.DB
	// InstantSet instant set setting, will affect current db
	InstantSet(name string, value interface{}) *gorm.DB
	// Get get setting by name
	Get(name string) (value interface{}, ok bool)
	// SetJoinTableHandler set a model's join table handler for a relation
	SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface)
	// AddError add error to the db
	AddError(err error) error
	// GetErrors get happened errors from the db
	GetErrors() []error
	// print log with logid
	SetCtx(context.Context) *gorm.DB
}

func (mysqlorm *MysqlORM) New() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.New()
}

func (mysqlorm *MysqlORM) Close() error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Close()
}

func (mysqlorm *MysqlORM) GetDB() *sql.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.DB()
}

func (mysqlorm *MysqlORM) CommonDB() gorm.SQLCommon {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.CommonDB()
}

func (mysqlorm *MysqlORM) Dialect() gorm.Dialect {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Dialect()
}

func (mysqlorm *MysqlORM) Callback() *gorm.Callback {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Callback()
}

func (mysqlorm *MysqlORM) LogMode(enable bool) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.LogMode(enable)
}

func (mysqlorm *MysqlORM) SetNowFuncOverride(nowFuncOverride func() time.Time) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.SetNowFuncOverride(nowFuncOverride)
}

func (mysqlorm *MysqlORM) BlockGlobalUpdate(enable bool) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.BlockGlobalUpdate(enable)
}

func (mysqlorm *MysqlORM) HasBlockGlobalUpdate() bool {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.HasBlockGlobalUpdate()
}

func (mysqlorm *MysqlORM) SingularTable(enable bool) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	mysqlorm.DB.SingularTable(enable)
}

func (mysqlorm *MysqlORM) NewScope(value interface{}) *gorm.Scope {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.NewScope(value)
}

func (mysqlorm *MysqlORM) Where(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Where(query, args...)
}

func (mysqlorm *MysqlORM) Or(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Or(query, args...)
}

func (mysqlorm *MysqlORM) Not(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Not(query, args...)
}

func (mysqlorm *MysqlORM) Limit(limit interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Limit(limit)
}

func (mysqlorm *MysqlORM) Offset(offset interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Offset(offset)
}

func (mysqlorm *MysqlORM) Order(value interface{}, reorder ...bool) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Order(value, reorder...)
}

func (mysqlorm *MysqlORM) Select(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Select(query, args...)
}

func (mysqlorm *MysqlORM) Omit(columns ...string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Omit(columns...)
}

func (mysqlorm *MysqlORM) Group(query string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Group(query)
}

func (mysqlorm *MysqlORM) Having(query interface{}, values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Having(query, values...)
}

func (mysqlorm *MysqlORM) Joins(query string, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Joins(query, args...)
}

func (mysqlorm *MysqlORM) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Scopes(funcs...)
}

func (mysqlorm *MysqlORM) Unscoped() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Unscoped()
}

func (mysqlorm *MysqlORM) Attrs(attrs ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Attrs(attrs...)
}

func (mysqlorm *MysqlORM) Assign(attrs ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Assign(attrs...)
}

func (mysqlorm *MysqlORM) First(out interface{}, where ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.First(out, where...)
}

func (mysqlorm *MysqlORM) Take(out interface{}, where ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Take(out, where...)
}

func (mysqlorm *MysqlORM) Last(out interface{}, where ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Last(out, where...)
}

func (mysqlorm *MysqlORM) Find(out interface{}, where ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Find(out, where...)
}

func (mysqlorm *MysqlORM) Preloads(out interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Preloads(out)
}

func (mysqlorm *MysqlORM) Scan(dest interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Scan(dest)
}

func (mysqlorm *MysqlORM) Row() *sql.Row {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Row()
}

func (mysqlorm *MysqlORM) Rows() (*sql.Rows, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Rows()
}

func (mysqlorm *MysqlORM) ScanRows(rows *sql.Rows, result interface{}) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.ScanRows(rows, result)
}

func (mysqlorm *MysqlORM) Pluck(column string, value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Pluck(column, value)
}

func (mysqlorm *MysqlORM) Count(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Count(value)
}

func (mysqlorm *MysqlORM) Related(value interface{}, foreignKeys ...string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Related(value, foreignKeys...)
}

func (mysqlorm *MysqlORM) FirstOrInit(out interface{}, where ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.FirstOrInit(out, where...)
}

func (mysqlorm *MysqlORM) FirstOrCreate(out interface{}, where ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.FirstOrCreate(out, where...)
}

func (mysqlorm *MysqlORM) Update(attrs ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Update(attrs...)
}

func (mysqlorm *MysqlORM) Updates(values interface{}, ignoreProtectedAttrs ...bool) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Updates(values, ignoreProtectedAttrs...)
}

func (mysqlorm *MysqlORM) UpdateColumn(attrs ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.UpdateColumn(attrs...)
}

func (mysqlorm *MysqlORM) UpdateColumns(values interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.UpdateColumns(values)
}

func (mysqlorm *MysqlORM) Save(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Save(value)
}

func (mysqlorm *MysqlORM) Create(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Create(value)
}

func (mysqlorm *MysqlORM) Delete(value interface{}, where ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Delete(value, where...)
}

func (mysqlorm *MysqlORM) Raw(sql string, values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Raw(sql, values...)
}

func (mysqlorm *MysqlORM) Exec(sql string, values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Exec(sql, values...)
}

func (mysqlorm *MysqlORM) Model(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Model(value)
}

func (mysqlorm *MysqlORM) Table(name string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Table(name)
}

func (mysqlorm *MysqlORM) Debug() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Debug()
}

func (mysqlorm *MysqlORM) Transaction(fc func(tx *gorm.DB) error) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Transaction(fc)
}

func (mysqlorm *MysqlORM) Begin() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Begin()
}

func (mysqlorm *MysqlORM) BeginTx(ctx context.Context, opts *sql.TxOptions) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.BeginTx(ctx, opts)
}

func (mysqlorm *MysqlORM) Commit() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Commit()
}

func (mysqlorm *MysqlORM) Rollback() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Rollback()
}

func (mysqlorm *MysqlORM) RollbackUnlessCommitted() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.RollbackUnlessCommitted()
}

func (mysqlorm *MysqlORM) NewRecord(value interface{}) bool {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.NewRecord(value)
}

func (mysqlorm *MysqlORM) RecordNotFound() bool {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.RecordNotFound()
}

func (mysqlorm *MysqlORM) CreateTable(models ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.CreateTable(models...)
}

func (mysqlorm *MysqlORM) DropTable(values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.DropTable(values...)
}

func (mysqlorm *MysqlORM) DropTableIfExists(values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.DropTableIfExists(values...)
}

func (mysqlorm *MysqlORM) HasTable(value interface{}) bool {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.HasTable(value)
}

func (mysqlorm *MysqlORM) AutoMigrate(values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.AutoMigrate(values...)
}

func (mysqlorm *MysqlORM) ModifyColumn(column string, typ string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.ModifyColumn(column, typ)
}

func (mysqlorm *MysqlORM) DropColumn(column string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.DropColumn(column)
}

func (mysqlorm *MysqlORM) AddIndex(indexName string, columns ...string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.AddIndex(indexName, columns...)
}

func (mysqlorm *MysqlORM) AddUniqueIndex(indexName string, columns ...string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.AddUniqueIndex(indexName, columns...)
}

func (mysqlorm *MysqlORM) RemoveIndex(indexName string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.RemoveIndex(indexName)
}

func (mysqlorm *MysqlORM) AddForeignKey(field string, dest string, onDelete string, onUpdate string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.AddForeignKey(field, dest, onDelete, onUpdate)
}

func (mysqlorm *MysqlORM) RemoveForeignKey(field string, dest string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.RemoveForeignKey(field, dest)
}

func (mysqlorm *MysqlORM) Association(column string) *gorm.Association {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Association(column)
}

func (mysqlorm *MysqlORM) Preload(column string, conditions ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Preload(column, conditions...)
}

func (mysqlorm *MysqlORM) Set(name string, value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Set(name, value)
}

func (mysqlorm *MysqlORM) InstantSet(name string, value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.InstantSet(name, value)
}

func (mysqlorm *MysqlORM) Get(name string) (value interface{}, ok bool) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Get(name)
}

func (mysqlorm *MysqlORM) SetJoinTableHandler(source interface{}, column string, handler gorm.JoinTableHandlerInterface) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	mysqlorm.DB.SetJoinTableHandler(source, column, handler)
}

func (mysqlorm *MysqlORM) AddError(err error) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.AddError(err)
}

func (mysqlorm *MysqlORM) GetErrors() []error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.GetErrors()
}

type LoggerCtx struct {
	ctx context.Context
}

func (logger LoggerCtx) Print(values ...interface{}) {
	foxlog.Info(logger.ctx, logFormatter(values...)...)
}

func (mysqlorm *MysqlORM) SetCtx(ctx context.Context) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	newDb := mysqlorm.DB.New()
	newDb.SetLogger(LoggerCtx{ctx})
	return newDb
}

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
}

var defaultMysqlOptions = mysqlOptions{
	connTimeout:      300 * time.Millisecond,
	readTimeout:      1 * time.Second,
	writeTimeout:     1 * time.Second,
	connMaxLifetime:  500 * time.Second,
	maxIdleConnCount: 50,
	maxOpenConnCount: 100,
	parseTime:        true,
	ormLogMode:       false,
}

type mysqlOptionsFunc func(*mysqlOptions)

func MysqlAddr(s string) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.addr = s
	}
}

func MysqlUser(s string) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.user = s
	}
}

func MysqlPasswd(s string) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.passwd = s
	}
}

func MysqlDbname(s string) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.dbname = s
	}
}

func MysqlConnTimeout(s int64) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.connTimeout = time.Duration(s) * time.Millisecond
	}
}

func MysqlReadTimeout(s int64) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.readTimeout = time.Duration(s) * time.Millisecond
	}
}

func MysqlWriteTimeout(s int64) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.writeTimeout = time.Duration(s) * time.Millisecond
	}
}
func MysqlConnMaxLifetime(s int64) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.connMaxLifetime = time.Duration(s) * time.Millisecond
	}
}

func MysqlMaxIdleConnCount(s int64) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.maxIdleConnCount = s
	}
}

func MysqlMaxOpenConnCount(s int64) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.maxOpenConnCount = s
	}
}

func MysqlParseTime(s bool) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.parseTime = s
	}
}

func MysqlOrmLogMode(s bool) mysqlOptionsFunc {
	return func(o *mysqlOptions) {
		o.ormLogMode = s
	}
}

func assemblyConfigAndRegisterDial(opt ...mysqlOptionsFunc) (*mysql.Config, *mysqlOptions, error) {
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

	return dbCfg, &opts, nil
}

func InitMysql(opt ...mysqlOptionsFunc) (DpMysql, error) {
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

	mysql.SetLogger(mysql.Logger(log.New(os.Stdout, "[mysql] ", log.Ldate|log.Ltime|log.Lshortfile)))
	return &Mysql{db}, nil
}

func InitMysqlX(db DpMysql) DpMysqlX {
	dbx := sqlx.NewDb(db.GetDB(), "mysql")
	return &MysqlX{dbx}
}

func InitMysqlORM(opt ...mysqlOptionsFunc) (DpMysqlORM, error) {
	dbCfg, opts, err := assemblyConfigAndRegisterDial(opt...)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open("mysql", dbCfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxOpenConns(int(opts.maxOpenConnCount))
	db.DB().SetMaxIdleConns(int(opts.maxIdleConnCount))
	db.DB().SetConnMaxLifetime(opts.connMaxLifetime)
	db.LogMode(opts.ormLogMode)
	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	return &MysqlORM{db}, nil
}
