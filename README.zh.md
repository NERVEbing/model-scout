# model-scout

model-scout 用于扫描 LLM 平台，并探测你的 API Key 可以使用哪些模型。

## 功能特性

- 获取平台模型列表，并用轻量聊天请求探测可用性。
- 支持并发探测，可配置 worker 数量与超时。
- 通过子串过滤不需要的模型。
- 输出 JSON 或 YAML，便于自动化处理。

## 环境要求

- Go 1.25.5（与 `go.mod` 保持一致）
- 平台 API Key（当前支持 DashScope，后续会接入更多平台）

## 构建

```
go build ./cmd/model-scout
```

二进制文件会生成在当前目录的 `./model-scout`。

## 使用

```
model-scout scan --platform dashscope --api-key $DASHSCOPE_API_KEY
```

也可以依赖默认环境变量 `DASHSCOPE_API_KEY`：

```
model-scout scan --platform dashscope
```

### 快速开始

运行扫描并输出 JSON：

```
model-scout scan --platform dashscope --out json
```

JSON 输出示例：

```json
[
  {
    "platform": "dashscope",
    "model": "qwen-plus",
    "status": "ok",
    "available": true,
    "capabilities": ["chat"]
  }
]
```

运行扫描并输出 YAML：

```
model-scout scan --platform dashscope --out yaml
```

YAML 输出示例：

```yaml
- platform: dashscope
  model: qwen-plus
  status: ok
  available: true
  capabilities:
    - chat
```

### 参数说明

- `--platform`（必填）：扫描的平台。支持：`dashscope`。
- `--api-key`：平台 API Key。为空时会读取 `--key-env` 指定的环境变量。
- `--key-env`：API Key 的环境变量名（默认：`DASHSCOPE_API_KEY`）。
- `--workers`：并发探测数（默认：4）。
- `--timeout`：HTTP 超时时间，如 `10s`（默认：`15s`）。
- `--out`：输出格式：`json` 或 `yaml`（默认：`json`）。
- `--output-file`：输出到文件（默认 stdout）。
- `--only-ok`：仅输出可用模型。
- `--exclude`：逗号分隔的排除子串。
- `--filter`：按 `key=value` 或 `key!=value` 过滤输出（可重复，值可用逗号分隔）。

### 过滤规则

支持精确匹配的字段：`available`、`status`、`model`、`platform`。

示例：

```
model-scout scan --platform dashscope --filter available=true
model-scout scan --platform dashscope --filter status=ok,active
model-scout scan --platform dashscope --filter platform=dashscope --filter status=ok
```

### 默认过滤

扫描时会跳过包含以下子串的模型 ID：

```
image, tts, asr, mt, ocr, rerank, embedding, realtime, livetranslate
```

可以使用 `--exclude` 增加其他子串。

## 输出

每条结果包含：

- `platform`：平台名称
- `model`：模型 ID
- `status`：`ok`、`denied`、`unsupported`、`fail` 或 `error`
- `available`：是否可用
- `reason`：失败原因或错误信息（若有）
- `capabilities`：目前成功探测会返回 `chat`

JSON 输出示例：

```json
[
  {
    "platform": "dashscope",
    "model": "qwen-plus",
    "status": "ok",
    "available": true,
    "capabilities": ["chat"]
  }
]
```

## 开发

```
go test ./...
```

## 平台支持

- DashScope（`dashscope`）
- 其他平台将陆续接入

## 安全提示

不要提交 API Key。运行时使用环境变量或 `--api-key`。
