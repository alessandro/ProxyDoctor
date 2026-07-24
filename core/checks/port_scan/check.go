package portscan

import (
	"fmt"
	"net/url"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

type PortScanCheck struct{}

func NewPortScanCheck() check.Checker {
	return &PortScanCheck{}
}

func (c *PortScanCheck) ID() string         { return "port_connectivity" }
func (c *PortScanCheck) Name() string        { return "Port Connectivity" }
func (c *PortScanCheck) Description() string { return "Tests TCP connectivity to common ports (80, 443, 8080, 8443) on the target host" }
func (c *PortScanCheck) Category() check.CheckCategory { return check.CategoryNetwork }
func (c *PortScanCheck) DependsOn() []string           { return []string{} }

func (c *PortScanCheck) Execute(ctx check.ExecutionContext) check.CheckResult {
	result := check.NewCheckResult(c.ID(), c.Category())
	startTime := time.Now()

	targetURL := ctx.GetURL()
	parsed, err := url.Parse(targetURL)
	if err != nil {
		result.SetExecutionTime(time.Since(startTime))
		return *result.WithStatus(check.StatusError, check.SeverityCritical).
			WithExplanation(fmt.Sprintf("Invalid URL: %v", err)).
			WithConfidence(0)
	}

	hostname := parsed.Hostname()

	adapter := ctx.GetProxyAdapter()
	if ctx.GetProxyConfig().Type == check.ProxyTypeDirect {
		adapter = ctx.GetDirectAdapter()
	}

	ports := []int{80, 443, 8080, 8443}
	timeout := 5 * time.Second
	openPorts := []int{}
	closedPorts := []int{}

	for _, port := range ports {
		open, _ := adapter.TestPort(hostname, port, timeout)
		if open {
			openPorts = append(openPorts, port)
		} else {
			closedPorts = append(closedPorts, port)
		}
	}

	ctx.SetSharedData("open_ports", openPorts)

	if len(openPorts) == 0 {
		result.SetExecutionTime(time.Since(startTime))
		return *result.WithStatus(check.StatusFailed, check.SeverityCritical).
			WithExplanation(fmt.Sprintf("No open ports detected on %s (tested: %v)", hostname, ports)).
			WithConfidence(0.9).
			AddEvidence("hostname", hostname).
			AddEvidence("open_ports", openPorts).
			AddEvidence("closed_ports", closedPorts).
			AddProbableCause("Host may be unreachable through the proxy").
			AddProbableCause("Firewall may be blocking connections").
			AddSuggestedAction("Verify the proxy is working and the target host is online")
	}

	result.SetExecutionTime(time.Since(startTime))
	return *result.WithStatus(check.StatusPassed, check.SeverityInfo).
		WithExplanation(fmt.Sprintf("%d open port(s) on %s: %v", len(openPorts), hostname, openPorts)).
		WithConfidence(0.9).
		AddEvidence("hostname", hostname).
		AddEvidence("open_ports", openPorts).
		AddEvidence("closed_ports", closedPorts)
}
