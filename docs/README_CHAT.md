# AI聊天模块使用说明

本模块实现了与AI服务的聊天功能，支持流式和非流式两种交互方式。

## 配置说明

在 `configs/config.yaml` 中配置聊天服务参数：

```yaml
chat:
  api_key: "your-openai-api-key"  # AI服务API Key
  base_url: "https://api.openai.com/v1"  # AI服务请求地址
  model: "gpt-3.5-turbo"  # 使用的模型
  max_tokens: 2048  # 最大token数
  temperature: 0.7  # 温度参数
  timeout: 30s  # 请求超时时间
  stream: true  # 是否启用流式输出
```

## API接口

### 1. 非流式聊天

**请求接口：** `POST /api/v1/chat`

**请求体：**
```json
{
  "message": "你好，请介绍一下自己",
  "temperature": 0.7,
  "max_tokens": 1000,
  "stream": false
}
```

**响应体：**
```json
{
  "id": "chatcmpl-123",
  "message": "你好！我是一个AI助手...",
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 50,
    "total_tokens": 60
  },
  "created_at": "2025-01-10T12:00:00Z",
  "finished": true
}
```

### 2. 流式聊天

**请求接口：** `POST /api/v1/chat/stream`

**请求体：**
```json
{
  "message": "请详细介绍人工智能",
  "temperature": 0.7,
  "max_tokens": 2000,
  "stream": true
}
```

**响应格式：** Server-Sent Events (SSE)

**响应示例：**
```
data: {"id":"chatcmpl-123","delta":"人工","done":false,"created":1641811200}

data: {"id":"chatcmpl-123","delta":"智能","done":false,"created":1641811200}

data: {"id":"chatcmpl-123","delta":"是...","done":false,"created":1641811200}

data: {"id":"chatcmpl-123","delta":"","done":true,"created":1641811200}
```

### 3. 获取模型列表

**请求接口：** `GET /api/v1/chat/models`

**响应体：**
```json
{
  "models": [
    {
      "id": "gpt-3.5-turbo",
      "name": "GPT-3.5 Turbo",
      "description": "快速响应的通用AI模型",
      "max_tokens": 4096
    },
    {
      "id": "gpt-4",
      "name": "GPT-4", 
      "description": "更强大的AI模型",
      "max_tokens": 8192
    }
  ],
  "total": 2
}
```

## 使用示例

### JavaScript/前端示例

**非流式聊天：**
```javascript
async function chat(message) {
  const response = await fetch('/api/v1/chat', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      message: message,
      temperature: 0.7,
      max_tokens: 1000
    })
  });
  
  const result = await response.json();
  console.log(result.message);
}
```

**流式聊天：**
```javascript
async function chatStream(message, onMessage) {
  const response = await fetch('/api/v1/chat/stream', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      message: message,
      temperature: 0.7,
      max_tokens: 2000
    })
  });

  const reader = response.body.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { value, done } = await reader.read();
    if (done) break;

    const chunk = decoder.decode(value);
    const lines = chunk.split('\n');
    
    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = JSON.parse(line.slice(6));
        if (data.done) {
          return;
        }
        onMessage(data.delta);
      }
    }
  }
}

// 使用示例
chatStream("请介绍人工智能", (delta) => {
  console.log(delta); // 实时输出AI回复的每个部分
});
```

### cURL示例

**非流式聊天：**
```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "你好",
    "temperature": 0.7,
    "max_tokens": 1000
  }'
```

**流式聊天：**
```bash
curl -X POST http://localhost:8080/api/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "message": "请介绍人工智能",
    "temperature": 0.7,
    "max_tokens": 2000
  }' \
  --no-buffer
```

## 错误处理

API可能返回的错误状态码：

- `400 Bad Request` - 请求参数错误
- `500 Internal Server Error` - 服务器内部错误

错误响应格式：
```json
{
  "error": "错误描述信息"
}
```

## 部署说明

1. 确保配置文件中的API Key和Base URL正确
2. 启动服务：`go run cmd/api/main.go`
3. 服务将在 `:8080` 端口启动
4. 访问 `http://localhost:8080/api/v1/chat/models` 验证服务是否正常

## 注意事项

1. API Key请妥善保管，不要泄露
2. 流式聊天会保持连接，注意处理连接断开的情况
3. 根据使用的AI服务调整base_url和model参数
4. 合理设置timeout避免长时间等待
