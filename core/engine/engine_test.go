package engine

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/francomano/proxydoctor/core/adapters"
	"github.com/francomano/proxydoctor/core/check"
	publicip "github.com/francomano/proxydoctor/core/checks/public_ip"
)

func TestDiagnosisOrchestrator(t *testing.T) {
	registry := NewCheckRegistry()
	registry.Register(publicip.NewPublicIPCheck())
	adapterFactory := adapters.NewAdapterFactory()
	orchestrator := NewDiagnosisOrchestrator(registry, adapterFactory, 2)

	request := DiagnosisRequest{
		URL: "https://example.com",
		ProxyConfig: check.ProxyConfig{
			Type: check.ProxyTypeDirect,
		},
		Timeout: 10 * time.Second,
	}

	report, err := orchestrator.Execute(request)

	if err != nil {
		t.Logf("Diagnosis execution: %v (this is expected for unit test)", err)
	}

	if report != nil {
		if report.ChecksExecuted == 0 {
			t.Error("No checks were executed")
		}
	}
}

func TestCheckRegistry(t *testing.T) {
	registry := NewCheckRegistry()

	check1 := publicip.NewPublicIPCheck()
	err := registry.Register(check1)
	if err != nil {
		t.Fatalf("Failed to register check: %v", err)
	}

	retrieved, ok := registry.GetCheck(check1.ID())
	if !ok {
		t.Error("Check not found in registry")
	}

	if retrieved.ID() != check1.ID() {
		t.Errorf("Check ID mismatch: got %s, want %s", retrieved.ID(), check1.ID())
	}

	checks := registry.ListChecks()
	if len(checks) != 1 {
		t.Errorf("List returned %d checks, want 1", len(checks))
	}

	err = registry.Register(check1)
	if err == nil {
		t.Error("Should not allow duplicate registration")
	}

	err = registry.Unregister(check1.ID())
	if err != nil {
		t.Fatalf("Failed to unregister check: %v", err)
	}

	checks = registry.ListChecks()
	if len(checks) != 0 {
		t.Errorf("After unregister, list returned %d checks, want 0", len(checks))
	}
}

func TestDependencyDAG(t *testing.T) {
	dag := NewDependencyDAG()

	check1 := publicip.NewPublicIPCheck()
	check2 := publicip.NewPublicIPCheck()

	dag.AddNode("check1", check1)
	dag.AddNode("check2", check2)
	dag.AddEdge("check1", "check2")

	order := dag.TopologicalSort()
	if len(order) != 2 {
		t.Errorf("TopologicalSort returned %d nodes, want 2", len(order))
	}

	check1Idx := -1
	check2Idx := -1
	for i, id := range order {
		if id == "check1" {
			check1Idx = i
		}
		if id == "check2" {
			check2Idx = i
		}
	}

	if check1Idx == -1 || check2Idx == -1 {
		t.Error("Not all checks found in topological sort")
	}

	if check1Idx >= check2Idx {
		t.Error("check1 should come before check2 in topological order")
	}

	if dag.HasCycle() {
		t.Error("DAG should not have a cycle")
	}
}

func TestExecutionContext(t *testing.T) {
	adapterFactory := adapters.NewAdapterFactory()
	directAdapter := adapterFactory.CreateAdapter(check.ProxyTypeDirect, check.ProxyConfig{})
	proxyAdapter := adapterFactory.CreateAdapter(check.ProxyTypeHTTP, check.ProxyConfig{
		Host: "localhost",
		Port: 8080,
	})

	ctx := NewExecutionContext(
		"https://example.com",
		check.ProxyConfig{
			Type: check.ProxyTypeHTTP,
			Host: "localhost",
			Port: 8080,
		},
		directAdapter,
		proxyAdapter,
		30*time.Second,
	)

	if ctx.GetURL() != "https://example.com" {
		t.Errorf("URL mismatch: got %s, want https://example.com", ctx.GetURL())
	}

	ctx.SetSharedData("test_key", "test_value")
	value := ctx.GetSharedData("test_key")
	if value != "test_value" {
		t.Errorf("Shared data mismatch: got %v, want test_value", value)
	}

	if ctx.GetTimeout() != 30*time.Second {
		t.Errorf("Timeout mismatch: got %v, want 30s", ctx.GetTimeout())
	}

	if ctx.IsCancelled() {
		t.Error("Should not be cancelled initially")
	}
}

