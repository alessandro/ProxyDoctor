package commands

import (
	"fmt"
	neturl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/francomano/proxydoctor/core/adapters"
	"github.com/francomano/proxydoctor/core/check"
	publicip "github.com/francomano/proxydoctor/core/checks/public_ip"
	"github.com/francomano/proxydoctor/core/engine"
)

var (
	url       string
	proxyStr  string
	proxyType string
	exportFmt string
	output    string
	compare   bool
)

// RootCmd is the main command
var RootCmd = &cobra.Command{
	Use:   "proxyctl",
	Short: "ProxyDoctor - Comprehensive proxy diagnostics tool",
	Long: `ProxyDoctor is a command-line tool for comprehensive proxy diagnostics.
It analyzes connectivity through proxies and identifies issues.`,
}

// diagnoseCmd is the main diagnose command
var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Run a comprehensive diagnosis on a URL",
	Long: `Run a comprehensive diagnosis on the given URL through your proxy or direct connection.
The diagnosis includes multiple checks for network connectivity, security, and leaks.`,
	RunE: runDiagnose,
}

func init() {
	// Add diagnose command to root
	RootCmd.AddCommand(diagnoseCmd)
	RootCmd.AddCommand(listChecksCmd)
	RootCmd.AddCommand(versionCmd)

	// Diagnose flags
	diagnoseCmd.Flags().StringVarP(&url, "url", "u", "", "URL to diagnose (required)")
	diagnoseCmd.Flags().StringVarP(&proxyStr, "proxy", "p", "", "Proxy URL (e.g., http://localhost:8080, socks5://localhost:1080)")
	diagnoseCmd.Flags().StringVar(&proxyType, "proxy-type", "auto", "Proxy type: auto, http, https, socks4, socks5")
	diagnoseCmd.Flags().StringVarP(&exportFmt, "export", "e", "text", "Export format: text, json, html, markdown")
	diagnoseCmd.Flags().StringVarP(&output, "output", "o", "", "Output file (empty = stdout)")
	diagnoseCmd.Flags().BoolVar(&compare, "compare", false, "Compare with direct connection")

	diagnoseCmd.MarkFlagRequired("url")
}

func runDiagnose(cmd *cobra.Command, args []string) error {
	fmt.Printf("🔍 ProxyDoctor v0.1 - Proxy Diagnostics Tool\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Parse proxy configuration
	proxyConfig, err := parseProxyConfig(proxyStr, proxyType)
	if err != nil {
		fmt.Printf("❌ Invalid proxy configuration: %v\n", err)
		return err
	}

	// Create registry and register checks
	registry := engine.NewCheckRegistry()
	registry.Register(publicip.NewPublicIPCheck())
	// TODO: Add more checks as they're implemented

	// Create adapter factory
	adapterFactory := adapters.NewAdapterFactory()

	// Create orchestrator
	orchestrator := engine.NewDiagnosisOrchestrator(registry, adapterFactory, 4)

	// Create request
	diagRequest := engine.DiagnosisRequest{
		URL:         url,
		ProxyConfig: proxyConfig,
		Timeout:     30 * time.Second,
	}

	fmt.Printf("📋 Running diagnosis for: %s\n", url)
	if proxyConfig.Type != check.ProxyTypeDirect {
		fmt.Printf("🔗 Via proxy: %s://%s:%d\n", proxyConfig.Type, proxyConfig.Host, proxyConfig.Port)
	}
	fmt.Printf("\n")

	// Execute diagnosis
	report, err := orchestrator.Execute(diagRequest)
	if err != nil {
		fmt.Printf("❌ Diagnosis failed: %v\n", err)
		return err
	}

	// Format and display results
	formatResults(report, exportFmt)

	return nil
}

func parseProxyConfig(proxyStr, proxyTypeStr string) (check.ProxyConfig, error) {
	if proxyStr == "" {
		return check.ProxyConfig{Type: check.ProxyTypeDirect}, nil
	}

	parsed, err := neturl.Parse(proxyStr)
	if err != nil {
		return check.ProxyConfig{}, fmt.Errorf("invalid proxy URL %q: %w", proxyStr, err)
	}
	if parsed.Host == "" {
		return check.ProxyConfig{}, fmt.Errorf("invalid proxy URL %q: missing host", proxyStr)
	}

	// Determine proxy type: explicit flag wins, otherwise infer from the URL scheme
	var proxyType check.ProxyType
	switch strings.ToLower(proxyTypeStr) {
	case "http":
		proxyType = check.ProxyTypeHTTP
	case "https":
		proxyType = check.ProxyTypeHTTPS
	case "socks4":
		proxyType = check.ProxyTypeSOCKS4
	case "socks5":
		proxyType = check.ProxyTypeSOCKS5
	case "auto", "":
		switch strings.ToLower(parsed.Scheme) {
		case "http":
			proxyType = check.ProxyTypeHTTP
		case "https":
			proxyType = check.ProxyTypeHTTPS
		case "socks4":
			proxyType = check.ProxyTypeSOCKS4
		case "socks5", "socks":
			proxyType = check.ProxyTypeSOCKS5
		default:
			return check.ProxyConfig{}, fmt.Errorf("cannot infer proxy type from scheme %q, pass --proxy-type explicitly", parsed.Scheme)
		}
	default:
		return check.ProxyConfig{}, fmt.Errorf("unknown --proxy-type %q", proxyTypeStr)
	}

	// Resolve host/port, applying sensible defaults per proxy type
	host := parsed.Hostname()
	portStr := parsed.Port()
	var port int
	if portStr != "" {
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return check.ProxyConfig{}, fmt.Errorf("invalid proxy port %q: %w", portStr, err)
		}
		port = p
	} else {
		switch proxyType {
		case check.ProxyTypeHTTP:
			port = 8080
		case check.ProxyTypeHTTPS:
			port = 443
		case check.ProxyTypeSOCKS4, check.ProxyTypeSOCKS5:
			port = 1080
		}
	}

	config := check.ProxyConfig{
		Type: proxyType,
		Host: host,
		Port: port,
	}
	if parsed.User != nil {
		config.Username = parsed.User.Username()
		if pw, ok := parsed.User.Password(); ok {
			config.Password = pw
		}
	}

	return config, nil
}


func formatResults(report *engine.DiagnosisReport, format string) {
	switch format {
	case "json":
		formatJSON(report)
	case "text":
		formatText(report)
	case "markdown":
		formatMarkdown(report)
	default:
		formatText(report)
	}
}

func formatText(report *engine.DiagnosisReport) {
	fmt.Printf("📊 Diagnosis Results\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	for i, result := range report.Results {
		status := "✅"
		if result.IsFailed() {
			status = "❌"
		} else if result.IsError() {
			status = "⚠️"
		}

		fmt.Printf("%d. %s %s\n", i+1, status, result.ID)
		fmt.Printf("   Status: %s | Severity: %s | Confidence: %.0f%%\n",
			result.Status, result.Severity, result.Confidence*100)
		fmt.Printf("   %s\n\n", result.Explanation)
	}

	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Checks Executed: %d | Failed: %d | Critical: %d\n",
		report.ChecksExecuted, report.ChecksFailed, report.CriticalFindings)
	fmt.Printf("Total Time: %s\n", report.ExecutionTime)
}

func formatJSON(report *engine.DiagnosisReport) {
	// TODO: Implement JSON formatting
	fmt.Println("JSON export not yet implemented")
}

func formatMarkdown(report *engine.DiagnosisReport) {
	// TODO: Implement Markdown formatting
	fmt.Println("Markdown export not yet implemented")
}