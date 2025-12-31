# ec-tmpl Agent Notes

## Overview
- Minimal API template using Go and Echo, designed with integrated testing, CI, and Docker in mind.
- Ships only the essentials: dependency injection using interfaces and implementation swapping, health and greeting routes, and test/CI/Docker wiring.
- The code is database-independent and follows Go conventions for test placement (adjacent to the target package).

## Design Philosophy
- Stay database-agnostic; add persistence only when the target project needs it.
- Use dependency injection with interfaces for service abstractions and factory pattern for implementations to maximize extensibility, maintainability, and testability.
- Keep settings and dependencies explicit via `Settings` struct and dependency builders.
- Maintain parity between local, Docker, and CI flows with a single source of truth (`just`, Go modules, `.env`).

## First Steps When Creating a Real API
1. Clone or copy the template and run `just setup` to install dependencies and set up the environment.
2. Rename the Go module from `example.com/ec-tmpl` if you need a project-specific namespace.
3. Extend `internal/api/routes.go` with domain routes and register required dependencies in `internal/dependencies/wiring_*.go`.
4. Update `.env.example` and documentation to reflect new environment variables or external services.

## Key Files
- `internal/dependencies/wiring_*.go`: central place to wire settings and service providers (dev and non-dev variants).
- `internal/api/httpserver.go`: Echo server instantiation; attach new routes here via `RegisterRoutes`.
- `internal/protocols/`: interface definitions for service abstractions.
- `internal/services/`: concrete service implementations.
- `tests/`: unit/intg/api/e2e layout kept light so additional checks can drop in without restructuring.

## Tooling Snapshot
- `justfile`: run/lint/test tasks (`dev`, `test`, `check`, `fix`) used locally and in CI. Prefer `just test` as the unified entrypoint.
- `go.mod` + `go.sum`: reproducible dependency graph; update with `go mod tidy` when deps change.
