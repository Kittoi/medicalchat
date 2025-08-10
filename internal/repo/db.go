package repo

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"medicalchat/pkg/utils"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	config := utils.GetConfig()

	// 构建 DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.PGVector.Host,
		config.PGVector.User,
		config.PGVector.Password,
		config.PGVector.DBName,
		config.PGVector.Port,
		config.PGVector.SSLMode,
	)

	log.Info().Str("dsn", dsn).Msg("数据库连接信息")

	// 配置 GORM 日志
	newLogger := logger.New(
		&GormZerologAdapter{}, // 使用自定义的 zerolog 适配器
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// 打开数据库连接
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 获取底层的 sqlDB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取 sqlDB 失败: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info().Int("maxIdleConns", 10).Int("maxOpenConns", 100).Dur("connMaxLifetime", time.Hour).Msg("设置数据库连接池参数")

	// 启用 pgvector 扩展
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS vector").Error; err != nil {
		return fmt.Errorf("启用 vector 扩展失败: %w", err)
	}

	log.Info().Msg("成功启用 pgvector 扩展")

	DB = db
	log.Info().Msg("数据库初始化完成")
	return nil
}

// GormZerologAdapter 实现 GORM 的日志接口
type GormZerologAdapter struct{}

func (g *GormZerologAdapter) Printf(format string, args ...interface{}) {
	log.Debug().Msgf(format, args...)
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}