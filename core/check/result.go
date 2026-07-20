package check

import (
	"fmt"
	"time"
)

// NewCheckResult creates a new CheckResult with common fields populated
func NewCheckResult(id string, category CheckCategory) *CheckResult {
	return &CheckResult{
		ID:        id,
		Category:  category,
		Timestamp: time.Now(),
		Evidence:  make(map[string]interface{}),
	}
}

// WithStatus sets the status and returns the result for chaining
func (r *CheckResult) WithStatus(status Status, severity Severity) *CheckResult {
	r.Status = status
	r.Severity = severity
	return r
}

// WithConfidence sets the confidence score
func (r *CheckResult) WithConfidence(confidence float64) *CheckResult {
	if confidence < 0 {
		confidence = 0
	}
	if confidence > 1 {
		confidence = 1
	}
	r.Confidence = confidence
	return r
}

// WithExplanation sets the human-readable explanation
func (r *CheckResult) WithExplanation(explanation string) *CheckResult {
	r.Explanation = explanation
	return r
}

// AddEvidence adds a piece of evidence to the result
func (r *CheckResult) AddEvidence(key string, value interface{}) *CheckResult {
	r.Evidence[key] = value
	return r
}

// AddProbableCause adds a probable cause
func (r *CheckResult) AddProbableCause(cause string) *CheckResult {
	r.ProbableCauses = append(r.ProbableCauses, cause)
	return r
}

// AddSuggestedAction adds a suggested action
func (r *CheckResult) AddSuggestedAction(action string) *CheckResult {
	r.SuggestedActions = append(r.SuggestedActions, action)
	return r
}

// AddReference adds a reference URL
func (r *CheckResult) AddReference(url string) *CheckResult {
	r.References = append(r.References, url)
	return r
}

// SetExecutionTime sets the execution time
func (r *CheckResult) SetExecutionTime(duration time.Duration) *CheckResult {
	r.ExecutionTime = duration
	return r
}

// String returns a formatted string representation
func (r *CheckResult) String() string {
	return fmt.Sprintf(
		"CheckResult{ID: %s, Status: %s, Severity: %s, Confidence: %.2f}",
		r.ID, r.Status, r.Severity, r.Confidence,
	)
}

// IsPassed returns true if the check passed
func (r *CheckResult) IsPassed() bool {
	return r.Status == StatusPassed
}

// IsFailed returns true if the check failed
func (r *CheckResult) IsFailed() bool {
	return r.Status == StatusFailed
}

// IsError returns true if the check had an error
func (r *CheckResult) IsError() bool {
	return r.Status == StatusError
}
