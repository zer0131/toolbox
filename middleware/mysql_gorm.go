package middleware

import (
	"context"
	"database/sql"
	"github.com/zer0131/toolbox"
	"github.com/zer0131/toolbox/stat"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	//"github.com/jinzhu/gorm"
	"gorm.io/gorm"
	"time"
)

func InitMysqlORM(opt ...MysqlOptionsFunc) (DpMysqlORM, error) {
	dbCfg, opts, err := assemblyConfigAndRegisterDial(opt...)
	if err != nil {
		return nil, err
	}

	//db, err := gorm.Open("mysql", dbCfg.FormatDSN())
	db, err := gorm.Open(mysql.Open(dbCfg.FormatDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetMaxOpenConns(int(opts.maxOpenConnCount))
	sqlDb.SetMaxIdleConns(int(opts.maxIdleConnCount))
	sqlDb.SetConnMaxLifetime(opts.connMaxLifetime)
	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}

	return &MysqlORM{opts, db}, nil
}

type MysqlORM struct {
	mysqlOptions
	*gorm.DB
}

// DpMysqlORM ...
type DpMysqlORM interface {

	// Session create new db session
	Session(config *gorm.Session) *gorm.DB

	// WithContext change current instance db's context to ctx
	WithContext(ctx context.Context) *gorm.DB

	// 去掉DB()方法，因为MysqlORM导致有field和method冲突的问题
	// DB get `*sql.DB` from current connection
	// If the underlying database connection is not a *sql.DB, returns nil
	// DB() *sql.DB

	// Debug start debug mode
	Debug() *gorm.DB

	// Set store value with key into current db instance's context
	Set(key string, value interface{}) *gorm.DB

	// Get get value with key from current db instance's context
	Get(key string) (interface{}, bool)

	// InstanceSet store value with key into current db instance's context
	InstanceSet(key string, value interface{}) *gorm.DB

	// InstanceGet get value with key from current db instance's context
	InstanceGet(key string) (interface{}, bool)

	// Create insert the value into database
	Create(value interface{}) *gorm.DB

	// CreateInBatches insert the value in batches into database
	CreateInBatches(value interface{}, batchSize int) *gorm.DB

	// Save update value in database, if the value doesn't have primary key, will insert it
	Save(value interface{}) *gorm.DB

	// First find first record that match given conditions, order by primary key
	First(dest interface{}, conds ...interface{}) *gorm.DB

	// Take return a record that match given conditions, the order will depend on the database implementation
	Take(dest interface{}, conds ...interface{}) *gorm.DB

	// Last find last record that match given conditions, order by primary key
	Last(dest interface{}, conds ...interface{}) *gorm.DB

	// Find find records that match given conditions
	Find(dest interface{}, conds ...interface{}) *gorm.DB

	// FindInBatches find records in batches
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB

	FirstOrInit(dest interface{}, conds ...interface{}) *gorm.DB

	FirstOrCreate(dest interface{}, conds ...interface{}) *gorm.DB

	// Update update attributes with callbacks, refer: https://gorm.io/docs/update.html#Update-Changed-Fields
	Update(column string, value interface{}) *gorm.DB

	// Updates update attributes with callbacks, refer: https://gorm.io/docs/update.html#Update-Changed-Fields
	Updates(values interface{}) *gorm.DB

	UpdateColumn(column string, value interface{}) *gorm.DB

	UpdateColumns(values interface{}) *gorm.DB

	// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
	Delete(value interface{}, conds ...interface{}) *gorm.DB

	Count(count *int64) *gorm.DB

	Row() *sql.Row

	Rows() (*sql.Rows, error)

	// Scan scan value to a struct
	Scan(dest interface{}) *gorm.DB

	// Pluck used to query single column from a model as a map
	//     var ages []int64
	//     db.Model(&users).Pluck("age", &ages)
	Pluck(column string, dest interface{}) *gorm.DB

	ScanRows(rows *sql.Rows, dest interface{}) error

	// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)

	// Begin begins a transaction
	Begin(opts ...*sql.TxOptions) *gorm.DB

	// Commit commit a transaction
	Commit() *gorm.DB

	// Rollback rollback a transaction
	Rollback() *gorm.DB

	SavePoint(name string) *gorm.DB

	RollbackTo(name string) *gorm.DB

	// Exec execute raw sql
	Exec(sql string, values ...interface{}) *gorm.DB

	// Model specify the model you would like to run db operations
	//    // update all users's name to `hello`
	//    db.Model(&User{}).Update("name", "hello")
	//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
	//    db.Model(&user).Update("name", "hello")
	Model(value interface{}) *gorm.DB

	// Clauses Add clauses
	Clauses(conds ...clause.Expression) *gorm.DB

	// Table specify the table you would like to run db operations
	Table(name string, args ...interface{}) *gorm.DB

	// Distinct specify distinct fields that you want querying
	Distinct(args ...interface{}) *gorm.DB

	// Select specify fields that you want when querying, creating, updating
	Select(query interface{}, args ...interface{}) *gorm.DB

	// Omit specify fields that you want to ignore when creating, updating and querying
	Omit(columns ...string) *gorm.DB

	// Where add conditions
	Where(query interface{}, args ...interface{}) *gorm.DB

	// Not add NOT conditions
	Not(query interface{}, args ...interface{}) *gorm.DB

	// Or add OR conditions
	Or(query interface{}, args ...interface{}) *gorm.DB

	// Joins specify Joins conditions
	//     db.Joins("Account").Find(&user)
	//     db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&user)
	//     db.Joins("Account", DB.Select("id").Where("user_id = users.id AND name = ?", "someName").Model(&Account{}))
	Joins(query string, args ...interface{}) *gorm.DB

	// Group specify the group method on the find
	Group(name string) *gorm.DB

	// Having specify HAVING conditions for GROUP BY
	Having(query interface{}, args ...interface{}) *gorm.DB

	// Order specify order when retrieve records from database
	//     db.Order("name DESC")
	//     db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
	Order(value interface{}) *gorm.DB

	// Limit specify the number of records to be retrieved
	Limit(limit int) *gorm.DB

	// Offset specify the number of records to skip before starting to return the records
	Offset(offset int) *gorm.DB

	// Scopes pass current database connection to arguments `func(DB) DB`, which could be used to add conditions dynamically
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
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB

	// Preload preload associations with given conditions
	//    db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
	Preload(query string, args ...interface{}) *gorm.DB

	Attrs(attrs ...interface{}) *gorm.DB

	Assign(attrs ...interface{}) *gorm.DB

	Unscoped() *gorm.DB

	Raw(sql string, values ...interface{}) *gorm.DB

	// 自定义方法，提供给mysqlx使用，也会暴露给rd，解决上面描述的问题
	GetDB() *sql.DB

	//设置上下文
	SetCtx(ctx context.Context) *gorm.DB
}

func (mysqlorm *MysqlORM) Session(config *gorm.Session) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Session(config)
}

