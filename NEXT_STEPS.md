# Next Steps — Development Roadmap

> **Built with AI assistance** — Contributions from GitHub Copilot, OpenCode (free open-source models), and Gemini Pro.

## Current Status (v0.2.0 Alpha)

✅ **Complete:**
- Core engine (check registry, dependency DAG, orchestrator)
- CLI with diagnose and list-checks commands
- HTTP/HTTPS proxy support in `diagnose --proxy` (fixed: was previously ignoring the flag entirely)
- **SOCKS4/SOCKS5 proxy support** — full protocol implementations in `core/adapters/socks.go`
  - SOCKS4/4a: custom dialer with domain name support via `0.0.0.1` extension
  - SOCKS5: `golang.org/x/net/proxy` with manual fallback, IPv4/IPv6/domain targets, username/password auth (RFC 1929)
  - All `NetworkAdapter` methods implemented: HTTP requests, redirects, DNS, port testing, TLS detection, public IP
- HTTP server wired to the core engine, with a web GUI at `/` and `/api/checks`, `/api/diagnose` JSON endpoints
- `--export json` and `--export markdown` (working)
- Unit tests (6/6 passing, engine coverage)
- Documentation (README, ARCHITECTURE, this file)

⚠️ **Incomplete / remaining gaps:**
- Only one check (`public_ip`) is registered; DNS leak, WebRTC leak, TLS checks not yet implemented
- `--export html` flag accepted but falls back to text (formatter not implemented)
- `--compare` flag declared but not wired into execution

## Immediate Next Steps (v0.2.0)

1. **Add more checks** (1-2 weeks)
   - Implement DNS leak detection (`core/checks/dns_leak/check.go`)
   - Implement WebRTC leak detection (`core/checks/webrtc_leak/check.go`)
   - Add TLS/certificate validation checks
   - Register them in both `cmd/cli/commands/diagnose.go` and `cmd/server/main.go`'s `newRegistry()` — currently only `public_ip` is registered in either

2. **CLI improvements** (3-5 days)
   - Implement `--export html` formatter
   - Add `--timeout` parameter for diagnosis (currently hardcoded to 30s)
   - Add `--checks` flag to select specific checks
   - Wire the `--compare` flag (run against direct + proxy, show diff)
   - Improve error messages and logging

3. **GUI/API enhancements** (3-5 days)
   - Show per-check progress instead of waiting for the full report
   - Add a "compare direct vs proxy" toggle in the form
   - Persist recent diagnosis results client-side (session only, no backend storage)

## Medium-term (v0.3.0 - v1.0)

4. **GUI Integration** (2-3 weeks)
   - Wire `gui/` (Tauri/React) to call HTTP server endpoints
   - Display check results with real-time updates
   - Add result export (JSON, PDF, HTML)

5. **Plugin System** (2-3 weeks)
   - Publish plugin SDK documentation
   - Create example plugin template
   - Test community contributions

6. **Performance & Optimization** (1-2 weeks)
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
git tag -a v0.2.0 -m "v0.2.0: SOCKS4/SOCKS5 proxy support"
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

Brief description of files left in the repository

- `README.md` — short project overview and quick start commands.
- `ARCHITECTURE.md` — architecture diagram, package structure, and what to keep testable.
- `VERSION` — current project version (simple single-line file).
- `CHANGELOG.md` — release notes template.
- `go.mod` — Go module and dependency list.
- `cmd/cli/` — CLI entry point and commands. This is where `go build ./cmd/cli` produces the CLI binary.
- `core/` — core engine, adapters and check implementations. Testable package with unit tests under `core/test`.

How to run tests locally

```bash
# from project root
go test ./...

# build CLI
go build -o bin/proxydoctor ./cmd/cli

# run CLI help
./bin/proxydoctor diagnose --help
```
