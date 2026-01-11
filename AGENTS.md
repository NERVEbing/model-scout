# Repository Guidelines

## Project Structure & Module Organization
- `cmd/model-scout/` holds the CLI entrypoint (`main.go`).
- `internal/` contains app logic and integrations (scout engine, platform clients, output encoders).
- `pkg/` exposes shared types meant for external use.
- Tests live alongside code as `*_test.go` files (currently under `internal/`).

## Build, Test, and Development Commands
- `go build ./cmd/model-scout` builds the CLI binary.
- `go run ./cmd/model-scout scan --platform dashscope --api-key $DASHSCOPE_API_KEY` runs a scan locally.
- `go test ./...` runs all unit tests.
- `go test ./internal/scout -run TestEngineScanFilters` runs a focused test.

## Coding Style & Naming Conventions
- Use standard Go formatting (`gofmt`), tabs for indentation, and Go idioms.
- Package names are short, lowercase, and singular (e.g., `scout`, `platform`).
- File names are lowercase with underscores where needed (e.g., `engine_test.go`).

## Testing Guidelines
- Use Go’s `testing` package with `TestXxx` naming.
- Prefer table-driven tests when multiple cases share setup.
- Keep tests close to the code they validate; ensure new features include coverage.

## Commit & Pull Request Guidelines
- Commits in this repo use short, imperative subjects (e.g., “Update module for Go 1.25.5”).
- Keep commits scoped and focused; avoid mixing unrelated changes.
- PRs should include: purpose summary, testing performed (`go test ./...`), and any relevant CLI output or sample JSON/YAML.

## Security & Configuration Tips
- API keys are read from `DASHSCOPE_API_KEY` by default or `--api-key`.
- Avoid committing secrets; use env vars or local shell profiles.
