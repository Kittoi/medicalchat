package handler

import (
	"medicalchat/internal/models"
	"medicalchat/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EthereumHandler struct {
	ethService *service.EthereumService
	logger     zerolog.Logger
}

func NewEthereumHandler(ethService *service.EthereumService, logger zerolog.Logger) *EthereumHandler {
	return &EthereumHandler{
		ethService: ethService,
		logger:     logger,
	}
}

// GetNetworkStatus 获取区块链网络状态
func (h *EthereumHandler) GetNetworkStatus(c *gin.Context) {
	if h.ethService == nil {
		h.logger.Error().Msg("Ethereum service not available")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Ethereum service not available",
		})
		return
	}

	ctx := c.Request.Context()

	status, err := h.ethService.GetNetworkStatus(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get network status")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get network status",
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// CreateCertificate 创建存证
func (h *EthereumHandler) CreateCertificate(c *gin.Context) {
	if h.ethService == nil {
		h.logger.Error().Msg("Ethereum service not available")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Ethereum service not available",
		})
		return
	}

	var req models.CertificateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid certificate request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	ctx := c.Request.Context()
	cert, err := h.ethService.CreateCertificate(ctx, req)
	if err != nil {
		h.logger.Error().Err(err).Str("cert_id", req.ID).Msg("Failed to create certificate")

		// 如果cert为nil，创建一个带错误状态的证书
		if cert == nil {
			cert = &models.Certificate{
				ID:        req.ID,
				Content:   req.Content,
				Status:    "failed",
				CreatedAt: time.Now(),
			}
		}

		c.JSON(http.StatusInternalServerError, models.CertificateResponse{
			Certificate: *cert,
			Message:     "Failed to create certificate: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.CertificateResponse{
		Certificate: *cert,
		Message:     "Certificate created successfully",
	})
}

// VerifyCertificate 验证存证
func (h *EthereumHandler) VerifyCertificate(c *gin.Context) {
	if h.ethService == nil {
		h.logger.Error().Msg("Ethereum service not available")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Ethereum service not available",
		})
		return
	}

	var req models.VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid verify request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	ctx := c.Request.Context()
	cert, valid, err := h.ethService.VerifyCertificate(ctx, req.CertID)
	if err != nil {
		h.logger.Error().Err(err).Str("cert_id", req.CertID).Msg("Failed to verify certificate")
		c.JSON(http.StatusNotFound, models.VerifyResponse{
			Valid:   false,
			Message: "Certificate not found",
		})
		return
	}

	message := "验证成功，存证有效"
	if !valid {
		message = "验证失败，存证无效"
	}

	c.JSON(http.StatusOK, models.VerifyResponse{
		Valid:       valid,
		Certificate: *cert,
		Message:     message,
	})
}

// GetCertificates 获取存证列表
func (h *EthereumHandler) GetCertificates(c *gin.Context) {
	log.Info().Msgf("Getting certificates")
	if h.ethService == nil {
		h.logger.Error().Msg("Ethereum service not available")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Ethereum service not available",
		})
		return
	}

	status := c.Query("status") // 可选的状态过滤参数

	ctx := c.Request.Context()
	certificates, err := h.ethService.GetCertificatesByStatus(ctx, status)
	if err != nil {
		h.logger.Error().Err(err).Str("status", status).Msg("Failed to get certificates")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get certificates",
		})
		return
	}

	// 统计各状态的存证数量
	stats := map[string]int{
		"total":     len(certificates),
		"pending":   0,
		"confirmed": 0,
		"failed":    0,
	}

	for _, cert := range certificates {
		stats[cert.Status]++
	}

	c.JSON(http.StatusOK, gin.H{
		"certificates": certificates,
		"statistics":   stats,
	})
}

// GetCertificateByID 根据ID获取存证详情
func (h *EthereumHandler) GetCertificateByID(c *gin.Context) {
	if h.ethService == nil {
		h.logger.Error().Msg("Ethereum service not available")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Ethereum service not available",
		})
		return
	}

	certID := c.Param("id")
	if certID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Certificate ID is required",
		})
		return
	}

	ctx := c.Request.Context()
	cert, valid, err := h.ethService.VerifyCertificate(ctx, certID)
	if err != nil {
		h.logger.Error().Err(err).Str("cert_id", certID).Msg("Certificate not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Certificate not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"certificate": cert,
		"valid":       valid,
	})
}
