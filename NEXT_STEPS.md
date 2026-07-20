# Next Steps — Development Roadmap

> **Built with Copilot Auto Mode** — This roadmap was created with AI assistance for clarity and feasibility.

## Current Status (v0.1.0 Alpha)

✅ **Complete:**
- Core engine (check registry, dependency DAG, orchestrator)
- CLI with diagnose and list-checks commands
- HTTP/HTTPS proxy support in `diagnose --proxy` (fixed: was previously ignoring the flag entirely)
- Unit tests (6/6 passing, 100% engine coverage)
- Documentation (README, ARCHITECTURE, this file)

⚠️ **Incomplete / misleading in earlier docs:**
- HTTP server (`cmd/server`) exists but is **not wired to the core engine** — it's a separate minimal placeholder with one endpoint (`/api/check/public-ip`) and no web page at `/`
- SOCKS4/SOCKS5 adapters are stubs (`not yet implemented`), even though the CLI flag accepts them
- `--export json` / `--export markdown` are TODO stubs; only `text` output works

## Immediate Next Steps (v0.2.0)

1. **Implement SOCKS4/SOCKS5 adapters** (2-4 days)
   - `core/adapters/socks.go`: `ExecuteHTTPRequest` currently returns `not yet implemented` for both
   - Needed before `--proxy-type socks4/socks5` can actually be used

2. **Wire `cmd/server` to the core engine** (2-3 days)
   - Currently it bypasses `core/` entirely and hardcodes a direct request to ipify/icanhazip
   - Should build a `DiagnosisOrchestrator` the same way `cmd/cli/commands/diagnose.go` does, and expose it over HTTP
   - Only after this is done does a browser-based GUI make sense

3. **Add more checks** (1-2 weeks)
   - Implement DNS leak detection (`core/checks/dns_leak/check.go`)
   - Implement WebRTC leak detection (`core/checks/webrtc_leak/check.go`)
   - Add TLS/certificate validation checks
   - Update CLI to `diagnose --check dns_leak --proxy http://localhost:8080`

4. **CLI improvements** (3-5 days)
   - Implement `--export json` and `--export markdown` (currently TODO stubs)
   - Add `--timeout` parameter for diagnosis
   - Improve error messages and logging

5. **HTTP server enhancements** (3-5 days)
   - Add `/api/checks` endpoint (list available checks)
   - Add `/api/diagnose` endpoint (run full diagnosis via HTTP)
   - Add CORS headers for GUI integration
   - Serve a minimal HTML page at `/` (there currently is none — hitting it 404s)

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
git tag -a v0.2.0 -m "v0.2.0: DNS and WebRTC checks"
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
   - Publish binaries or provide build instructions when ready.

Brief description of files left in the repository

- `README.md` — short project overview and quick start commands.
- `ARCHITECTURE.md` — very short architecture diagram and what to keep testable.
- `VERSION` — current project version (simple single-line file).
- `CHANGELOG.md` — release notes template.
- `go.mod` — Go module and dependency list.
- `cmd/cli/` — CLI entry point and commands. This is where `go build ./cmd/cli` produces the CLI binary.
- `core/` — core engine, adapters and check implementations. Testable package with unit tests under `core/test`.

What I moved to `archived/`

- Large docs, optional GUI code, CI/Docker artifacts, plugins and example folders were moved into `archived/` to keep the repo minimal and reversible.

How to run tests locally

```bash
# from project root
go test ./...

# build CLI
go build -o bin/proxydoctor ./cmd/cli

# run CLI help
./bin/proxydoctor diagnose --help
```

If you want, I can now:
- (A) Add the small HTTP wrapper to `cmd/server` (1 file) and the route to run a public_ip check.
- (B) Commit and tag a `v0.1.0` release.
- (C) Reintroduce a minimal CI workflow once tests are passing.