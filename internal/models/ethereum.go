package models

import "time"

// NetworkStatus 区块链网络状态
type NetworkStatus struct {
	Network     string `json:"network"`      // 网络名称 (如 Ethereum Mainnet)
	Latency     int64  `json:"latency"`      // 延迟 (毫秒)
	BlockHeight uint64 `json:"block_height"` // 区块高度
	IsConnected bool   `json:"is_connected"` // 是否连接成功
}

// Certificate 存证记录
type Certificate struct {
	ID          string     `json:"id"`           // 存证ID
	Hash        string     `json:"hash"`         // 区块链哈希
	Status      string     `json:"status"`       // 状态 (pending/confirmed/failed)
	Content     string     `json:"content"`      // 存证内容
	TxHash      string     `json:"tx_hash"`      // 交易哈希
	BlockNumber uint64     `json:"block_number"` // 区块号
	CreatedAt   time.Time  `json:"created_at"`   // 创建时间
	ConfirmedAt *time.Time `json:"confirmed_at"` // 确认时间
}

// CertificateRequest 存证请求
type CertificateRequest struct {
	ID      string `json:"id" binding:"required"`      // 存证ID
	Content string `json:"content" binding:"required"` // 存证内容
}

// CertificateResponse 存证响应
type CertificateResponse struct {
	Certificate Certificate `json:"certificate"`
	Message     string      `json:"message"`
}

// VerifyRequest 验证请求
type VerifyRequest struct {
	CertID string `json:"cert_id" binding:"required"` // 存证ID
}

// VerifyResponse 验证响应
type VerifyResponse struct {
	Valid       bool        `json:"valid"`       // 是否有效
	Certificate Certificate `json:"certificate"` // 存证信息
	Message     string      `json:"message"`     // 验证信息
}
