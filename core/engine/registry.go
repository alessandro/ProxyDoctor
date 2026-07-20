package engine

import (
	"fmt"
	"sync"

	"github.com/francomano/proxydoctor/core/check"
)

// CheckRegistry manages all available checks
type CheckRegistry struct {
	checks map[string]check.Checker
	mu     sync.RWMutex
}

// NewCheckRegistry creates a new check registry
func NewCheckRegistry() *CheckRegistry {
	return &CheckRegistry{
		checks: make(map[string]check.Checker),
	}
}

// Register registers a new check
func (r *CheckRegistry) Register(checker check.Checker) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.checks[checker.ID()]; exists {
		return fmt.Errorf("check with ID %s already registered", checker.ID())
	}

	r.checks[checker.ID()] = checker
	return nil
}

// Unregister removes a check from the registry
func (r *CheckRegistry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.checks[id]; !exists {
		return fmt.Errorf("check with ID %s not found", id)
	}

	delete(r.checks, id)
	return nil
}

// GetCheck returns a check by ID
func (r *CheckRegistry) GetCheck(id string) (check.Checker, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	checker, ok := r.checks[id]
	return checker, ok
}

// ListChecks returns all registered checks
func (r *CheckRegistry) ListChecks() map[string]check.Checker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]check.Checker)
	for id, checker := range r.checks {
		result[id] = checker
	}
	return result
}

// GetCheckMetadata returns metadata for a check
func (r *CheckRegistry) GetCheckMetadata(id string) *CheckMetadata {
	if checker, ok := r.GetCheck(id); ok {
		return &CheckMetadata{
			ID:          checker.ID(),
			Name:        checker.Name(),
			Description: checker.Description(),
			Category:    string(checker.Category()),
			DependsOn:   checker.DependsOn(),
		}
	}
	return nil
}

// CheckMetadata contains metadata about a check
type CheckMetadata struct {
	ID          string
	Name        string
	Description string
	Category    string
	DependsOn   []string
}
