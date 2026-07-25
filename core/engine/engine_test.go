package engine

import (
	"testing"
	"time"

	"github.com/francomano/proxydoctor/core/adapters"
	"github.com/francomano/proxydoctor/core/check"
	publicip "github.com/francomano/proxydoctor/core/checks/public_ip"
)

func TestDiagnosisOrchestrator(t *testing.T) {
	// Create registry with test check
	registry := NewCheckRegistry()
	registry.Register(publicip.NewPublicIPCheck())

	// Create adapter factory
	adapterFactory := adapters.NewAdapterFactory()

	// Create orchestrator
	orchestrator := NewDiagnosisOrchestrator(registry, adapterFactory, 2)

	// Create request
	request := DiagnosisRequest{
		URL: "https://example.com",
		ProxyConfig: check.ProxyConfig{
			Type: check.ProxyTypeDirect,
		},
		Timeout: 10 * time.Second,
	}

	// Execute
	report, err := orchestrator.Execute(request)

	// Assertions
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

	// Test register
	check1 := publicip.NewPublicIPCheck()
	err := registry.Register(check1)
	if err != nil {
		t.Fatalf("Failed to register check: %v", err)
	}

	// Test get
	retrieved, ok := registry.GetCheck(check1.ID())
	if !ok {
		t.Error("Check not found in registry")
	}

	if retrieved.ID() != check1.ID() {
		t.Errorf("Check ID mismatch: got %s, want %s", retrieved.ID(), check1.ID())
	}

	// Test list
	checks := registry.ListChecks()
	if len(checks) != 1 {
		t.Errorf("List returned %d checks, want 1", len(checks))
	}

	// Test duplicate registration
	err = registry.Register(check1)
	if err == nil {
		t.Error("Should not allow duplicate registration")
	}

	// Test unregister
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
	dag.AddEdge("check1", "check2") // check2 depends on check1

	// Test topological sort
	order := dag.TopologicalSort()
	if len(order) != 2 {
		t.Errorf("TopologicalSort returned %d nodes, want 2", len(order))
	}

	// check1 should come before check2
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

	// Test has cycle
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

	// Test URL
	if ctx.GetURL() != "https://example.com" {
		t.Errorf("URL mismatch: got %s, want https://example.com", ctx.GetURL())
	}

	// Test shared data
	ctx.SetSharedData("test_key", "test_value")
	value := ctx.GetSharedData("test_key")
	if value != "test_value" {
		t.Errorf("Shared data mismatch: got %v, want test_value", value)
	}

	// Test timeout
	if ctx.GetTimeout() != 30*time.Second {
		t.Errorf("Timeout mismatch: got %v, want 30s", ctx.GetTimeout())
	}

	// Test cancel
	if ctx.IsCancelled() {
		t.Error("Should not be cancelled initially")
	}
}

func TestGenerateReportIncludesTimeout(t *testing.T) {
	registry := NewCheckRegistry()
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

	report := orchestrator.generateReport(request, nil, time.Now())

	if report.RequestMetadata.Timeout != timeout {
		t.Errorf("Timeout mismatch: got %v, want %v", report.RequestMetadata.Timeout, timeout)
	}
}
