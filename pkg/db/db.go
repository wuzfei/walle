package db

import (
	"errors"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const Mysql = "mysql"
const Postgresql = "postgres"
const Sqlite3 = "sqlite3"

var ErrDB = errs.Class("DB")

type Config struct {
	Driver string `help:"数据库驱动" default:"sqlite3"`
	//Dsn string `help:"数据库连接"  default:"ljg:abcd123456@tcp(192.168.43.90:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"`
	//Dsn string `help:"数据库连接"  default:"host=myhost port=myport user=gorm dbname=gorm password=mypassword"`
	//Dsn string `help:"数据库连接"  default:"/tmp/gorm.db"`
	Dsn      string `help:"数据库连接"  default:"/Users/wuxin/worker/yema.dev/yema_dev.db"`
	LogLevel string `help:"数据库日志打印级别,默认为空,可选[error|warn|info]" devDefault:"" default:"warn"`
}

func (conf *Config) Dialector() (dial gorm.Dialector, err error) {
	switch conf.Driver {
	case Mysql:
		dial = mysql.New(mysql.Config{
			DSN:                       conf.Dsn,
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		})
	case Postgresql:
		dial = postgres.New(postgres.Config{
			DSN: conf.Dsn,
		})
	case Sqlite3:
		dial = sqlite.Open(conf.Dsn)
	default:
		return nil, errors.New("database url error")
	}
	return
}

func NewDB(cfg *Config, zapLog *zap.Logger) (*gorm.DB, error) {
	dail, err := cfg.Dialector()
	if err != nil {
		return nil, ErrDB.Wrap(err)
	}
	return gorm.Open(dail, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger:                                   getLogInterface(zapLog, cfg.LogLevel),
	})
}
