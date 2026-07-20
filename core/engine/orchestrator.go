package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/francomano/proxydoctor/core/check"
)

// DiagnosisRequest represents the input for a diagnosis
type DiagnosisRequest struct {
	URL        string
	ProxyConfig check.ProxyConfig
	CheckIDs   []string // empty = all checks
	Timeout    time.Duration
}

// DiagnosisReport represents the final output report
type DiagnosisReport struct {
	ID                string
	RequestMetadata   RequestMetadata
	Results           []check.CheckResult
	ExecutionTime     time.Duration
	ChecksExecuted    int
	ChecksFailed      int
	CriticalFindings  int
	WarningFindings   int
}

// RequestMetadata contains information about the diagnosis request
type RequestMetadata struct {
	URL              string
	ProxyType        check.ProxyType
	StartedAt        time.Time
	CompletedAt      time.Time
	UserAgent        string
}

// DiagnosisOrchestrator orchestrates the execution of checks
type DiagnosisOrchestrator struct {
	registry  *CheckRegistry
	adapters  AdapterFactory
	maxWorkers int
}

// NewDiagnosisOrchestrator creates a new orchestrator
func NewDiagnosisOrchestrator(registry *CheckRegistry, adapters AdapterFactory, maxWorkers int) *DiagnosisOrchestrator {
	if maxWorkers <= 0 {
		maxWorkers = 4
	}
	return &DiagnosisOrchestrator{
		registry:  registry,
		adapters:  adapters,
		maxWorkers: maxWorkers,
	}
}

// Execute runs the diagnosis with the given request
func (o *DiagnosisOrchestrator) Execute(req DiagnosisRequest) (*DiagnosisReport, error) {
	startTime := time.Now()

	// Validate input
	if err := o.validateRequest(req); err != nil {
		return nil, err
	}

	// Get checks to execute
	checksToRun := o.getChecksToRun(req.CheckIDs)
	if len(checksToRun) == 0 {
		return nil, fmt.Errorf("no checks to execute")
	}

	// Create network adapters
	directAdapter := o.adapters.CreateAdapter(check.ProxyTypeDirect, check.ProxyConfig{})
	proxyAdapter := o.adapters.CreateAdapter(req.ProxyConfig.Type, req.ProxyConfig)

	// Create execution context
	ctx := NewExecutionContext(req.URL, req.ProxyConfig, directAdapter, proxyAdapter, req.Timeout)

	// Build dependency graph
	dag := o.buildDependencyGraph(checksToRun)

	// Execute checks
	results := o.executeChecksParallel(ctx, dag, checksToRun)

	// Generate report
	report := o.generateReport(req, results, startTime)

	return report, nil
}

func (o *DiagnosisOrchestrator) validateRequest(req DiagnosisRequest) error {
	if req.URL == "" {
		return fmt.Errorf("URL is required")
	}
	if req.Timeout == 0 {
		req.Timeout = 30 * time.Second
	}
	return nil
}

func (o *DiagnosisOrchestrator) getChecksToRun(checkIDs []string) map[string]check.Checker {
	if len(checkIDs) == 0 {
		return o.registry.ListChecks()
	}

	result := make(map[string]check.Checker)
	for _, id := range checkIDs {
		if c, ok := o.registry.GetCheck(id); ok {
			result[id] = c
		}
	}
	return result
}

func (o *DiagnosisOrchestrator) buildDependencyGraph(checks map[string]check.Checker) DependencyDAG {
	dag := NewDependencyDAG()
	for id, checker := range checks {
		dag.AddNode(id, checker)
		for _, depID := range checker.DependsOn() {
			dag.AddEdge(depID, id)
		}
	}
	return dag
}

func (o *DiagnosisOrchestrator) executeChecksParallel(
	ctx check.ExecutionContext,
	dag DependencyDAG,
	checks map[string]check.Checker,
) []check.CheckResult {
	results := make([]check.CheckResult, 0)
	resultsChan := make(chan check.CheckResult, len(checks))
	var wg sync.WaitGroup

	// Get execution order from DAG
	executionOrder := dag.TopologicalSort()

	// Semaphore for limiting concurrent executions
	semaphore := make(chan struct{}, o.maxWorkers)

	for _, id := range executionOrder {
		if checker, ok := checks[id]; ok {
			wg.Add(1)
			go func(checkID string, c check.Checker) {
				defer wg.Done()
				semaphore <- struct{}{}        // Acquire
				defer func() { <-semaphore }() // Release

				if execCtx, ok := ctx.(*DefaultExecutionContext); ok {
					if execCtx.IsCancelled() {
						return
					}
				}

				startTime := time.Now()
				result := c.Execute(ctx)
				result.ExecutionTime = time.Since(startTime)
				result.Timestamp = time.Now()

				resultsChan <- result
			}(id, checker)
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func (o *DiagnosisOrchestrator) generateReport(
	req DiagnosisRequest,
	results []check.CheckResult,
	startTime time.Time,
) *DiagnosisReport {
	report := &DiagnosisReport{
		ID: fmt.Sprintf("diag_%d", time.Now().Unix()),
		RequestMetadata: RequestMetadata{
			URL:       req.URL,
			ProxyType: req.ProxyConfig.Type,
			StartedAt: startTime,
			CompletedAt: time.Now(),
		},
		Results:       results,
		ExecutionTime: time.Since(startTime),
		ChecksExecuted: len(results),
	}

	// Count failures and findings
	for _, r := range results {
		if r.IsFailed() {
			report.ChecksFailed++
		}
		if r.Severity == check.SeverityCritical {
			report.CriticalFindings++
		} else if r.Severity == check.SeverityWarning {
			report.WarningFindings++
		}
	}

	return report
}
