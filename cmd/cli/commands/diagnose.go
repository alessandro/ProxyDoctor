package commands

import (
	"fmt"
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

	// Simple parser - would be more robust in production
	config := check.ProxyConfig{
		Type: check.ProxyTypeHTTP,
		Host: "localhost",
		Port: 8080,
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
