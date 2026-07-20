# ProxyDoctor — Proxy Diagnostics

**A minimal, testable Go tool for running network diagnostics through proxies**

[![License: GPL-3.0](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)

> **Built with Copilot Auto Mode** — This project was simplified, tested, and documented with GitHub Copilot.

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
- ✅ HTTP/HTTPS proxy support in `diagnose --proxy` (fixed: URL is now actually parsed, was previously hardcoded)
- 🔲 SOCKS4/SOCKS5 proxy support — flag exists but adapters are not implemented yet (`not yet implemented` error)
- ⚠️ HTTP server (`cmd/server`) — minimal placeholder, **not connected to the core engine**. Only exposes `/api/check/public-ip` via a hardcoded direct request. Does not run real diagnoses and does not serve any page on `/`.
- 🔲 Web GUI — not implemented. There is no HTML page served at `http://localhost:8080`; hitting `/` returns 404. Only the JSON endpoint above exists.
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

### Try the HTTP API (no GUI yet)

There is currently no web page — only a single JSON endpoint. `http://localhost:8080/` will 404.

```bash
# Start the server
./run.sh server

# In another terminal, hit the actual endpoint (note the /api/... path)
curl http://localhost:8080/api/check/public-ip
# Returns: {"ip":"1.2.3.4"}
```

Note: this server does **not** use the core engine — it's a separate, minimal placeholder. Real diagnoses (with proxy support, multiple checks, etc.) currently only work via the CLI below.

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
- `cmd/server/` — tiny HTTP server that runs a single check (`public_ip`) and returns JSON.
- `core/` — engine, adapters and checks implementations (testable packages).
- `core/checks/public_ip/` — example check used by the server.
- `ARCHITECTURE.md` — minimal architecture diagram.
- `VERSION`, `CHANGELOG.md`, `NEXT_STEPS.md` — versioning and high-level next steps.

## Known Issues / Limitations

- `cmd/server` is disconnected from `core/` — it's a standalone placeholder, not the real diagnosis engine exposed over HTTP.
- SOCKS4/SOCKS5 proxy types are accepted by the CLI but the underlying adapters (`core/adapters/socks.go`) are stubs that return `not yet implemented`.
- `--export json` and `--export markdown` are not implemented yet (only `text` output works).
- Only one check (`public_ip`) is implemented; all others under `core/checks/` in `ARCHITECTURE.md` are planned, not built.