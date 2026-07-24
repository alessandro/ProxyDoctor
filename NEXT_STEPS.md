# Next Steps — Development Roadmap

> **Built with AI assistance** — Contributions from GitHub Copilot, OpenCode (free open-source models), and Gemini Pro.

## Current Status (v0.2.0 Alpha)

✅ **Complete:**
- Core engine (check registry, dependency DAG, orchestrator)
- CLI with diagnose and list-checks commands
- HTTP/HTTPS proxy support in `diagnose --proxy`
- SOCKS4/SOCKS5 proxy support — full protocol implementations in `core/adapters/socks.go`
  - SOCKS4/4a: custom dialer with domain name support via `0.0.0.1` extension
  - SOCKS5: `golang.org/x/net/proxy` with manual fallback, IPv4/IPv6/domain targets, username/password auth (RFC 1929)
  - All `NetworkAdapter` methods implemented: HTTP requests, redirects, DNS, port testing, TLS detection, public IP
- HTTP server wired to the core engine, with a web GUI at `/` and `/api/checks`, `/api/diagnose` JSON endpoints
- 4 built-in checks: `public_ip`, `dns_resolve`, `tls_certificate`, `port_connectivity`
- `--export json` and `--export markdown` (working)
- Proxy URL parsing: supports `scheme://host:port`, bare `host:port` (with explicit type), and bare `host` (with explicit type)
- Unit tests (6/6 passing, engine coverage)
- Documentation (README, ARCHITECTURE, CHANGELOG, this file)

⚠️ **Incomplete / remaining gaps:**
- `--export html` flag accepted but falls back to text (formatter not implemented)
- `--compare` flag declared but not wired into execution
- No DNS leak, WebRTC leak, geolocation, or IP reputation checks yet

## Immediate Next Steps (v0.3.0)

1. **Add leak detection checks** (1-2 weeks)
   - Implement DNS leak detection (`core/checks/dns_leak/check.go`)
   - Implement WebRTC leak detection (`core/checks/webrtc_leak/check.go`)
   - Implement IPv6 leak detection
   - Register them in CLI and server

2. **CLI improvements** (3-5 days)
   - Implement `--export html` formatter
   - Add `--timeout` parameter for diagnosis (currently hardcoded to 30s)
   - Add `--checks` flag to select specific checks
   - Wire the `--compare` flag (run against direct + proxy, show diff)
   - Improve error messages and logging

3. **More checks** (1-2 weeks)
   - TLS certificate chain validation
   - Redirect analysis (track 301/302 chains)
   - HTTP security headers analysis
   - Geolocation (MaxMind GeoIP2 or free alternative)
   - IP reputation (AbuseIPDB)

4. **GUI/API enhancements** (3-5 days)
   - Show per-check progress instead of waiting for the full report
   - Add a "compare direct vs proxy" toggle in the form
   - Persist recent diagnosis results client-side (session only, no backend storage)

## Medium-term (v0.4.0 - v1.0)

5. **GUI Integration** (2-3 weeks)
   - Wire `gui/` (Tauri/React) to call HTTP server endpoints
   - Display check results with real-time updates
   - Add result export (JSON, PDF, HTML)

6. **Plugin System** (2-3 weeks)
   - Publish plugin SDK documentation
   - Create example plugin template
   - Test community contributions

7. **Performance & Optimization** (1-2 weeks)
   - Profile and optimize check execution
   - Add result caching layer
   - Implement parallel check execution where safe

## Long-term (Future)

- VPN routing support
- Browser automation checks (Playwright)
- Mobile app (React Native)
- Cloud-based result storage and comparison
- Community check registry

## Development Workflow

```bash
# Setup environment
./setup.sh

# Run tests after every change
go test -v ./...

# Build binaries
go build -o bin/proxydoctor ./cmd/cli
go build -o bin/proxydoctor-server ./cmd/server

# Or use convenience scripts
./run.sh cli --help
./run.sh server

# Tag releases
git tag -a v0.2.0 -m "v0.2.0: SOCKS4/SOCKS5, new checks, proxy UX"
git push --tags
```

## Release Checklist

Before tagging a new version:
- [ ] All tests pass (`go test ./...`)
- [ ] Binaries build without warnings (`go build ./cmd/cli` and `./cmd/server`)
- [ ] Update `VERSION` file
- [ ] Update `CHANGELOG.md` with new features
- [ ] Update `README.md` if needed
- [ ] Tag with `git tag -a vX.Y.Z -m "vX.Y.Z: description"`
- [ ] Push tags with `git push --tags`

## Repository Files

- `README.md` — project overview, quick start, check list, proxy input formats.
- `ARCHITECTURE.md` — architecture diagram, package structure, planned checks.
- `VERSION` — current project version.
- `CHANGELOG.md` — release notes.
- `NEXT_STEPS.md` — this file.
- `go.mod` — Go module and dependency list.
- `cmd/cli/` — CLI entry point and commands.
- `cmd/server/` — HTTP server + web GUI.
- `core/` — engine, adapters, checks, and utils (testable packages).

## How to Run Tests

```bash
# from project root
go test ./...

# build CLI
go build -o bin/proxydoctor ./cmd/cli

# run CLI help
./bin/proxydoctor diagnose --help
```
