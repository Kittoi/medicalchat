package handler

import (
	"context"
	"io"
	"medicalchat/internal/models"
	"medicalchat/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}


// ChatStream 流式聊天接口
// @Summary 流式聊天
// @Description 向AI发送消息并获取流式回复
// @Tags chat
// @Accept json
// @Produce text/event-stream
// @Param request body models.ChatRequest true "聊天请求"
// @Success 200 {object} models.StreamResponse "流式响应"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /api/chat/stream [post]
func (ch *ChatHandler) ChatStream(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("绑定请求参数失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数格式错误: " + err.Error(),
		})
		return
	}

	// 设置流式响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Cache-Control")

	// 创建响应通道
	responseChan := make(chan models.StreamResponse, 100)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 启动流式聊天协程
	go func() {
		if err := ch.chatService.ChatStream(ctx, req, responseChan); err != nil {
			log.Error().Err(err).Msg("流式聊天服务出错")
			// 发送错误消息
			select {
			case responseChan <- models.StreamResponse{
				Delta: "抱歉，服务出现错误：" + err.Error(),
				Done:  true,
			}:
			default:
			}
		}
	}()

	// 监听客户端断开连接
	clientGone := c.Writer.CloseNotify()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			log.Info().Msg("客户端断开连接")
			cancel() // 取消上下文，停止AI服务调用
			return false
		case response, ok := <-responseChan:
			if !ok {
				log.Debug().Msg("响应通道已关闭")
				return false
			}

			// 发送SSE格式的数据
			c.SSEvent("message", response.Delta)
			return !response.Done
		case <-ctx.Done():
			log.Info().Msg("请求超时")
			c.SSEvent("error", map[string]string{
				"error": "请求超时",
			})
			return false
		}
	})
}

