package db

import (
	"ai-knowledge/internal/config"
	"ai-knowledge/internal/logger"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mysql 连接

var (
	GormHandler *gorm.DB
)

func InitDB(debug bool, cfg *config.DbConfig) {
	if cfg == nil {
		log.Panicln("db config is nil")
		return
	}
	var err error
	GormHandler, err = NewDbClient(debug, cfg)
	if err != nil {
		log.Panicln("db init error", err)
		return
	}
}

// NewDbClient 创建数据库连接
func NewDbClient(debug bool, cfg *config.DbConfig) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("The database configuration file can not be empty.")
	}
	// 拼接连接数据库字符串
	connStr := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.User,
		cfg.Password,
		cfg.Address,
		cfg.Port,
		cfg.DbName)
	// 连接数据库
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		DisableAutomaticPing: false,
		AllowGlobalUpdate:    false,
	})
	if err != nil {
		return nil, err
	}
	// 是否开启debug模式
	if debug {
		// debug 模式
		db = db.Debug()
	}
	// 连接池最大连接数
	_db, err := db.DB()
	if err != nil {
		return nil, err
	}
	_db.SetMaxIdleConns(cfg.MaxIdleConns)
	// 默认打开连接数
	_db.SetMaxOpenConns(cfg.MaxOpenConns)
	// 开启协程ping MySQL数据库查看连接状态
	go func() {
		for {
			// ping
			err = _db.Ping()
			if err != nil {
				logger.Logger.Errorw("mysql ping error", "err", err)
			}
			// 间隔30s ping一次
			time.Sleep(time.Second * 30)
		}
	}()
	// 防止我条件修改和删除操作
	db.AllowGlobalUpdate = false
	return db, err
}
