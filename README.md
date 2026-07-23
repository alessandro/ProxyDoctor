# ProxyDoctor — Proxy Diagnostics

**A minimal, testable Go tool for running network diagnostics through proxies**

[![License: GPL-3.0](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)

![Logo](images/proxydoctor-logo.png)

> **Built with a big contribution of Copilot Auto Mode and Claude Sonnet 5** — This project was developed and documented with GitHub Copilot and Claude.

## 🤔 Why ProxyDoctor?

**The Problem:**
You're behind a corporate proxy or using a VPN, and things break mysteriously:
- Sites don't load (proxy misconfiguration?)
- Your IP leaks even through a "private" proxy
- You can't tell if the problem is your proxy, DNS, or the site itself

**The Solution:**
ProxyDoctor is a **lightweight diagnostic tool** that:
- Runs network checks through any proxy (HTTP, SOCKS4/5)
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
- Run network checks (DNS, IP detection, port connectivity, TLS)
- Compare results between direct connections and proxy-routed connections
- Identify connectivity issues and proxy misconfigurations

## Project Status

**v0.1.0 (Alpha)**
- ✅ Core engine with check registry and dependency DAG
- ✅ CLI with diagnose and list-checks commands
- ✅ HTTP/HTTPS proxy support in `diagnose --proxy` (URL is properly parsed: scheme, host, port, credentials)
- ✅ HTTP server (`cmd/server`) — now wired to the real core engine (`core/engine.DiagnosisOrchestrator`), same code path as the CLI
- ✅ Web GUI — single-page form at `http://localhost:8080/`, runs a real diagnosis via `POST /api/diagnose` and renders results
- 🔲 SOCKS4/SOCKS5 proxy support — flag/option exists but adapters are not implemented yet (`not yet implemented` error)
- ✅ Unit tests (6 passing tests, 100% coverage of engine)
- 🔲 Plugin system (scaffolding in place, not functional)
- 🔲 `--export json` / `--export markdown` — flags exist but formatters are unimplemented (TODO stubs)

## Requirements

- **Go** >= 1.21
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
- Run all tests (6 passing tests ✅)
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
  -d '{"url":"https://example.com","proxy":"http://127.0.0.1:3128","proxy_type":"http"}'
```

### Use the CLI

```bash
# Get help
./run.sh cli --help

# List available checks
./run.sh cli list-checks

# Run diagnostics (--url is required)
./run.sh cli diagnose --url https://example.com

# Run diagnostics through an HTTP proxy
./run.sh cli diagnose --url https://example.com --proxy http://127.0.0.1:3128 --proxy-type http
```

### Run the Server

```bash
# Start HTTP server on :8080
./run.sh server

# In another terminal, test the endpoint
curl http://localhost:8080/api/check/public-ip
# Returns: {"ip":"1.2.3.4"}
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
│   └── server/           ← HTTP API server (minimal, stdlib only)
├── core/
│   ├── engine/           ← Orchestration engine (tests included)
│   ├── check/            ← Result types and interfaces (tests included)
│   ├── adapters/         ← Proxy implementations (HTTP, SOCKS)
│   └── checks/           ← Diagnostic checks (public_ip, dns_leak, etc.)
├── gui/                  ← Optional Tauri/React GUI
└── go.mod, go.sum        ← Go modules (Cobra, Viper)
```

What the main files do (very short)
- `cmd/cli/` — CLI entrypoint and commands (diagnose, list, version).
- `cmd/server/` — tiny HTTP server (stdlib only) wired to the core engine; serves the web GUI at `/` and the `/api/checks`, `/api/diagnose` JSON endpoints.
- `core/` — engine, adapters and checks implementations (testable packages).
- `core/checks/public_ip/` — example check used by the server.
- `ARCHITECTURE.md` — minimal architecture diagram.
- `VERSION`, `CHANGELOG.md`, `NEXT_STEPS.md` — versioning and high-level next steps.

## Known Issues / Limitations

- SOCKS4/SOCKS5 proxy types are accepted by the CLI and GUI but the underlying adapters (`core/adapters/socks.go`) are stubs that return `not yet implemented`.
- `--export json` and `--export markdown` are not implemented yet (only `text` output works in the CLI; the GUI/API always returns full JSON).
- Only one check (`public_ip`) is implemented; all others under `core/checks/` in `ARCHITECTURE.md` are planned, not built.