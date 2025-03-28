package model

import (
	"errors"
	"go-admin/config"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/sohaha/zlsgo/zlog"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	g_db *gorm.DB //数据库
)

type BASE_MODEL struct {
	ID        uint           `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

type ClearDB struct {
	TableName    string
	CompareField string
	Interval     string
}

func OnInit() error {

	err := initDb()
	if err != nil {
		return err
	}

	// 缺省设置
	modelSetup := SysSetup{}
	err = modelSetup.NoFoundCreate("report_form_layout", "{}", "report form 布局")
	if err != nil {
		zlog.Error("model initDb", zap.Error(err))
		return err
	}

	err = modelSetup.NoFoundCreate("report_id_format", "{ \"pre\": \"a1a-\", \"date\": \"YY\", \"len\": 6 }", "报告ID的格式")
	if err != nil {
		zlog.Error("model initDb", zap.Error(err))
		return err
	}

	modelUser := SysUser{}
	err = modelUser.NoFoundCreate("admin", config.Config.Api.AdminPassword, "admin", true, "", false)
	if err != nil {
		zlog.Error("model initDb", zap.Error(err))
		return err
	}

	return nil
}

func initDb() error {

	dbtype := config.Config.System.DbType
	zlog.Debugf("model initDb %+v\n", config.Config.System)
	zlog.Debug("model initDb", "LoadConfig database : "+dbtype)

	db, err := getGorm(dbtype)
	if err != nil {
		zlog.Error("get gorm error: " + err.Error())
		return err
	}

	g_db = db

	registerTables()

	return nil
}

func getGorm(dbType string) (*gorm.DB, error) {
	switch dbType {
	case "mysql":
		return gormMysql()
		/*
			case "pgsql":
				return GormPgSql()
			case "oracle":
				return GormOracle()
			case "mssql":
				return GormMssql()
		*/
	case "sqlite":
		return gormSqlite()
	default:
		return gormMysql()
	}
}

func registerTables() {
	zlog.Debug("model registerTables enter")

	err := g_db.AutoMigrate(
		SysOperatRecord{},
		SysUser{},
		SysAuthInfo{},
		SysSetup{},
	)
	if err != nil {
		zlog.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	zlog.Debug("model registerTables leave")

}

type DBBASE interface {
	GetLogMode() string
}

var Gorm = new(_gorm)

type _gorm struct{}

// Config gorm 自定义配置
// Author [SliverHorn](https://github.com/SliverHorn)
func (g *_gorm) Config(prefix string, singular bool) *gorm.Config {
	gorm_config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	_default := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	var logMode DBBASE
	switch config.Config.System.DbType {
	case "mysql":
		logMode = &config.Config.Mysql
		/*
			case "pgsql":
				logMode = &config.Config.Pgsql
			case "oracle":
				logMode = &config.Config.Oracle
		*/
	default:
		logMode = &config.Config.Mysql
	}

	switch logMode.GetLogMode() {
	case "silent", "Silent":
		gorm_config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		gorm_config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		gorm_config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		gorm_config.Logger = _default.LogMode(logger.Info)
	default:
		gorm_config.Logger = _default.LogMode(logger.Info)
	}
	return gorm_config
}

// GormMysql 初始化Mysql数据库
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
func gormMysql() (*gorm.DB, error) {
	m := config.Config.Mysql
	if m.Dbname == "" {
		return nil, errors.New("not config mysql")
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil, err
	} else {
		db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)

		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db, err
	}
}

func gormSqlite() (*gorm.DB, error) {
	s := config.Config.Sqlite

	if s.Dbname == "" {
		return nil, errors.New("not config sqlite")
	}

	if db, err := gorm.Open(sqlite.Open(s.Dsn()), Gorm.Config(s.Prefix, s.Singular)); err != nil {
		return nil, errors.New("not config sqlite")
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)

		return db, nil
	}
}

type SearchInfo struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Page           int        `json:"page" form:"page"`         // 页码
	PageSize       int        `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword1       string     `json:"keyword1" form:"keyword1"` //关键字1
	Keyword2       string     `json:"keyword2" form:"keyword2"` //关键字2
}

type PageInfo struct {
	Page     int `json:"page" form:"page"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}

func ErrRecordNotFound(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}