func (mysqlorm *MysqlORM) WithContext(ctx context.Context) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.WithContext(ctx)
}

func (mysqlorm *MysqlORM) Debug() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Debug()
}

func (mysqlorm *MysqlORM) Set(key string, value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Set(key, value)
}

func (mysqlorm *MysqlORM) Get(key string) (interface{}, bool) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Get(key)
}

func (mysqlorm *MysqlORM) InstanceSet(key string, value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.InstanceSet(key, value)
}

func (mysqlorm *MysqlORM) InstanceGet(key string) (interface{}, bool) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.InstanceGet(key)
}

func (mysqlorm *MysqlORM) Create(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Create(value)
}

func (mysqlorm *MysqlORM) CreateInBatches(value interface{}, batchSize int) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.CreateInBatches(value, batchSize)
}

func (mysqlorm *MysqlORM) Save(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Save(value)
}

func (mysqlorm *MysqlORM) First(dest interface{}, conds ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.First(dest, conds...)
}

func (mysqlorm *MysqlORM) Take(dest interface{}, conds ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Take(dest, conds...)
}

func (mysqlorm *MysqlORM) Last(dest interface{}, conds ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Last(dest, conds...)
}

func (mysqlorm *MysqlORM) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Find(dest, conds...)
}

func (mysqlorm *MysqlORM) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.FindInBatches(dest, batchSize, fc)
}

