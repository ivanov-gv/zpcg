# Montenegro Railways Timetable Bot

A Telegram bot that provides instant access to train timetables for the Montenegro railway system. Stateless,
cost-effective, running on minimal cloud resources (1vCPU, 128MB RAM) with zero external dependencies.

## Design principles

- **Stateless & in-memory only** — No persistence, no external database or cache. Timetable compiled into the binary.
- **Cost-effective cloud-native** — Runs on Google Cloud Run with automatic scaling. ~€0.03/month for 150 users.
- **Custom path-finding algorithm** — Leverages Montenegro's railway topology (tree structure with Podgorica as hub).
  Faster than Dijkstra for this specific problem.
- **User-friendly interface** — Supports Cyrillic & Latin input, handles typos, clear error messages, multi-language
  support.

## Language and tools

- **Language**: Go 1.26.1
- **Build**: `make build` (Docker); `make new_version` (full pipeline: parse timetable + build + push + deploy)
- **Test**: `make test_unit` (unit tests); `make test_integration` (with TDLib); `make test` (all tests)
- **Lint**: golangci-lint (via `.golangci.yml`)
- **Logging**: zerolog
- **Telegram API**: gotgbot/v2 + go-tdlib
- **Testing**: testify + mockery for mocks
- **Deployment**: Docker + Google Cloud Run

## Conventions

- Standard Go project layout: `cmd/` (exporter, tg-init, tg-server), `internal/` (core logic), `test/` (integration
  tests), `gen/` (generated code)
- Configuration via environment variables and `.env` files
- `cmd/exporter` generates `gen/timetable/timetable.gen.go` by parsing the official ZPCG website
- `cmd/tg-server` and `cmd/tg-init` handle Telegram bot operations
- Tests require race detection: `go test -race`; integration tests require TDLib
- CI testing: use `act` locally to test workflows before pushing (see Makefile targets)