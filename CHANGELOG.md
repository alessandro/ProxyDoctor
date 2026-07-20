# Changelog

All notable changes to this project will be documented in this file.

## Unreleased
### Added
- Web GUI at `http://localhost:8080/`: a single-page form (URL, proxy, proxy type) that runs a real diagnosis and renders results in the browser.
- `POST /api/diagnose`: runs a diagnosis through the real core engine (`core/engine.DiagnosisOrchestrator`), same code path as `cli diagnose`. Accepts `{"url", "proxy", "proxy_type"}` JSON, returns a full `DiagnosisReport`.
- `GET /api/checks`: lists registered checks as JSON (web equivalent of `cli list-checks`).
- `core/utils.ParseProxyConfig`: proxy URL parsing extracted into a shared helper used by both the CLI and the server, so both behave identically.
- JSON tags on `engine.DiagnosisReport` / `engine.RequestMetadata` (snake_case, consistent with `check.CheckResult`).

### Changed
- `cmd/server`: rewritten from a standalone placeholder (direct calls to ipify/icanhazip, no core import) into a thin HTTP layer over the core engine. The old `GET /api/check/public-ip` endpoint was removed in favor of `POST /api/diagnose`, which covers the same case and more.

### Fixed
- `cli diagnose --proxy`: the flag value was parsed and then discarded — `parseProxyConfig` always returned a hardcoded `localhost:8080` HTTP config regardless of input. Now the proxy URL passed via `--proxy` is actually parsed (scheme, host, port, optional credentials), and `--proxy-type` (`auto`, `http`, `https`, `socks4`, `socks5`) is respected. `auto` infers the type from the URL scheme.

## v0.1.0 - YYYY-MM-DD
- Initial project snapshot: CLI and core components.