# ProxyDoctor — Architecture

> **Built with Copilot Auto Mode** — Simplified and tested with GitHub Copilot.

## Overview

```
CLI (cmd/cli)
    ↓
Diagnosis Engine (core/engine)
    ├─ CheckRegistry → stores/resolves checks
    ├─ DependencyDAG → sorts checks by dependency
    └─ DiagnosisOrchestrator → executes checks in order
    ↓
Network Adapters (core/adapters)
    ├─ Direct adapter (no proxy)
    ├─ HTTP/HTTPS proxy adapters
    ├─ SOCKS4/SOCKS5 adapters
    └─ Adapter factory for dynamic creation
    ↓
Checks (core/checks/*)
    ├─ public_ip — fetch and return public IP
    ├─ dns_leak — (scaffolding)
    ├─ webrtc_leak — (scaffolding)
    └─ ... (more checks to come)
    ↓
Results & Reporting (core/check)
    └─ CheckResult — fluent builder, status, confidence, evidence
```

## Package Structure

```
core/
├── check/              ✅ Type definitions (CheckResult, Adapter interface)
├── engine/             ✅ Orchestration (Registry, DAG, Orchestrator)
├── adapters/           ⚠️ Direct + HTTP/HTTPS proxy implemented; SOCKS4/SOCKS5 are stubs (not yet implemented)
├── checks/
│   └── public_ip/      ✅ Public IP check (working)
│       └── check.go    → implements Checker interface
└── utils/              → helper functions

cmd/
├── cli/                ✅ CLI application (Cobra + Viper)
│   ├── commands/       ✅ diagnose, list-checks, version
│   └── main.go         ✅ CLI entry point
└── server/             ⚠️ HTTP server (stdlib only) — NOT wired to core/, standalone placeholder
    └── main.go         → /api/check/public-ip endpoint only, no page served at "/"

gui/                   📦 Tauri/React (archived, optional, not started)
```

## Key Design Decisions

- **No external HTTP libraries** — `cmd/server` uses stdlib only to avoid dependency bloat
- **Fluent builder for results** — Type-safe, readable check result construction
- **Dependency DAG for checks** — Automatic ordering respects implicit dependencies
- **Interface-based adapters** — Easy to add new proxy types (VPN, Tor, etc.)
- **All tests in same package** — Easier to maintain and refactor
- **Copilot-assisted development** — This architecture was refined with AI assistance for clarity and maintainability
```

#### 3.1.4 Modulo `checks` (Built-in checks)
```
core/checks/
├── public_ip/
│   ├── check.go
│   └── resolver.go (multiple services: ipify, ifconfig.me, etc.)
│
├── geolocation/
│   ├── check.go
│   ├── maxmind.go (MaxMind GeoIP2)
│   └── fallback_resolver.go
│
├── asn/
│   ├── check.go
│   └── whois_resolver.go
│
├── dns_leak/
│   ├── check.go
│   └── detector.go (compare DNS queries with direct vs proxy)
│
├── webrtc_leak/
│   ├── check.go
│   ├── ice_gathering.go
│   └── stun_server_monitor.go
│
├── ipv6_leak/
│   ├── check.go
│   └── ipv6_resolver.go
│
├── http2_support/
│   ├── check.go
│   └── protocol_detector.go
│
├── http3_quic_support/
│   ├── check.go
│   └── quic_connector.go
│
├── tls_handshake/
│   ├── check.go
│   ├── cipher_suite_analyzer.go
│   └── certificate_validator.go
│
├── certificates/
│   ├── check.go
│   ├── chain_validator.go
│   └── revocation_checker.go (OCSP, CRL)
│
├── redirects/
│   ├── check.go
│   └── redirect_follower.go (track 301, 302, etc.)
│
├── headers/
│   ├── check.go
│   ├── security_headers_analyzer.go
│   └── tracking_headers_detector.go
│
├── cookies/
│   ├── check.go
│   └── cookie_analyzer.go
│
├── browser_fingerprint/
│   ├── check.go
│   ├── canvas_fingerprint.go
│   ├── webgl_fingerprint.go
│   ├── navigator_properties.go
│   └── font_detection.go
│
├── streaming_compatibility/
│   ├── check.go
│   └── cdn_detector.go (Netflix, YouTube, etc.)
│
└── ip_reputation/
    ├── check.go
    └── reputation_checker.go (Abuseipdb, Shodan, etc.)