func (mysqlorm *MysqlORM) FirstOrInit(dest interface{}, conds ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.FirstOrInit(dest, conds...)
}

func (mysqlorm *MysqlORM) FirstOrCreate(dest interface{}, conds ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.FirstOrCreate(dest, conds...)
}

func (mysqlorm *MysqlORM) Update(column string, value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Update(column, value)
}

func (mysqlorm *MysqlORM) Updates(values interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Updates(values)
}

func (mysqlorm *MysqlORM) UpdateColumn(column string, value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.UpdateColumn(column, value)
}

func (mysqlorm *MysqlORM) UpdateColumns(values interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.UpdateColumns(values)
}

func (mysqlorm *MysqlORM) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Delete(value, conds...)
}

func (mysqlorm *MysqlORM) Count(count *int64) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Count(count)
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

func (mysqlorm *MysqlORM) Scan(dest interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Scan(dest)
}

func (mysqlorm *MysqlORM) Pluck(column string, dest interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Pluck(column, dest)
}

func (mysqlorm *MysqlORM) ScanRows(rows *sql.Rows, dest interface{}) error {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.ScanRows(rows, dest)
}

func (mysqlorm *MysqlORM) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Transaction(fc, opts...)
}

func (mysqlorm *MysqlORM) Begin(opts ...*sql.TxOptions) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Begin(opts...)
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

func (mysqlorm *MysqlORM) SavePoint(name string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.SavePoint(name)
}

func (mysqlorm *MysqlORM) RollbackTo(name string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.RollbackTo(name)
}

func (mysqlorm *MysqlORM) Exec(sql string, values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Exec(sql, values)
}

func (mysqlorm *MysqlORM) Model(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Model(value)
}

func (mysqlorm *MysqlORM) Clauses(conds ...clause.Expression) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Clauses(conds...)
}

func (mysqlorm *MysqlORM) Table(name string, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Table(name, args...)
}

func (mysqlorm *MysqlORM) Distinct(args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Distinct(args...)
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

func (mysqlorm *MysqlORM) Where(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Where(query, args...)
}

func (mysqlorm *MysqlORM) Not(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Not(query, args...)
}

func (mysqlorm *MysqlORM) Or(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Or(query, args...)
}

func (mysqlorm *MysqlORM) Joins(query string, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Joins(query, args...)
}

func (mysqlorm *MysqlORM) Group(name string) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Group(name)
}

func (mysqlorm *MysqlORM) Having(query interface{}, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Having(query, args...)
}

func (mysqlorm *MysqlORM) Order(value interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Order(value)
}

func (mysqlorm *MysqlORM) Limit(limit int) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Limit(limit)
}

func (mysqlorm *MysqlORM) Offset(offset int) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Offset(offset)
}

func (mysqlorm *MysqlORM) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Scopes(funcs...)
}

func (mysqlorm *MysqlORM) Preload(query string, args ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Preload(query, args...)
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

func (mysqlorm *MysqlORM) Unscoped() *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Unscoped()
}

func (mysqlorm *MysqlORM) Raw(sql string, values ...interface{}) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	return mysqlorm.DB.Raw(sql, values...)
}

func (mysqlorm *MysqlORM) GetDB() *sql.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	sqlDb, _ := mysqlorm.DB.DB()
	return sqlDb
}

func (mysqlorm *MysqlORM) SetCtx(ctx context.Context) *gorm.DB {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.MysqlORM), startTime)
	newLogger := NewLoggerMe(logger.Config{
		SlowThreshold:             mysqlorm.slowThreshold, // 慢 SQL 阈值
		LogLevel:                  logger.Silent,          // 日志级别
		IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  true,                   // 彩色打印
	})
	db := mysqlorm.DB.Session(&gorm.Session{Context: ctx, Logger: newLogger})
	return db
}