func TestExecutePropagatesCustomTimeout(t *testing.T) {
	registry := NewCheckRegistry()
	registry.Register(publicip.NewPublicIPCheck())
	adapterFactory := adapters.NewAdapterFactory()
	orchestrator := NewDiagnosisOrchestrator(registry, adapterFactory, 2)

	timeout := 45 * time.Second
	request := DiagnosisRequest{
		URL: "https://example.com",
		ProxyConfig: check.ProxyConfig{
			Type: check.ProxyTypeDirect,
		},
		Timeout: timeout,
	}

	report, err := orchestrator.Execute(request)
	if err != nil {
		t.Logf("Execute returned error (expected in test env): %v", err)
	}
	if report == nil {
		t.Fatal("report should not be nil when error is non-fatal")
	}
	if report.RequestMetadata.Timeout != timeout {
		t.Errorf("custom timeout not propagated: got %v, want %v", report.RequestMetadata.Timeout, timeout)
	}
}

func TestExecuteAppliesDefaultTimeout(t *testing.T) {
	registry := NewCheckRegistry()
	registry.Register(publicip.NewPublicIPCheck())
	adapterFactory := adapters.NewAdapterFactory()
	orchestrator := NewDiagnosisOrchestrator(registry, adapterFactory, 2)

	request := DiagnosisRequest{
		URL: "https://example.com",
		ProxyConfig: check.ProxyConfig{
			Type: check.ProxyTypeDirect,
		},
		Timeout: 0,
	}

	report, err := orchestrator.Execute(request)
	if err != nil {
		t.Logf("Execute returned error (expected in test env): %v", err)
	}
	if report == nil {
		t.Fatal("report should not be nil when error is non-fatal")
	}
	if report.RequestMetadata.Timeout != DefaultDiagnosisTimeout {
		t.Errorf("default timeout not applied: got %v, want %v", report.RequestMetadata.Timeout, DefaultDiagnosisTimeout)
	}
}

func TestGetChecksToRunFiltersSelectedIDs(t *testing.T) {
	registry := NewCheckRegistry()
	registry.Register(newTestCheck("public_ip"))
	registry.Register(newTestCheck("dns_resolve"))
	registry.Register(newTestCheck("tls_certificate"))

	orchestrator := NewDiagnosisOrchestrator(registry, adapters.NewAdapterFactory(), 2)

	checksToRun, err := orchestrator.getChecksToRun([]string{"dns_resolve", "public_ip"})
	if err != nil {
		t.Fatalf("getChecksToRun returned error: %v", err)
	}

	got := make([]string, 0, len(checksToRun))
	for id := range checksToRun {
		got = append(got, id)
	}
	sort.Strings(got)

	want := []string{"dns_resolve", "public_ip"}
	if len(got) != len(want) {
		t.Fatalf("selected %d checks, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("selected checks mismatch: got %v, want %v", got, want)
		}
	}
}

func TestGetChecksToRunRejectsUnknownID(t *testing.T) {
	registry := NewCheckRegistry()
	registry.Register(newTestCheck("public_ip"))

	orchestrator := NewDiagnosisOrchestrator(registry, adapters.NewAdapterFactory(), 2)

	if _, err := orchestrator.getChecksToRun([]string{"missing_check"}); err == nil {
		t.Fatal("expected an error for an unknown check ID")
	}
}

