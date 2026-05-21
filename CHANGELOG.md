# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added

- Initial repository structure
- `.gitignore` with secrets, Python, Terraform, Helm, and editor layers
- `README.md` with architecture overview, tech stack, and pipeline documentation
- `CHANGELOG.md` following Keep a Changelog standard
- `pull_request_template.md` for standardized PR process
- `.pre-commit-config.yaml` for local enforcement of code quality
- Go application (FastAPI replaced with Go + Chi router)
- CRUD API with `/api/v1/items/` endpoints
- Kubernetes probes: `/healthz` (liveness) and `/readyz` (readiness)
- Security middleware with OWASP headers
- Request body size limit (1MB) and Content-Type validation
- Structured JSON logging with `slog`
- Graceful shutdown with SIGTERM/SIGINT handling
- 16 unit tests with 98.3% handler coverage
- Multi-stage Dockerfile with `distroless/static:nonroot` runtime (~3.7MB image)
- `.golangci.yml` with gosec, errcheck, bodyclose, and gocritic linters
- `.env.example` documenting environment variables

---

## [0.1.0] - 2026-05-21

### Added

- Project scaffolding and repository initialization

---

[Unreleased]: https://github.com/edsuarezs/devsecops-pipeline-demo/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/edsuarezs/devsecops-pipeline-demo/releases/tag/v0.1.0
