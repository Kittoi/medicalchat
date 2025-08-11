package models

import "time"

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Message     string  `json:"message" binding:"required"`
	Temperature float32 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Usage     Usage     `json:"usage,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Finished  bool      `json:"finished"`
}

// StreamResponse 流式响应结构
type StreamResponse struct {
	ID      string `json:"id"`
	Delta   string `json:"delta"`
	Done    bool   `json:"done"`
	Created int64  `json:"created"`
}

// Usage token使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAI API相关结构
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float32         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Stream      bool            `json:"stream"`
}

type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Usage   Usage  `json:"usage"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type OpenAIStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}
