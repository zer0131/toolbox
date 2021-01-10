package middleware

//ToDo(zer):gorm需要全面升级
import (
	"context"
	"database/sql"
	"github.com/jinzhu/gorm"
	//"gorm.io/gorm"
	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/stat"
	"time"
)

func InitMysqlORM(opt ...MysqlOptionsFunc) (DpMysqlORM, error) {
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

func (mysqlorm *MysqlORM) SetCtx(ctx context.Context) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	newDb := mysqlorm.DB.New()
	newDb.SetLogger(LoggerCtx{ctx})
	return newDb
}
