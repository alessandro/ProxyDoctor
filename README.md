# ProxyDoctor — Proxy Diagnostics

**A minimal, testable Go tool for running network diagnostics through proxies**

[![License: GPL-3.0](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Good First Issue](https://img.shields.io/badge/good%20first-issue-brightgreen.svg)](https://github.com/francomano/ProxyDoctor/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8.svg?logo=go)](https://go.dev/)

<p align="center">
  <img src="images/proxydoctor-logo.png" alt="ProxyDoctor Logo" width="150">
</p>

> **Built with AI assistance** — This project was developed and documented with contributions from:
> - **GitHub Copilot** (Auto Mode, inline suggestions)
> - **OpenCode** with free open-source models (big-pickle / opencode models)
> - **Google Gemini Pro** (architecture design, code review)

## Why ProxyDoctor?

**The Problem:**
You're behind a corporate proxy or using a VPN, and things break mysteriously:
- Sites don't load (proxy misconfiguration?)
- Your IP leaks even through a "private" proxy
- You can't tell if the problem is your proxy, DNS, or the site itself

**The Solution:**
ProxyDoctor is a **lightweight diagnostic tool** that:
- Runs network checks through any proxy (HTTP, HTTPS, SOCKS4, SOCKS5)
- Compares results between direct connection and proxied connection
- Identifies which specific layer is failing (DNS? TLS? IP leak?)
- Gives you actionable insights in seconds

**Perfect for:**
- Developers debugging proxy issues
- DevOps engineers troubleshooting VPN connectivity
- Security teams validating proxy implementations
- Anyone tired of guessing what's broken

## What It Does

ProxyDoctor is a CLI-first tool to:
- Run network checks (DNS resolution, IP detection, TLS certificate validation, port connectivity)
- Compare results between direct connections and proxy-routed connections
- Identify connectivity issues and proxy misconfigurations

## Project Status

**v0.2.0 (Alpha)**
- ✅ Core engine with check registry and dependency DAG
- ✅ CLI with diagnose and list-checks commands
- ✅ HTTP/HTTPS proxy support in `diagnose --proxy`
- ✅ SOCKS4/SOCKS5 proxy support (full protocol implementation, SOCKS4a domain support, SOCKS5 auth per RFC 1929)
- ✅ HTTP server wired to the core engine, with a web GUI at `/` and `/api/checks`, `/api/diagnose` JSON endpoints
- ✅ Web GUI — single-page form at `http://localhost:8080/`, runs a real diagnosis via `POST /api/diagnose` and renders results
- ✅ 4 built-in checks: public_ip, dns_resolve, tls_certificate, port_connectivity
- ✅ `--export json` and `--export markdown` (working)
- ✅ Unit tests (6 passing tests, engine coverage)
- 🔲 More checks (DNS leak, WebRTC leak, geolocation, IP reputation planned)
- 🔲 Plugin system (scaffolding in place, not functional)

## Requirements

- **Go** >= 1.25
- **Git**

## Quick Start

### Setup (first time only)

```bash
git clone https://github.com/francomano/proxydoctor
cd ProxyDoctor
./setup.sh
```

This script will:
- Verify Go installation
- Download and verify dependencies
- Run all tests (6 passing tests)
- Build CLI and server binaries

### Try the Web GUI

```bash
# Start the server
./run.sh server

# Open in browser
open http://localhost:8080
```

Fill in the URL (and optionally a proxy + proxy type), hit "Run diagnosis" — it runs the same `core/engine.DiagnosisOrchestrator` the CLI uses and renders the results as cards.

Two JSON endpoints back the GUI, and can be called directly:

```bash
# List available checks
curl http://localhost:8080/api/checks

# Run a diagnosis
curl -X POST http://localhost:8080/api/diagnose \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","proxy":"socks5://77.245.76.107:1080","proxy_type":"socks5"}'
```

### Use the CLI

```bash
# Get help
./run.sh cli --help

# List available checks
./run.sh cli list-checks

# Run diagnostics (direct connection)
./run.sh cli diagnose --url https://example.com

# Run diagnostics through an HTTP proxy
./run.sh cli diagnose --url https://example.com --proxy http://127.0.0.1:3128 --proxy-type http

# Run diagnostics through a SOCKS5 proxy (with scheme)
./run.sh cli diagnose --url https://example.com --proxy socks5://127.0.0.1:1080 --proxy-type socks5

# Run diagnostics through a SOCKS5 proxy (bare host:port + type)
./run.sh cli diagnose --url https://example.com --proxy 127.0.0.1:1080 --proxy-type socks5

# Export results as JSON
./run.sh cli diagnose --url https://example.com --export json --output report.json
```

### Run the Server

```bash
# Start HTTP server on :8080
./run.sh server
```

### Run Tests

```bash
# Run all tests
./run.sh test

# Or directly
go test -v ./...
```

## File Structure

```
ProxyDoctor/
├── setup.sh              ← One-time setup (install deps, test, build)
├── run.sh                ← Convenience launcher (cli, server, test)
├── cmd/
│   ├── cli/              ← CLI application (diagnose, list-checks)
│   └── server/           ← HTTP API server (minimal, stdlib only) + web GUI
├── core/
│   ├── engine/           ← Orchestration engine (tests included)
│   ├── check/            ← Result types and interfaces (tests included)
│   ├── adapters/         ← Proxy implementations (Direct, HTTP, HTTPS, SOCKS4, SOCKS5)
│   ├── checks/           ← Diagnostic checks (public_ip, dns_resolve, tls_cert, port_scan)
│   └── utils/            ← Shared helpers (proxy URL parsing)
├── go.mod, go.sum        ← Go modules
├── README.md
├── ARCHITECTURE.md
├── CHANGELOG.md
├── NEXT_STEPS.md
└── VERSION
```

## Built-in Checks

| Check | Category | Description |
|---|---|---|
| `public_ip` | network | Detects public IP address via ipify.org, icanhazip.com, ifconfig.me |
| `dns_resolve` | network | Resolves hostname to IP addresses through the current connection |
| `tls_certificate` | tls | Validates TLS certificate (issuer, expiry, cipher suite, TLS version) |
| `port_connectivity` | network | Tests TCP connectivity to ports 80, 443, 8080, 8443 |

## Proxy Input Formats

The CLI and GUI accept proxy URLs in multiple formats:

| Format | Example | Notes |
|---|---|---|
| `scheme://host:port` | `socks5://77.245.76.107:1080` | Auto-detects type from scheme |
| `host:port` + type | `77.245.76.107:1080` + `--proxy-type socks5` | Requires explicit type |
| `host` + type | `77.245.76.107` + `--proxy-type socks5` | Uses default port (1080 for SOCKS, 8080 for HTTP) |
| With auth | `socks5://user:pass@host:port` | Credentials extracted from URL |

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

- 🐛 [Report a Bug](https://github.com/francomano/ProxyDoctor/issues/new?template=bug_report.md)
- ✨ [Request a Feature](https://github.com/francomano/ProxyDoctor/issues/new?template=feature_request.md)
- 💡 [Good First Issues](https://github.com/francomano/ProxyDoctor/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

## Known Issues / Limitations

- `--export html` flag is accepted but silently falls back to text output (html formatter not implemented).
- `--compare` flag is declared but not yet wired into the diagnosis execution.
- DNS leak, WebRTC leak, geolocation, and IP reputation checks are planned but not yet implemented.