func TestExecuteComparisonRunsDirectAndProxyDiagnoses(t *testing.T) {
	registry := NewCheckRegistry()
	if err := registry.Register(proxyAwareCheck{}); err != nil {
		t.Fatalf("Failed to register check: %v", err)
	}

	orchestrator := NewDiagnosisOrchestrator(registry, fakeAdapterFactory{}, 1)
	report, err := orchestrator.ExecuteComparison(DiagnosisRequest{
		URL: "https://example.com",
		ProxyConfig: check.ProxyConfig{
			Type: check.ProxyTypeHTTP,
			Host: "proxy.local",
			Port: 8080,
		},
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatalf("ExecuteComparison returned error: %v", err)
	}

	if report.DirectReport.RequestMetadata.ProxyType != check.ProxyTypeDirect {
		t.Fatalf("direct report used proxy type %s", report.DirectReport.RequestMetadata.ProxyType)
	}
	if report.ProxyReport.RequestMetadata.ProxyType != check.ProxyTypeHTTP {
		t.Fatalf("proxy report used proxy type %s", report.ProxyReport.RequestMetadata.ProxyType)
	}
	if len(report.Differences) == 0 {
		t.Fatal("expected comparison differences")
	}

	foundIPDiff := false
	for _, diff := range report.Differences {
		if diff.CheckID == "public_ip" && diff.Field == "ip_address" {
			foundIPDiff = true
			break
		}
	}
	if !foundIPDiff {
		t.Fatalf("expected public_ip ip_address diff, got %#v", report.Differences)
	}
}

func TestExecuteComparisonRequiresProxy(t *testing.T) {
	orchestrator := NewDiagnosisOrchestrator(NewCheckRegistry(), fakeAdapterFactory{}, 1)
	_, err := orchestrator.ExecuteComparison(DiagnosisRequest{
		URL:         "https://example.com",
		ProxyConfig: check.ProxyConfig{Type: check.ProxyTypeDirect},
	})
	if err == nil {
		t.Fatal("expected error for direct-only comparison")
	}
}

type testCheck struct {
	id string
}

func newTestCheck(id string) testCheck {
	return testCheck{id: id}
}

func (c testCheck) ID() string                    { return c.id }
func (c testCheck) Name() string                  { return c.id }
func (c testCheck) Description() string           { return "test check" }
func (c testCheck) Category() check.CheckCategory { return check.CategoryNetwork }
func (c testCheck) DependsOn() []string           { return nil }
func (c testCheck) Execute(ctx check.ExecutionContext) check.CheckResult {
	return check.CheckResult{
		ID:         c.id,
		Category:   c.Category(),
		Status:     check.StatusPassed,
		Severity:   check.SeverityInfo,
		Confidence: 1,
	}
}

type proxyAwareCheck struct{}

func (proxyAwareCheck) ID() string {
	return "public_ip"
}

func (proxyAwareCheck) Name() string {
	return "Public IP Detection"
}

func (proxyAwareCheck) Description() string {
	return "detects public IP"
}

func (proxyAwareCheck) Category() check.CheckCategory {
	return check.CategoryNetwork
}

func (proxyAwareCheck) DependsOn() []string {
	return nil
}

func (proxyAwareCheck) Execute(ctx check.ExecutionContext) check.CheckResult {
	ip := "203.0.113.10"
	if ctx.GetProxyConfig().Type != check.ProxyTypeDirect {
		ip = "198.51.100.20"
	}

	result := check.NewCheckResult("public_ip", check.CategoryNetwork)
	result.WithStatus(check.StatusPassed, check.SeverityInfo).
		WithConfidence(1).
		WithExplanation(fmt.Sprintf("Public IP detected: %s", ip)).
		AddEvidence("ip_address", ip)
	return *result
}

type fakeAdapterFactory struct{}

func (fakeAdapterFactory) CreateAdapter(check.ProxyType, check.ProxyConfig) check.NetworkAdapter {
	return fakeNetworkAdapter{}
}

type fakeNetworkAdapter struct{}

func (fakeNetworkAdapter) Type() check.ProxyType {
	return check.ProxyTypeDirect
}

func (fakeNetworkAdapter) ExecuteHTTPRequest(*check.HTTPRequest) (*check.HTTPResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (fakeNetworkAdapter) FollowRedirects(string, int) ([]check.RedirectStep, error) {
	return nil, fmt.Errorf("not implemented")
}

func (fakeNetworkAdapter) ResolveDNS(string) ([]string, error) {
	return nil, fmt.Errorf("not implemented")
}

func (fakeNetworkAdapter) GetPublicIP() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (fakeNetworkAdapter) TestPort(string, int, time.Duration) (bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (fakeNetworkAdapter) GetTLSCertificate(string) (*check.CertificateInfo, error) {
	return nil, fmt.Errorf("not implemented")
}

func (fakeNetworkAdapter) GetTLSCipherSuite(string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (fakeNetworkAdapter) GetTLSVersion(string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
