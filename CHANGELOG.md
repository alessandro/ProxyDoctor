# Changelog

All notable changes to this project will be documented in this file.

## Unreleased
### Fixed
- `cli diagnose --proxy`: the flag value was parsed and then discarded — `parseProxyConfig` always returned a hardcoded `localhost:8080` HTTP config regardless of input. Now the proxy URL passed via `--proxy` is actually parsed (scheme, host, port, optional credentials), and `--proxy-type` (`auto`, `http`, `https`, `socks4`, `socks5`) is respected. `auto` infers the type from the URL scheme.

## v0.1.0 - YYYY-MM-DD
- Initial project snapshot: CLI and core components.