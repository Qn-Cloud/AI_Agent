package common

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Config struct {
	Username                  string // 用户名
	Password                  string // 密码
	Host                      string // 链接地址
	Port                      string // 端口
	Database                  string // 数据库名
	Charset                   string // 字符集
	ParseTime                 string // 是否解析时间
	Loc                       string // 时区默认 Local
	SingularTable             bool   // 是否使用单数表名 true 是 false 否
	PrepareStmt               bool   //  在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
	SkipDefaultTransaction    bool
	ConnMaxLifetime           int64 // 设置了连接可复用的最大时间
	MaxIdleConn               int   // 连接池里面的连接最大空闲连接数
	ConnMaxIdleTime           int64 //连接池里面的连接最大空闲时长(秒)
	MaxOpenConn               int   // 设置打开数据库连接的最大数量
	IgnoreRecordNotFoundError bool  // 忽略ErrRecordNotFound（记录未找到）错误
}

// GetDB 获取数据库连接实例
func GetDB(c Config) (DB *gorm.DB) {
	var err error

	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Loc)

	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: c.SingularTable, //  使用单数表名，启用该选项，禁用表名复数
		},
		PrepareStmt:            c.PrepareStmt,
		QueryFields:            true, //打印sql
		SkipDefaultTransaction: c.SkipDefaultTransaction,
	})
	if err != nil {
		panic(err.Error())
	}

	err = DB.Use(tracing.NewPlugin())
	if err != nil {
		panic(err)
	}

	sqlDb, errMsg := DB.DB()
	if errMsg != nil {
		panic(errMsg.Error())
	}
	sqlDb.SetMaxIdleConns(c.MaxIdleConn)                                     //设置最大空闲连接数
	sqlDb.SetMaxOpenConns(c.MaxOpenConn)                                     //设置最大的空闲连接数
	sqlDb.SetConnMaxLifetime(time.Minute * time.Duration(c.ConnMaxLifetime)) //设置最大连接时长
	sqlDb.SetConnMaxIdleTime(time.Minute * time.Duration(c.ConnMaxIdleTime)) //设置最大的空闲连接时长
	data, _ := json.Marshal(sqlDb.Stats())                                   //获得当前的SQL配置情况
	fmt.Println(string(data))
	return
}
