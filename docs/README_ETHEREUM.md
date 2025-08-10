# Medical Chat - 以太坊工具包

基于 Go 语言开发的以太坊区块链工具包，提供区块链网络状态查询和存证功能。

## 功能特性

### 区块链网络状态查询
- 获取当前区块链网络信息
- 查看网络延迟和连接状态
- 显示当前区块高度

### 区块链存证功能
- 创建数据存证记录
- 验证存证的有效性
- 查询存证列表和详情
- 支持按状态过滤存证

## 项目结构

```
medicalchat/
├── configs/
│   └── config.yaml              # 配置文件
├── cmd/
│   ├── api/main.go             # 原主程序
│   └── ethereum/main.go        # 以太坊工具包主程序
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── ethereum.go         # 以太坊处理器
│   │   │   └── ethereum_test.go    # 处理器测试
│   │   └── routes/
│   │       └── ethereum.go         # 以太坊路由
│   ├── models/
│   │   ├── config.go            # 配置模型
│   │   └── ethereum.go          # 以太坊模型
│   └── service/
│       ├── ethereum.go          # 以太坊服务
│       └── ethereum_test.go     # 服务测试
├── test/
│   └── ethereum_integration_test.go # 集成测试
└── scripts/
    └── test_ethereum_api.sh     # API测试脚本
```

## 配置说明

在 `configs/config.yaml` 中配置以太坊相关参数：

```yaml
ethereum:
  rpc_url: "https://mainnet.infura.io/v3/YOUR_PROJECT_ID"  # 以太坊RPC节点地址
  private_key: ""  # 私钥(用于交易签名，测试环境使用)
  contract_address: ""  # 存证合约地址
  gas_limit: 21000  # Gas限制
  gas_price: 20000000000  # Gas价格 (20 Gwei)
  chain_id: 1  # 链ID (1为主网，5为Goerli测试网，11155111为Sepolia测试网)
  timeout: 30s  # 超时时间
```

### 配置参数说明

- **rpc_url**: 以太坊节点RPC地址
  - 主网: `https://mainnet.infura.io/v3/YOUR_PROJECT_ID`
  - Goerli测试网: `https://goerli.infura.io/v3/YOUR_PROJECT_ID`
  - Sepolia测试网: `https://sepolia.infura.io/v3/YOUR_PROJECT_ID`
  
- **private_key**: 用于签名交易的私钥（可选，仅创建存证时需要）

- **chain_id**: 网络链ID
  - 1: 以太坊主网
  - 5: Goerli测试网
  - 11155111: Sepolia测试网

## API 端点

### 获取区块链网络状态
```http
GET /api/blockchain/status
```

响应示例：
```json
{
  "network": "Ethereum Mainnet",
  "latency": 234,
  "block_height": 18450024,
  "is_connected": true
}
```

### 创建存证
```http
POST /api/certificate
Content-Type: application/json

{
  "id": "cert_001",
  "content": "存证内容"
}
```

### 验证存证
```http
POST /api/certificate/verify
Content-Type: application/json

{
  "cert_id": "cert_001"
}
```

### 获取存证列表
```http
GET /api/certificate
GET /api/certificate?status=confirmed
```

### 获取存证详情
```http
GET /api/certificate/{id}
```

## 安装和运行

### 1. 安装依赖
```bash
go mod tidy
```

### 2. 配置文件
编辑 `configs/config.yaml`，配置你的以太坊 RPC URL 和其他参数。

### 3. 运行服务
```bash
# 运行以太坊工具包服务
go run cmd/ethereum/main.go
```

服务将在 `http://localhost:8080` 启动。