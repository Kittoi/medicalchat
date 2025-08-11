package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"medicalchat/internal/models"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type ChatService struct {
	config models.ChatConfig
	client *http.Client
}

func NewChatService(config models.ChatConfig) *ChatService {
	fmt.Printf("创建聊天服务: %+v\n", config)
	return &ChatService{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// ChatStream 流式聊天
func (cs *ChatService) ChatStream(ctx context.Context, req models.ChatRequest, responseChan chan<- models.StreamResponse) error {
	defer close(responseChan)

	openAIReq := models.OpenAIChatRequest{
		Model: cs.config.Model,
		Messages: []models.OpenAIMessage{
			{
				Role:    "user",
				Content: req.Message,
			},
		},
		Temperature: cs.getTemperature(req.Temperature),
		MaxTokens:   cs.getMaxTokens(req.MaxTokens),
		Stream:      true,
	}

	resp, err := cs.makeStreamRequest(ctx, openAIReq)
	if err != nil {
		return fmt.Errorf("请求AI流式服务失败: %w", err)
	}
	defer resp.Body.Close()
	allDeltaByAI := ""
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			responseChan <- models.StreamResponse{
				Done: true,
			}
			// 发送结束标志
			fmt.Println("接收到结束标志，退出循环")
			break
		}

		var streamResp models.OpenAIStreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			log.Error().Err(err).Str("data", data).Msg("解析流式响应失败")
			continue
		}

		if len(streamResp.Choices) > 0 {
			delta := streamResp.Choices[0].Delta.Content
			if delta != "" {
				allDeltaByAI += delta
				responseChan <- models.StreamResponse{
					ID:      streamResp.ID,
					Delta:   delta,
					Done:    false,
					Created: streamResp.Created,
				}
			}
		}
	}
	fmt.Println("所有Delta:", allDeltaByAI)
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取流式响应失败: %w", err)
	}

	return nil
}


func (cs *ChatService) makeStreamRequest(ctx context.Context, req models.OpenAIChatRequest) (*http.Response, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", cs.config.BaseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cs.config.APIKey)
	httpReq.Header.Set("Accept", "text/event-stream")
	httpReq.Header.Set("Cache-Control", "no-cache")

	resp, err := cs.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("AI服务返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (cs *ChatService) getTemperature(temp float32) float32 {
	if temp > 0 {
		return temp
	}
	return cs.config.Temperature
}

func (cs *ChatService) getMaxTokens(maxTokens int) int {
	if maxTokens > 0 {
		return maxTokens
	}
	return cs.config.MaxTokens
}
