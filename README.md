# ec-tmpl

`ec-tmpl` is a minimal API template using Go and Echo, designed with integrated testing, CI, and Docker in mind. The code is database-independent, based on dependency injection (interface + implementation swapping), and follows Go conventions for test placement (adjacent to the target package).

## Prerequisites
- Go 1.22+
- Docker / Docker Compose (used for Docker-based tests and compose execution)
- `just` (used for command aggregation)

## Execution Flow (`just`)
- `just dev` starts the API locally with hot reload (uses `EC_TMPL_DEV_PORT` as the port).
- `just test` runs unit, integration, and Docker (api / e2e) tests.
- `just check` runs gofmt checks and static analysis.
- `just fix` runs gofmt and automatic fixes, which may modify the code.

## Endpoints
This template provides the following endpoints:

```http
GET /health -> {"status":"ok"}
GET /hello/:name -> {"message":"Hello, <name>"}
```

## Configuration (Environment Variables)
Environment variables can be loaded from `.env` (see `.env.example` for a template).

- `EC_TMPL_APP_NAME`: Application name
- `EC_TMPL_BIND_IP` / `EC_TMPL_BIND_PORT`: Server bind settings
- `EC_TMPL_DEV_PORT`: Port for local development (used by `just dev`)
- `EC_TMPL_USE_MOCK_GREETING`: Swap greeting implementation (dev build only)
- `EC_TMPL_SHUTDOWN_GRACE_SEC`: Grace period in seconds after receiving SIGTERM

## Testing Policy
- Unit tests are placed as `*_test.go` files adjacent to the target package.
- Integration tests are `*_test.go` files with `//go:build intg` and are enabled with `go test -tags=intg ./...`.
- Docker-based black-box tests are located in `tests/api` and `tests/e2e`, distinguished by `//go:build api` / `//go:build e2e`.
- Docker-based tests start the API container via `testcontainers-go` using the image specified by `EC_TMPL_E2E_IMAGE`.
- `EC_TMPL_TEST_BASE_URL` overrides container startup and targets an existing API base URL.
