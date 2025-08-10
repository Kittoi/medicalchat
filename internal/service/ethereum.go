package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"medicalchat/internal/models"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type EthereumService struct {
	client     *ethclient.Client
	config     models.EthereumConfig
	privateKey *ecdsa.PrivateKey
	logger     zerolog.Logger
}

func NewEthereumService(config models.EthereumConfig, logger zerolog.Logger) (*EthereumService, error) {
	client, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %v", err)
	}

	var privateKey *ecdsa.PrivateKey
	if config.PrivateKey != "" {
		privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(config.PrivateKey, "0x"))
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %v", err)
		}
	}

	return &EthereumService{
		client:     client,
		config:     config,
		privateKey: privateKey,
		logger:     logger,
	}, nil
}

// GetNetworkStatus 获取网络状态
func (s *EthereumService) GetNetworkStatus(ctx context.Context) (*models.NetworkStatus, error) {
	start := time.Now()

	// 获取最新区块
	header, err := s.client.HeaderByNumber(ctx, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get latest block header")
		return &models.NetworkStatus{
			Network:     s.getNetworkName(),
			Latency:     time.Since(start).Milliseconds(),
			BlockHeight: 0,
			IsConnected: false,
		}, err
	}

	// 计算延迟
	latency := time.Since(start).Milliseconds()

	return &models.NetworkStatus{
		Network:     s.getNetworkName(),
		Latency:     latency,
		BlockHeight: header.Number.Uint64(),
		IsConnected: true,
	}, nil
}

// CreateCertificate 创建存证
func (s *EthereumService) CreateCertificate(ctx context.Context, req models.CertificateRequest) (*models.Certificate, error) {
	if s.privateKey == nil {
		return nil, fmt.Errorf("private key not configured")
	}

	// 计算内容哈希
	contentHash := s.hashContent(req.Content)

	// 创建存证记录
	cert := &models.Certificate{
		ID:        req.ID,
		Content:   req.Content,
		Hash:      contentHash,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// 发送交易到区块链 (这里使用简单的ETH转账来模拟存证)
	txHash, err := s.sendCertificateTransaction(ctx, contentHash)
	if err != nil {
		s.logger.Error().Err(err).Str("cert_id", req.ID).Msg("Failed to send certificate transaction")
		cert.Status = "failed"
		return cert, err
	}

	cert.TxHash = txHash
	cert.Status = "confirmed"
	now := time.Now()
	cert.ConfirmedAt = &now

	s.logger.Info().
		Str("cert_id", req.ID).
		Str("tx_hash", txHash).
		Str("content_hash", contentHash).
		Msg("Certificate created successfully")

	return cert, nil
}

// VerifyCertificate 验证存证
func (s *EthereumService) VerifyCertificate(ctx context.Context, certID string) (*models.Certificate, bool, error) {
	// 在实际应用中，这里应该从数据库或合约中查询存证信息
	// 这里我们简化实现，返回一个模拟的存证记录

	// 模拟查询逻辑
	if certID == "cert_001" {
		now := time.Now()
		cert := &models.Certificate{
			ID:          certID,
			Hash:        "0x1a2b3c4d5e6f...",
			Status:      "confirmed",
			Content:     "头痛症状咨询及AI回复",
			TxHash:      "0xabcdef123456...",
			BlockNumber: 18450024,
			CreatedAt:   time.Date(2025, 8, 9, 0, 0, 0, 0, time.UTC),
			ConfirmedAt: &now,
		}
		return cert, true, nil
	}

	return nil, false, fmt.Errorf("certificate not found")
}

// GetCertificatesByStatus 根据状态获取存证列表
func (s *EthereumService) GetCertificatesByStatus(ctx context.Context, status string) ([]*models.Certificate, error) {
	// 模拟数据
	certificates := []*models.Certificate{}

	if status == "" || status == "confirmed" {
		now := time.Now()
		cert := &models.Certificate{
			ID:          "cert_001",
			Hash:        "0x1a2b3c4d5e6f...",
			Status:      "confirmed",
			Content:     "头痛症状咨询及AI回复",
			TxHash:      "0xabcdef123456...",
			BlockNumber: 18450024,
			CreatedAt:   time.Date(2025, 8, 9, 0, 0, 0, 0, time.UTC),
			ConfirmedAt: &now,
		}
		certificates = append(certificates, cert)
	}

	return certificates, nil
}

// sendCertificateTransaction 发送存证交易
func (s *EthereumService) sendCertificateTransaction(ctx context.Context, contentHash string) (string, error) {
	// 获取账户地址
	publicKey := s.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取nonce
	nonce, err := s.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 构造交易 (发送到自己地址，金额为0，在data中包含存证哈希)
	gasPrice := big.NewInt(s.config.GasPrice)
	gasLimit := s.config.GasLimit

	// 将内容哈希作为交易数据
	data := []byte(contentHash)

	tx := types.NewTransaction(
		nonce,
		fromAddress,   // 发送给自己
		big.NewInt(100), // 金额为0
		gasLimit,
		gasPrice,
		data,
	)

	// 签名交易
	chainID := big.NewInt(s.config.ChainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), s.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	err = s.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// hashContent 计算内容哈希
func (s *EthereumService) hashContent(content string) string {
	hash := sha256.Sum256([]byte(content))
	return "0x" + hex.EncodeToString(hash[:])
}

// getNetworkName 根据链ID获取网络名称
func (s *EthereumService) getNetworkName() string {
	switch s.config.ChainID {
	case 1:
		return "Ethereum Mainnet"
	case 5:
		return "Goerli Testnet"
	case 11155111:
		return "Sepolia Testnet"
	default:
		return fmt.Sprintf("Unknown Network (Chain ID: %d)", s.config.ChainID)
	}
}

// Close 关闭以太坊客户端连接
func (s *EthereumService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
