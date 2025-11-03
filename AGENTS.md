# Agent Guidelines

## Build & Test
- **Build Lambda**: `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -C cmd/lambda -o ../../dist/bootstrap`
- **Test all**: `make test` (runs with `-race -count=1`)
- **Test single package**: `go test -race -count=1 ./internal/app`
- **Test single function**: `go test -race -count=1 ./internal/app -run TestFunctionName`
- **Run sample locally**: `go run -C cmd/sample .` (requires `.env` file)
- **Lint**: `go vet ./...` and `gofmt -l .`

## Code Style
- **Imports**: stdlib, then blank line, then third-party, then local (e.g., `internal/`)
- **Naming**: Go standard - `PascalCase` exports, `camelCase` private, `ALL_CAPS` for env vars prefixed with `APP_`
- **Error handling**: return errors up the stack; use `fmt.Errorf` for wrapping
- **Structs**: define types in package, constructors as `New()` or `NewTypeName()`; all methods must be public (PascalCase)
- **Interfaces**: keep minimal (e.g., `SecurityHubEvent` has 2 methods)
- **Formatting**: use `gofmt` (tabs for indentation)
- **Comments**: rare, lowercase, short, concise; code should be self-documenting
- **Code smells**: keep to minimum; prefer clear naming over comments

## Architecture
- `cmd/lambda/main.go` - Lambda handler entry point
- `cmd/sample/main.go` - Local development runner using fixtures
- `internal/app/` - Core application logic and configuration
- `internal/events/` - OCSF event parsing and Slack message formatting
- `fixtures/samples.json` - Sample Security Hub v2 OCSF findings for testing

## Important Notes
- This project is specifically for **AWS Security Hub v2** which uses OCSF (Open Cybersecurity Schema Framework) format
- It is NOT compatible with the original AWS Security Hub (now called Security Hub CSPM) ASFF format
- Security Hub v2 centralizes findings from GuardDuty, Inspector, Macie, IAM Access Analyzer, and Security Hub CSPM
- Events use OCSF fields like `finding_info`, `metadata`, `severity`, `class_name`, etc.
