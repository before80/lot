package db

import (
	"fmt"
	"github.com/before80/lot/cfg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Default.Db.User,
		cfg.Default.Db.Password,
		cfg.Default.Db.Host,
		cfg.Default.Db.Port,
		cfg.Default.Db.Database,
	)
	var gormConfig gorm.Config
	if cfg.Default.EnvIsLocal == 1 {
		gormConfig = gorm.Config{}
	} else {
		gormConfig = gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		}
	}
	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gormConfig)
	if err != nil {
		panic(fmt.Sprintf("连接数据库出现错误：%v\n", err))
	}
}
