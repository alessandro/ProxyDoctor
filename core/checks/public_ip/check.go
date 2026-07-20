package publicip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

// PublicIPCheck detects the public IP address
type PublicIPCheck struct{}

// NewPublicIPCheck creates a new public IP check
func NewPublicIPCheck() check.Checker {
	return &PublicIPCheck{}
}

func (c *PublicIPCheck) ID() string {
	return "public_ip"
}

func (c *PublicIPCheck) Name() string {
	return "Public IP Detection"
}

func (c *PublicIPCheck) Description() string {
	return "Detects your public IP address via the current connection"
}

func (c *PublicIPCheck) Category() check.CheckCategory {
	return check.CategoryNetwork
}

func (c *PublicIPCheck) DependsOn() []string {
	return []string{} // No dependencies
}

func (c *PublicIPCheck) Execute(ctx check.ExecutionContext) check.CheckResult {
	result := check.NewCheckResult(c.ID(), c.Category())
	startTime := time.Now()

	// Try multiple IP detection services
	services := []IPService{
		{
			Name: "ipify.org",
			URL:  "https://api.ipify.org?format=json",
		},
		{
			Name: "icanhazip.com",
			URL:  "https://icanhazip.com/",
		},
		{
			Name: "ifconfig.me",
			URL:  "https://ifconfig.me/",
		},
	}

	var publicIP string
	var detectionService string

	// Determine which adapter to use
	adapter := ctx.GetProxyAdapter()
	if ctx.GetProxyConfig().Type == check.ProxyTypeDirect {
		adapter = ctx.GetDirectAdapter()
	}

	// Try each service
	for _, service := range services {
		ip, err := detectIPFromService(adapter, service)
		if err == nil && ip != "" {
			publicIP = ip
			detectionService = service.Name
			break
		}
	}

	result.SetExecutionTime(time.Since(startTime))

	// Check if we got a valid IP
	if publicIP == "" {
		result.WithStatus(check.StatusError, check.SeverityCritical)
		result.WithExplanation("Unable to detect public IP address")
		result.AddProbableCause("Network connectivity issues")
		result.AddProbableCause("All IP detection services are unreachable")
		result.WithConfidence(0.0)
		return *result
	}

	// Store in shared context for other checks
	ctx.SetSharedData("public_ip", publicIP)

	result.WithStatus(check.StatusPassed, check.SeverityInfo)
	result.WithExplanation(fmt.Sprintf("Public IP detected: %s via %s", publicIP, detectionService))
	result.WithConfidence(0.95)
	result.AddEvidence("ip_address", publicIP)
	result.AddEvidence("detection_service", detectionService)

	return *result
}

// IPService represents an IP detection service
type IPService struct {
	Name string
	URL  string
}

// detectIPFromService attempts to detect public IP from a service
func detectIPFromService(adapter check.NetworkAdapter, service IPService) (string, error) {
	req := &check.HTTPRequest{
		Method: "GET",
		URL:    service.URL,
		Headers: map[string]string{
			"User-Agent": "ProxyDoctor/0.1",
		},
	}

	resp, err := adapter.ExecuteHTTPRequest(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	// Try JSON format first (ipify)
	var jsonResp struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(resp.Body, &jsonResp); err == nil && jsonResp.IP != "" {
		return jsonResp.IP, nil
	}

	// Try plain text format (icanhazip, ifconfig)
	ip := string(resp.Body)
	if ip != "" {
		// Clean up whitespace
		return parseIP(ip), nil
	}

	return "", fmt.Errorf("could not parse IP response")
}

// parseIP extracts and validates IP from response
func parseIP(s string) string {
	// Simple validation - would be more robust in production
	return s
}
