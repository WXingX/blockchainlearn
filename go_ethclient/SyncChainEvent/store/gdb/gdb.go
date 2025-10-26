package gdb

import (
	"errors"
	"fmt"
	"time"

	"SyncChainEvent/utils/timeutils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Config struct {
	// 用户
	User string `toml:"user" json:"user"`
	// 密码
	Password string `toml:"password" json:"password"`
	// 地址
	Host string `toml:"host" json:"host"`
	// 端口
	Port int `toml:"port" json:"port"`
	// 数据库
	Database string `toml:"database" json:"database"`
	// 最大空闲连接数
	MaxIdleConns int `toml:"max_idle_conns" mapstructure:"max_idle_conns" json:"max_idle_conns"`
	// 最大打开连接数
	MaxOpenConns int `toml:"max_open_conns" mapstructure:"max_open_conns" json:"max_open_conns"`
	// 连接复用时间
	MaxConnMaxLifetime int64 `toml:"max_conn_max_lifetime" mapstructure:"max_conn_max_lifetime" json:"max_conn_max_lifetime"`
	// 日志级别，枚举（info、warn、error和silent）
	LogLevel string `toml:"log_level" mapstructure:"log_level" json:"log_level"`
}

func (cfg *Config) CreateDatabase() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	// TODO
	// err = gdb.Exec("CREATE DATABASE IF NOT EXISTS gdb;").Error

	return err
}

// GetDataSource 获取GORM Data Source信息
func (cfg *Config) GetDataSource() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

// GetMysqlConfig 获取GORM MySQL相关配置
func (cfg *Config) GetMysqlConfig() mysql.Config {
	return mysql.Config{
		DSN:                       cfg.GetDataSource(),
		DefaultStringSize:         255,  // string类型字段默认长度
		DisableDatetimePrecision:  true, // 禁用datetime精度
		DontSupportRenameIndex:    true, // 禁用重命名索引
		DontSupportRenameColumn:   true, // 禁用重命名列名
		SkipInitializeWithVersion: true, // 禁用根据当前mysql版本自动配置
	}
}

// GetGormConfig 获取GORM相关配置
func (cfg *Config) GetGormConfig() *gorm.Config {
	gc := &gorm.Config{
		QueryFields: true, // 根据字段名称查询
		PrepareStmt: true, // 缓存预编译语句
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 数据表名单数
		},
		NowFunc: func() time.Time {
			return timeutils.Now() // 当前时间载入时区
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
	}

	logLevel := logger.Warn
	switch cfg.LogLevel {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	case "silent":
		logLevel = logger.Silent
	}

	// gc.Logger = logger.Default.LogMode(logLevel)
	gc.Logger = NewLogger(logLevel, 200*time.Millisecond) // 设置日志记录器

	return gc
}

// NewDB 新建gorm.DB对象
func NewDB(cfg *Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("gdb: illegal gdb configure")
	}

	db, err := gorm.Open(mysql.New(cfg.GetMysqlConfig()), cfg.GetGormConfig())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "gdb: open database connection err", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "gdb: get database instance err", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	if cfg.MaxConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.MaxConnMaxLifetime))
	}

	return db, nil
}

// MustNewDB 新建gorm.DB对象
func MustNewDB(c *Config) *gorm.DB {
	db, err := NewDB(c)
	if err != nil {
		panic(err)
	}

	return db
}
