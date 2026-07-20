package check

import (
	"testing"
	"time"
)

func TestCheckResultBuilder(t *testing.T) {
	result := NewCheckResult("test_id", CategoryNetwork)

	result.
		WithStatus(StatusPassed, SeverityInfo).
		WithConfidence(0.95).
		WithExplanation("Test passed successfully").
		AddEvidence("key1", "value1").
		AddProbableCause("cause1").
		AddSuggestedAction("action1").
		AddReference("https://example.com").
		SetExecutionTime(100 * time.Millisecond)

	if result.ID != "test_id" {
		t.Errorf("ID mismatch: got %s, want test_id", result.ID)
	}

	if result.Status != StatusPassed {
		t.Errorf("Status mismatch: got %v, want %v", result.Status, StatusPassed)
	}

	if result.Confidence != 0.95 {
		t.Errorf("Confidence mismatch: got %f, want 0.95", result.Confidence)
	}

	if result.IsPassed() == false {
		t.Error("IsPassed() should return true")
	}

	if len(result.Evidence) != 1 {
		t.Errorf("Evidence count mismatch: got %d, want 1", len(result.Evidence))
	}
}

func TestCheckResultConfidenceValidation(t *testing.T) {
	result := NewCheckResult("test", CategoryNetwork)

	// Test confidence > 1
	result.WithConfidence(1.5)
	if result.Confidence != 1.0 {
		t.Errorf("Confidence should be capped at 1.0, got %f", result.Confidence)
	}

	// Test confidence < 0
	result.WithConfidence(-0.5)
	if result.Confidence != 0.0 {
		t.Errorf("Confidence should be floored at 0.0, got %f", result.Confidence)
	}
}