```

#### 3.1.5 Modulo `utilities`
```
core/utilities/
├── http_client.go (with proxy support)
├── dns_resolver.go (with leak detection)
├── tls_parser.go
├── geolocation_resolver.go
├── external_api_client.go (MaxMind, AbuseIPDB, etc.)
├── config_manager.go
├── logger.go
└── cache.go (in-memory, optional Redis)
```

### 3.2 Plugin System (`/plugins`)

```
plugins/
├── README.md (plugin development guide)
├── sdk/
│   └── check_sdk.{go,ts,rs} (template per linguaggio)
│
├── community-checks/
│   ├── tor-exit-node-check/
│   │   ├── plugin.json
│   │   ├── check.go
│   │   ├── test.go
│   │   └── README.md
│   │
│   ├── cloudflare-cdn-check/
│   │   ├── plugin.json
│   │   ├── check.go
│   │   └── test.go
│   │
│   └── bittorrent-leak-check/
│       ├── plugin.json
│       ├── check.go
│       └── test.go
│
└── examples/
    └── simple-port-check/
        ├── plugin.json
        ├── check.go
        └── test.go
```

**Plugin Contract:**
```
plugin.json
{
  "id": "tor-exit-node",
  "name": "Tor Exit Node Detector",
  "version": "1.0.0",
  "author": "community",
  "language": "go",
  "api_version": "1.0",
  "dependencies": [],
  "entry_point": "check.go",
  "description": "Detects if the current IP is a known Tor exit node"
}
```

### 3.3 CLI (`/cmd/cli`)

**Stack consigliato:** Rust (clap) o Go (cobra), oppure Node.js (yargs/commander)

```
cmd/cli/
├── main.go (or .rs, .ts)
├── commands/
│   ├── diagnose.go (main command)
│   │   ├── --url string (required)
│   │   ├── --proxy string (optional, es: http://localhost:8080)
│   │   ├── --proxy-type string (auto, http, https, socks4, socks5)
│   │   ├── --checks string (comma-separated IDs, or "all")
│   │   ├── --compare (run against direct connection too)
│   │   ├── --export string (json, html, csv, md)
│   │   ├── --output string (file path)
│   │   ├── --json (shorthand for --export json)
│   │   ├── --verbose (-v, -vv)
│   │   └── --config string (path to config file)
│   │
│   ├── list-checks.go
│   │   └── Lists all available checks with descriptions
│   │
│   ├── plugin.go
│   │   ├── install <plugin-url>
│   │   ├── list
│   │   ├── uninstall <plugin-id>
│   │   └── verify <plugin-id>
│   │
│   ├── config.go
│   │   ├── show
│   │   ├── set <key> <value>
│   │   ├── init (guided setup)
│   │   └── reset
│   │
│   ├── version.go
│   ├── help.go
│   └── completion.go (bash, zsh, fish)
│
├── formatters/
│   ├── table_formatter.go (colorized table output)
│   ├── json_formatter.go
│   ├── markdown_formatter.go
│   └── html_formatter.go
│
├── styles/
│   └── colors.go (ANSI colors for terminal)
│
└── config/
    └── default_config.yaml
```

**Esempio di utilizzo CLI:**
```bash
# Basic usage
proxyctl diagnose --url https://example.com

# With proxy
proxyctl diagnose --url https://example.com --proxy http://localhost:8080

# Compare direct vs proxy
proxyctl diagnose --url https://example.com --proxy socks5://localhost:1080 --compare

# Specific checks only
proxyctl diagnose --url https://example.com --checks "public_ip,dns_leak,webrtc_leak"

# Export results
proxyctl diagnose --url https://example.com --export json --output report.json

# List available checks
proxyctl list-checks

# Install community plugin
proxyctl plugin install https://github.com/community/tor-check
```

### 3.4 GUI (`/cmd/gui`)

**Stack consigliato:**
- **Tauri + React/Vue** (Rust backend, nativo, leggero, ~150MB)
- **Electron + React** (Node.js, più memoria, più comunità)
- **Qt/PySide** (Python, meno moderno ma stabile)

**Consiglio per MVP:** Tauri + React

```
cmd/gui/
├── src-tauri/ (Rust backend for Tauri)
│   ├── src/
│   │   ├── main.rs
│   │   ├── commands/ (Tauri IPC commands)
│   │   │   ├── diagnose.rs (invokes core engine)
│   │   │   ├── list_checks.rs
│   │   │   ├── plugin_manager.rs
│   │   │   └── config.rs
│   │   └── utils/
│   │       └── logger.rs
│   │
│   ├── tauri.conf.json
│   └── Cargo.toml
│
└── src/ (React frontend)
    ├── components/
    │   ├── DiagnosisForm.jsx
    │   │   ├── URLInput
    │   │   ├── ProxySelector
    │   │   └── CheckOptions
    │   │
    │   ├── ResultsView.jsx
    │   │   ├── CheckResultCard
    │   │   ├── SeverityFilter
    │   │   ├── DetailedView
    │   │   └── ExportButton
    │   │
    │   ├── ComparisonView.jsx (direct vs proxy)
    │   └── SettingsPanel.jsx
    │
    ├── hooks/
    │   ├── useDiagnosis.js
    │   ├── useResults.js
    │   └── usePlugins.js
    │
    ├── styles/
    │   ├── global.css
    │   ├── theme.css
    │   └── dark-mode.css
    │
    ├── utils/
    │   ├── api.js (Tauri IPC client)
    │   └── formatter.js
    │
    ├── App.jsx
    ├── index.html
    └── package.json
```

**Features GUI:**
- ✅ Real-time progress indicator
- ✅ Dark mode
- ✅ Detailed view with evidence
- ✅ Side-by-side comparison (direct vs proxy)
- ✅ Export results (JSON, HTML, PDF)
- ✅ History of diagnostics
- ✅ Plugin management UI
- ✅ Configuration wizard

---

## 4. Flusso Completo di una Diagnosi

### 4.1 Scenario: User esegue diagnosi

```
┌─────────────────────────────────────────────────────────────────┐
│ User Input (CLI/GUI)                                             │
│ ├─ URL: https://example.com                                      │
│ ├─ Proxy: http://proxy.local:8080                               │
│ └─ Checks: [all] or [custom list]                               │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ DiagnosisEngine.Start()                                          │
├─ Parse & validate input                                         │
├─ Load proxy config (resolve hostname, test connectivity)        │
├─ Initialize network adapters (direct + proxy)                   │
├─ Load all checks from registry                                  │
├─ Build dependency DAG                                           │
└─ Initialize execution context                                   │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ Check Execution Phase (Parallel by DAG)                          │
├─ Check 1: Public IP (no dependencies)         ──┐               │
├─ Check 2: Geolocation (depends on Public IP) ──┤               │
├─ Check 3: DNS Leak (independent)             ──┤               │
├─ Check 4: WebRTC Leak (independent)          ──┤               │
├─ Check 5: TLS Handshake (independent)        ──┼─ PARALLEL    │
├─ Check 6: HTTP/2 Support (independent)       ──┤               │
├─ Check 7: Redirect Follow (independent)      ──┤               │
├─ Check 8: IP Reputation (depends on IP)      ──┤               │
└─ ...                                           ──┘               │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ Each Check Execution                                             │
├─ Execute check.Execute(ctx)                                     │
│   ├─ Use appropriate network adapter (direct/proxy)             │
│   ├─ Collect evidence                                           │
│   ├─ Populate CheckResult struct                                │
│   ├─ Handle errors gracefully (set Status: failed/error)        │
│   └─ Return CheckResult                                         │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ Result Aggregation                                               │
├─ Collect all CheckResult objects                                │
├─ Sort by severity (critical, warning, info)                     │
├─ Identify patterns & cross-check correlations                   │
│   ├─ E.g.: "If DNS leak + WebRTC leak → likely no encryption" │
│   └─ E.g.: "If redirect to different IP → possible MITM"      │
├─ Generate executive summary                                     │
└─ Create final DiagnosisReport                                   │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ Report Generation & Output                                       │
├─ Formatter selection (JSON, HTML, Markdown, CSV)                │
├─ Add metadata (timestamp, duration, user agent)                 │
├─ Save or stream to user                                         │
└─ If GUI: render interactive results                             │
└─ If CLI: display colorized output                               │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│ User Output (CLI/GUI)                                            │
├─ Summary of findings                                             │
├─ Critical issues highlighted                                    │
├─ Recommended actions                                            │
└─ Option to export or share results                              │
└─────────────────────────────────────────────────────────────────┘
```

### 4.2 Data Flow Tra Moduli

```
CLI/GUI Input
    │
    ▼
┌─────────────────────────┐
│ DiagnosisOrchestrator   │
├─ parseInput()           │
├─ validateConfig()       │
└─ startDiagnosis()       │
    │
    ├─────────────────────────────────────────┐
    │                                         │
    ▼                                         ▼
┌──────────────────┐              ┌──────────────────────┐
│ DirectAdapter    │              │ ProxyAdapter         │
│ ─────────────    │              │ ─────────────────    │
│ .DNS()           │              │ .HTTP()              │
│ .TLS()           │              │ .DNS() [monitored]   │
│ .HTTP()          │              │ .TLS()               │
└────────┬─────────┘              └──────────┬───────────┘
         │                                   │
         └───────────┬───────────────────────┘
                     │
                     ▼
            ┌──────────────────┐
            │ Check Registry   │
            ├─ checks[]       │
            │ ├─ PublicIPCheck│
            │ ├─ DNSLeakCheck │
            │ ├─ WebRTCCheck  │
            │ └─ ...          │
            └────────┬─────────┘
                     │
         ┌───────────┴───────────┐
         │                       │
         ▼                       ▼
    ┌─────────────┐          ┌─────────────┐
    │ Check 1     │          │ Check N     │
    │ Execute()   │          │ Execute()   │
    │ ↓           │          │ ↓           │
    │ CheckResult │          │ CheckResult │
    └────┬────────┘          └────┬────────┘
         │                        │
         └────────┬───────────────┘
                  │
                  ▼
         ┌──────────────────┐
         │ ResultAggregator │
         ├─ correlate()     │
         ├─ prioritize()    │
         ├─ summarize()     │
         └─ generate()      │
                  │
                  ▼
         ┌──────────────────┐
         │ DiagnosisReport  │
         ├─ results[]      │
         ├─ summary        │
         ├─ metadata       │
         └─ exported_at    │
                  │
                  ▼
    ┌────────────┴────────────┐
    │                         │
    ▼                         ▼
┌──────────────┐        ┌──────────────┐
│ CLI Output   │        │ GUI Display  │
│ (formatted)  │        │ (interactive)│
└──────────────┘        └──────────────┘
```

---

## 5. API Interne Tra Moduli

### 5.1 Core → Checks

```
interface Checker {
    // Metadata
    ID() string
    Name() string
    Description() string
    Category() CheckCategory
    
    // Dependencies
    DependsOn() []string
    
    // Execution
    Execute(ctx ExecutionContext) CheckResult
}

interface ExecutionContext {
    URL string
    ProxyConfig ProxyConfig
    DirectAdapter NetworkAdapter
    ProxyAdapter NetworkAdapter
    Metadata map[string]interface{}
    GetSharedData(key string) interface{}
    SetSharedData(key string, value interface{})
    GetTimeout() Duration
    IsCancelled() bool
}

interface CheckResult {
    ID string
    Category CheckCategory
    Status Status // passed, failed, skipped, error
    Severity Severity // info, warning, critical
    Confidence float64 // 0.0-1.0
    Evidence map[string]interface{}
    Explanation string
    ProbableCauses []string
    SuggestedActions []string
    References []string
    ExecutionTime Duration
    Timestamp time.Time
}
```

### 5.2 DiagnosisEngine → Checks

```
// Registry pattern
type CheckRegistry {
    LoadChecks(pluginDir string) error
    GetCheck(id string) Checker
    ListChecks() []Checker
    RegisterCheck(check Checker) error
    UnregisterCheck(id string) error
    GetDependencyGraph() DAG
}

// Execution
type DiagnosisEngine {
    LoadChecks() error
    ValidateInput(url string, proxy string) error
    ExecuteChecks(checkIDs []string) ([]CheckResult, error)
    ExecuteChecksDependencyOrder() ([]CheckResult, error)
    CancelExecution() error
}
```

### 5.3 Checks → NetworkAdapters

```
interface NetworkAdapter {
    // HTTP
    ExecuteHTTPRequest(req *http.Request) (*http.Response, error)
    FollowRedirects(url string, maxRedirects int) ([]RedirectStep, error)
    
    // DNS
    ResolveDNS(hostname string) ([]net.IP, error)
    ResolveDNSWithDetails() (DNSResponse, error)
    
    // Network
    GetPublicIP() (string, error)
    TestPort(host string, port int, timeout Duration) (bool, error)
    
    // TLS
    GetTLSCertificate(hostname string) (*x509.Certificate, error)
    GetTLSChain(hostname string) ([]*x509.Certificate, error)
    GetTLSCipherSuite(hostname string) (string, error)
    GetTLSVersion(hostname string) (string, error)
    
    // WebRTC (browser-based)
    ExecuteJavaScript(script string) (interface{}, error)
    GetBrowserContext() BrowserContext
    
    // Browser (headless)
    OpenBrowser(url string) BrowserTab
    NavigateTo(url string) error
    ExecuteScript(script string) (interface{}, error)
    GetCookies() []Cookie
    GetHeaders() http.Header
}
```

### 5.4 CLI/GUI → Core

```
// CLI Interface
type CoreAPIClient {
    Diagnose(req DiagnosisRequest) (DiagnosisReport, error)
    ListChecks() ([]CheckMetadata, error)
    ValidateURL(url string) error
    ValidateProxy(proxy string) error
}

type DiagnosisRequest {
    URL string
    ProxyConfig ProxyConfig
    CheckIDs []string // or [] for all
    Compare bool // run against direct too
    Timeout Duration
    ExportFormat string // json, html, csv, md
}

type DiagnosisReport {
    ID string
    RequestMetadata RequestMetadata
    Results []CheckResult
    Summary ExecutiveSummary
    Correlations []Finding
    GeneratedAt time.Time
    Duration Duration
}
```

---

## 6. Tecnologie Consigliate

### 6.1 Backend (Core Engine)

| Aspetto | Consigliato | Alternativa |
|---------|------------|-------------|
| **Linguaggio** | Go | Rust, Node.js/TypeScript |
| **HTTP Client** | `net/http`, `resty` | `reqwest` (Rust) |
| **DNS** | `net.LookupHost()` | `trust-dns` (Rust) |
| **TLS** | `crypto/tls` | `rustls` |
| **Async/Concurrency** | goroutines | tokio (Rust) |
| **CLI Framework** | `cobra` | `clap` (Rust) |
| **Proxy Support** | `net.Dialer` + SOCKS | `hyper-socks5` |
| **Configuration** | `viper` | `config-rs` (Rust) |

**Perché Go?**
- Single binary (facile distribuzione)
- Concorrenza nativa (goroutines)
- Standard library eccellente per networking
- Cross-platform compilation triviale
- Memoria contenuta
- Dipendenze minime

### 6.2 GUI (Frontend)

| Aspetto | Consigliato | Alternativa |
|---------|------------|-------------|
| **Framework** | Tauri + React | Electron, Qt/PySide |
| **UI Library** | shadcn/ui, Material-UI | Chakra UI |
| **State Management** | Zustand, Jotai | Redux, Recoil |
| **HTTP Client** | `@tauri-apps/api` | axios |
| **Charting** | Recharts | Chart.js, Plotly |
| **Icons** | Heroicons, Lucide | Bootstrap Icons |
| **CSS** | Tailwind CSS | styled-components |
| **Build Tool** | Vite | Webpack |

**Perché Tauri?**
- Binario nativo ~150MB (vs 300+ Electron)
- Rust backend sicuro
- Consumo memoria inferiore
- API moderne (IPC)
- Distribuibile senza runtime esterno

### 6.3 Browser Automation

| Use Case | Consigliato |
|----------|-------------|
| WebRTC leak detection | Playwright (Go, TS) |
| Cookie analysis | Playwright, headless Chromium |
| JavaScript fingerprint | Playwright |
| Browser context override | Playwright |

### 6.4 External Services

| Servizio | Uso |
|----------|-----|
| MaxMind GeoIP2 | Geolocalizzazione IP |
| AbuseIPDB API | IP reputation |
| OCSP Stapling | Revocation check |
| Whois API | ASN lookup |
| Public IP Services | `api.ipify.org`, `icanhazip.com` |

---

## 7. Albero delle Cartelle

```
ProxyDoctor/
├── README.md (progetto overview)
├── ARCHITECTURE.md (questo file)
├── CONTRIBUTING.md
├── LICENSE (GPL-3.0)
├── .gitignore
├── go.mod / package.json / Cargo.toml (dipendenze)
│
├── core/
│   ├── go.mod
│   ├── check/ (interface definition)
│   │   ├── check.go
│   │   └── result.go
│   │
│   ├── engine/
│   │   ├── orchestrator.go
│   │   ├── execution_context.go
│   │   ├── registry.go
│   │   ├── dependency.go
│   │   └── reporter.go
│   │
│   ├── adapters/
│   │   ├── network_adapter.go (interface)
│   │   ├── direct.go
│   │   ├── http_proxy.go
│   │   ├── socks.go
│   │   ├── vpn.go
│   │   └── browser.go
│   │
│   ├── checks/
│   │   ├── public_ip/ {check.go, resolver.go}
│   │   ├── geolocation/ {check.go, maxmind.go}
│   │   ├── asn/ {check.go, whois.go}
│   │   ├── dns_leak/ {check.go, detector.go}
│   │   ├── webrtc_leak/ {check.go, monitor.go}
│   │   ├── ipv6_leak/ {check.go}
│   │   ├── http2_support/ {check.go}
│   │   ├── http3_quic/ {check.go}
│   │   ├── tls_handshake/ {check.go, analyzer.go}
│   │   ├── certificates/ {check.go, validator.go}
│   │   ├── redirects/ {check.go, follower.go}
│   │   ├── headers/ {check.go}
│   │   ├── cookies/ {check.go}
│   │   ├── browser_fingerprint/ {check.go, ...}
│   │   ├── streaming/ {check.go, cdn_detector.go}
│   │   └── ip_reputation/ {check.go, reputation.go}
│   │
│   ├── utils/
│   │   ├── http_client.go
│   │   ├── dns_resolver.go
│   │   ├── tls_parser.go
│   │   ├── config.go
│   │   ├── logger.go
│   │   └── cache.go
│   │
│   └── test/ (unit tests per ogni modulo)
│       ├── check_test.go
│       ├── engine_test.go
│       ├── adapters_test.go
│       └── checks_test.go
│
├── cmd/
│   ├── cli/
│   │   ├── go.mod
│   │   ├── main.go
│   │   ├── commands/
│   │   │   ├── diagnose.go
│   │   │   ├── list_checks.go
│   │   │   ├── plugin.go
│   │   │   ├── config.go
│   │   │   └── version.go
│   │   ├── formatters/
│   │   │   ├── table.go
│   │   │   ├── json.go
│   │   │   └── markdown.go
│   │   ├── styles/
│   │   │   └── colors.go
│   │   └── config/
│   │       └── default.yaml
│   │
│   └── gui/
│       ├── src-tauri/ (Rust backend)
│       │   ├── src/
│       │   │   ├── main.rs
│       │   │   ├── commands/
│       │   │   │   ├── diagnose.rs
│       │   │   │   ├── list_checks.rs
│       │   │   │   └── plugin.rs
│       │   │   └── utils/
│       │   ├── tauri.conf.json
│       │   └── Cargo.toml
│       │
│       └── src/ (React frontend)
│           ├── components/
│           │   ├── DiagnosisForm.jsx
│           │   ├── ResultsView.jsx
│           │   ├── ComparisonView.jsx
│           │   └── Settings.jsx
│           ├── hooks/
│           │   ├── useDiagnosis.js
│           │   └── usePlugins.js
│           ├── styles/
│           │   ├── global.css
│           │   └── theme.css
│           ├── utils/
│           │   └── api.js
│           ├── App.jsx
│           ├── index.html
│           └── package.json
│
├── plugins/
│   ├── README.md (plugin development guide)
│   ├── sdk/
│   │   └── go/ (plugin template Go)
│   │       └── check_sdk.go
│   │
│   ├── community/
│   │   ├── tor-exit-check/
│   │   │   ├── plugin.json
│   │   │   ├── check.go
│   │   │   └── test.go
│   │   │
│   │   └── cloudflare-check/
│   │       ├── plugin.json
│   │       ├── check.go
│   │       └── test.go
│   │
│   └── examples/
│       └── custom_port_check/
│           ├── plugin.json
│           ├── check.go
│           └── README.md
│
├── tests/
│   ├── integration/
│   │   ├── full_diagnosis_test.go
│   │   ├── proxy_detection_test.go
│   │   └── plugin_loading_test.go
│   │
│   ├── fixtures/
│   │   ├── test_urls.txt
│   │   ├── mock_responses.json
│   │   └── test_certificates/
│   │
│   └── e2e/ (end-to-end tests)
│       ├── cli_test.go
│       └── gui_test.go
│
├── docs/
│   ├── INSTALLATION.md
│   ├── QUICK_START.md
│   ├── USER_GUIDE.md
│   ├── CLI_REFERENCE.md
│   ├── PLUGIN_DEVELOPMENT.md
│   ├── API_REFERENCE.md
│   └── FAQ.md
│
├── scripts/
│   ├── build.sh (compile for all platforms)
│   ├── test.sh (run all tests)
│   ├── lint.sh (code quality)
│   └── release.sh (create release)
│
├── config/
│   ├── default.yaml (default settings)
│   ├── checks_registry.yaml (built-in checks list)
│   └── external_services.yaml (API endpoints)
│
├── .github/
│   ├── workflows/
│   │   ├── test.yml
│   │   ├── build.yml
│   │   ├── lint.yml
│   │   └── release.yml
│   │
│   └── ISSUE_TEMPLATE/
│       └── bug_report.md
│
└── ROADMAP.md
```

---

## 8. Roadmap: v0.1 → v1.0

### **v0.1 (MVP - Core Foundation)**
**Timeline:** 4-6 settimane | **Focus:** Core + CLI

- [ ] Struttura base progetto Go
- [ ] Interfaccia `CheckResult` standardizzata
- [ ] `DiagnosisEngine` e orchestration
- [ ] Network adapters: Direct + HTTP proxy + SOCKS5
- [ ] Checks built-in essenziali:
  - [ ] Public IP
  - [ ] Geolocation
  - [ ] DNS leak
  - [ ] WebRTC leak
  - [ ] TLS handshake basic
- [ ] CLI básico con `cobra`
- [ ] Plugin system skeleton
- [ ] Unit tests (80% coverage)
- [ ] Documentation base

**Rilasci:** `v0.1.0-alpha`

---

### **v0.2 (Extended Checks)**
**Timeline:** 3-4 settimane | **Focus:** More checks

- [ ] Aggiungi checks:
  - [ ] IPv6 leak
  - [ ] HTTP/2 support
  - [ ] Certificate validation
  - [ ] Redirect analysis
  - [ ] HTTP headers
  - [ ] Streaming CDN detection
- [ ] Miglior error handling
- [ ] Advanced caching system
- [ ] Performance optimization
- [ ] Integration tests

**Rilasci:** `v0.2.0`

---

### **v0.3 (Plugin Ecosystem)**
**Timeline:** 3-4 settimane | **Focus:** Plugin maturity

- [ ] Plugin SDK completo (Go)
- [ ] Plugin manager (install/uninstall/update)
- [ ] Plugin marketplace skeleton
- [ ] Example community plugins:
  - [ ] Tor exit node detector
  - [ ] Custom port scanner
  - [ ] BGP hijack detector
- [ ] Plugin security (sandboxing basics)
- [ ] Plugin versioning

**Rilasci:** `v0.3.0`

---

### **v0.4 (GUI Beta)**
**Timeline:** 5-6 settimane | **Focus:** Tauri GUI

- [ ] Setup Tauri + React
- [ ] UI base (dark mode, responsive)
- [ ] Diagnosis form & results view
- [ ] Real-time progress
- [ ] Comparison view (direct vs proxy)
- [ ] Export (JSON, HTML)
- [ ] Settings panel
- [ ] Plugin management UI
- [ ] Result history

**Rilasci:** `v0.4.0-beta`

---

### **v0.5 (Completeness)**
**Timeline:** 3-4 settimane | **Focus:** Missing checks

- [ ] Aggiungi remaining checks:
  - [ ] HTTP/3 QUIC detection
  - [ ] Browser fingerprint (advanced)
  - [ ] IP reputation (AbuseIPDB)
  - [ ] Cookie analysis advanced
- [ ] Report generation (PDF via CLI)
- [ ] Scheduling (periodic diagnostics)
- [ ] Batch diagnosis (multiple URLs)
- [ ] Advanced filtering & sorting results

**Rilasci:** `v0.5.0`

---

### **v1.0 (Production Ready)**
**Timeline:** 4-5 settimane | **Focus:** Polish & hardening

- [ ] Code audit & security review
- [ ] Performance profiling & optimization
- [ ] Comprehensive documentation
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Binary releases (Windows, macOS, Linux)
- [ ] Installer (chocolatey, brew, apt)
- [ ] Community feedback integration
- [ ] Telemetry & analytics (optional, privacy-first)
- [ ] Cloud API (optional, separate service)

**Rilasci:** `v1.0.0`

---

### **Post-v1.0 (Future)**
- [ ] Web Dashboard (separate service)
- [ ] Mobile app (API + native frontend)
- [ ] VS Code extension
- [ ] IDE integrations
- [ ] Advanced reporting & compliance
- [ ] Machine learning anomaly detection
- [ ] Real-time monitoring mode

---

## 9. Cosa Evitare (Pitfalls)

### ❌ **Non Fare:**

1. **Mixed Responsibilities**
   - Non mescolare network adapter + checks logic
   - Ogni modulo ha UN'unica responsabilità
   - ✅ **Soluzione:** Separation of concerns rigorosa

2. **Hard-coded Checks**
   - Non aggiungere nuovi check direttamente in `engine/`
   - ✅ **Soluzione:** Plugin system fin da v0.1

3. **GUI Tightly Coupled to Core**
   - Non importare direttamente il core nella GUI
   - ✅ **Soluzione:** IPC layer (Tauri commands) tra GUI e core

4. **No Standardization**
   - Non lasciare check ritornare formati diversi
   - ✅ **Soluzione:** `CheckResult` interface obbligatorio

5. **Synchronous Execution**
   - Non eseguire check uno per uno
   - ✅ **Soluzione:** DAG-based parallel execution

6. **No External Service Abstraction**
   - Non hard-coded API calls sparse nel codice
   - ✅ **Soluzione:** Centralized `ExternalServiceClient`

7. **Verbose Output by Default**
   - Non saturare l'output con debug logs
   - ✅ **Soluzione:** Configurable log levels (info, debug, trace)

8. **No Configuration Management**
   - Non hard-coded settings
   - ✅ **Soluzione:** viper/config-rs + env overrides

9. **Circular Dependencies**
   - Non creare dipendenze circolari tra moduli
   - ✅ **Soluzione:** Dependency injection pattern

10. **No Testing Strategy**
    - Non code senza unit + integration tests
    - ✅ **Soluzione:** TDD approach, 80%+ coverage target

11. **Over-engineering**
    - Non aggiungere features non necessarie in v0.1
    - ✅ **Soluzione:** MVP-first, iterate on feedback

12. **Documentation Debt**
    - Non rinviare docs, scrivi man mano
    - ✅ **Soluzione:** Docs as code, live aggiornati

---

## 10. Recommandazioni Finali

### **Suggerimenti Implementativi**

1. **Inizia con il Core**
   - Implementa prima la base (check interface, engine, adapters)
   - CLI viene naturale dopo
   - GUI è decorativa nel primo stadio

2. **Valida il Concetto**
   - Crea subito un check "proof of concept"
   - Assicurati che il flusso `Check → Result → Report` funzioni
   - Poi scala

3. **Plugin System Early**
   - Non rimandare a v0.5
   - Anche un plugin system semplice (v0.1) rende il progetto "future-proof"

4. **Testing from Day One**
   - Ogni check ha un test
   - Mock di external services (MaxMind, etc.)
   - Fixtures per test URLs

5. **Documentation Living**
   - Scrivi ARCHITECTURE.md mentre buildi
   - Aggiorna man mano
   - Non è un'attività finale

6. **Binary Distribution**
   - Priorità: Go per il core (single binary)
   - Tauri per GUI (nativa, leggera)
   - Evita "app store hell" all'inizio

7. **Community First**
   - GitHub discussions per feedback
   - Make plugin contribution easy (buona docs + examples)
   - Open source da giorno 1

---

## 11. Conclusione

**ProxyDoctor** ha il potenziale per diventare lo strumento definitivo di diagnostica proxy se costruito attorno a:
- ✅ **Modularità rigorosa** (ogni check è un plugin)
- ✅ **Standardizzazione** (`CheckResult` interface)
- ✅ **Estendibilità** (plugin system maturo)
- ✅ **Semplicità** (core lean, no bloat)

Il fondamento solido di questa architettura permette di crescere da MVP a strumento production-grade **senza refactor epici**.

**Next Step:** Implementare il core in Go, partendo dalla `CheckResult` interface e dal primo check (Public IP detection).