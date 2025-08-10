package models

type Config struct {
	PGVector PGVectorConfig `mapstructure:"pgvector"`
	Log      LogConfig      `mapstructure:"log"`
}

type PGVectorConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type LogConfig struct {
	Dir   string `mapstructure:"dir"`    // 日志目录
	Level string `mapstructure:"level"`  // 日志级别
}