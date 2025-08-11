package models

import "time"

type Config struct {
	PGVector PGVectorConfig `mapstructure:"pgvector"`
	Log      LogConfig      `mapstructure:"log"`
	Ethereum EthereumConfig `mapstructure:"ethereum"`
	Chat     ChatConfig     `mapstructure:"chat"`
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
	Dir   string `mapstructure:"dir"`   // 日志目录
	Level string `mapstructure:"level"` // 日志级别
}

type EthereumConfig struct {
	RPCURL          string        `mapstructure:"rpc_url"`
	PrivateKey      string        `mapstructure:"private_key"`
	ContractAddress string        `mapstructure:"contract_address"`
	GasLimit        uint64        `mapstructure:"gas_limit"`
	GasPrice        int64         `mapstructure:"gas_price"`
	ChainID         int64         `mapstructure:"chain_id"`
	Timeout         time.Duration `mapstructure:"timeout"`
}

type ChatConfig struct {
	APIKey      string        `mapstructure:"api_key"`
	BaseURL     string        `mapstructure:"base_url"`
	Model       string        `mapstructure:"model"`
	MaxTokens   int           `mapstructure:"max_tokens"`
	Temperature float32       `mapstructure:"temperature"`
	Timeout     time.Duration `mapstructure:"timeout"`
	Stream      bool          `mapstructure:"stream"`
}
