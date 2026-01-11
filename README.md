# model-scout

model-scout scans LLM platforms and probes which models are available for your API key.

[English](README.md) | [中文](README.zh.md)

## Features

- List models from a platform and probe availability with a lightweight chat request.
- Concurrent probing with configurable workers and timeout.
- Filter out unwanted models by substring.
- JSON or YAML output for automation.

## Requirements

- Go 1.25.5 (matches `go.mod`)
- A platform API key (DashScope and DeepSeek today; more platforms coming)

## Build

```
go build ./cmd/model-scout
```

The binary is written to `./model-scout` in the current directory.

## Usage

```
model-scout scan --platform dashscope --api-key $DASHSCOPE_API_KEY
```

You can also rely on the default environment variable `DASHSCOPE_API_KEY`:

```
model-scout scan --platform dashscope
```

For DeepSeek, the default environment variable is `DEEPSEEK_API_KEY`:

```
model-scout scan --platform deepseek
```

### Quickstart

Run a scan and output JSON:

```
model-scout scan --platform dashscope --out json
```

Example JSON output:

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

Run a scan and output YAML:

```
model-scout scan --platform dashscope --out yaml
```

Example YAML output:

```yaml
- platform: dashscope
  model: qwen-plus
  status: ok
  available: true
  capabilities:
    - chat
```

### Flags

- `--platform` (required): platform to scan. Supported: `dashscope`.
- `--api-key`: platform API key. If empty, the platform default environment variable is used.
- `--workers`: number of concurrent probes (default: 4).
- `--timeout`: HTTP timeout, e.g. `10s` (default: `15s`).
- `--out`: output format: `json` or `yaml` (default: `json`).
- `--output-file`: write output to a file (defaults to stdout).
- `--exclude`: comma-separated substrings to exclude.
- `--filter`: filter output with `key=value` or `key!=value` (repeatable, values can be comma-separated).

### Filters

Filters support exact matching on these keys: `available`, `status`, `model`, `platform`.

Examples:

```
model-scout scan --platform dashscope --filter available=true
model-scout scan --platform dashscope --filter status=ok,active
model-scout scan --platform dashscope --filter platform=dashscope --filter status=ok
```

### Default filters

The scout skips model IDs containing:

```
image, tts, asr, mt, ocr, rerank, embedding, realtime, livetranslate
```

Use `--exclude` to add more substrings.

## Output

Each result includes:

- `platform`: platform name
- `model`: model ID
- `status`: `ok`, `denied`, `unsupported`, `fail`, or `error`
- `available`: boolean
- `reason`: error or failure message (if any)
- `capabilities`: currently `chat` for successful probes

Example JSON output:

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

## Development

```
go test ./...
```

## Platforms

- DashScope (`dashscope`)
- DeepSeek (`deepseek`)
- More platforms will be added

## Security

Do not commit API keys. Use environment variables or `--api-key` at runtime.
