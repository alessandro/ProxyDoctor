package commands

import (
	"strings"
	"testing"
	"time"

	"github.com/francomano/proxydoctor/core/check"
	"github.com/francomano/proxydoctor/core/engine"
)

func TestFormatComparisonOutputsDifferences(t *testing.T) {
	report := &engine.ComparisonReport{
		ID: "compare_test",
		DirectReport: &engine.DiagnosisReport{
			ChecksExecuted:   1,
			ChecksFailed:     0,
			CriticalFindings: 0,
			ExecutionTime:    time.Second,
		},
		ProxyReport: &engine.DiagnosisReport{
			ChecksExecuted:   1,
			ChecksFailed:     1,
			CriticalFindings: 1,
			ExecutionTime:    2 * time.Second,
		},
		Differences: []engine.ComparisonDifference{
			{
				CheckID:     "public_ip",
				Field:       "ip_address",
				DirectValue: "203.0.113.10",
				ProxyValue:  "198.51.100.20",
				Summary:     "public_ip IP changed from 203.0.113.10 to 198.51.100.20",
			},
		},
		ExecutionTime: 3 * time.Second,
	}

	text := formatComparisonText(report)
	if !strings.Contains(text, "Direct Connection") ||
		!strings.Contains(text, "Proxied Connection") ||
		!strings.Contains(text, "public_ip IP changed") {
		t.Fatalf("text comparison output missing expected content:\n%s", text)
	}

	markdown := formatComparisonMarkdown(report)
	if !strings.Contains(markdown, "# ProxyDoctor Comparison Report") ||
		!strings.Contains(markdown, "**public_ip**") {
		t.Fatalf("markdown comparison output missing expected content:\n%s", markdown)
	}

	html := formatComparisonHTML(report)
	if !strings.Contains(html, "<h1>ProxyDoctor Comparison Report</h1>") ||
		!strings.Contains(html, "public_ip IP changed") {
		t.Fatalf("html comparison output missing expected content:\n%s", html)
	}

	json := formatComparisonJSON(report)
	if !strings.Contains(json, `"direct_report"`) ||
		!strings.Contains(json, `"differences"`) ||
		!strings.Contains(json, `"proxy_value": "198.51.100.20"`) {
		t.Fatalf("json comparison output missing expected content:\n%s", json)
	}
}

func TestFormatComparisonOutputsNoDifferences(t *testing.T) {
	report := &engine.ComparisonReport{
		DirectReport: &engine.DiagnosisReport{
			Results: []check.CheckResult{},
		},
		ProxyReport: &engine.DiagnosisReport{
			Results: []check.CheckResult{},
		},
	}

	if got := formatComparisonText(report); !strings.Contains(got, "No differences detected") {
		t.Fatalf("expected no differences message, got:\n%s", got)
	}
}
