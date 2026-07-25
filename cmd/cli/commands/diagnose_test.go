package commands

import (
	"sort"
	"testing"
	"time"

	dnsresolve "github.com/francomano/proxydoctor/core/checks/dns_resolve"
	portscan "github.com/francomano/proxydoctor/core/checks/port_scan"
	publicip "github.com/francomano/proxydoctor/core/checks/public_ip"
	tlscert "github.com/francomano/proxydoctor/core/checks/tls_cert"
	"github.com/francomano/proxydoctor/core/engine"
)

func TestParseDiagnosisTimeout(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:  "empty uses default",
			value: "",
			want:  engine.DefaultDiagnosisTimeout,
		},
		{
			name:  "seconds",
			value: "10s",
			want:  10 * time.Second,
		},
		{
			name:  "minutes",
			value: "5m",
			want:  5 * time.Minute,
		},
		{
			name:    "below minimum",
			value:   "500ms",
			wantErr: true,
		},
		{
			name:    "above maximum",
			value:   "6m",
			wantErr: true,
		},
		{
			name:    "invalid duration",
			value:   "slow",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDiagnosisTimeout(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCheckFiltersByID(t *testing.T) {
	got, err := parseCheckFilters(" public_ip, dns_resolve ", newTestRegistry())
	if err != nil {
		t.Fatalf("parseCheckFilters returned error: %v", err)
	}

	want := []string{"dns_resolve", "public_ip"}
	assertStringSlicesEqual(t, got, want)
}

func TestParseCheckFiltersByCategory(t *testing.T) {
	got, err := parseCheckFilters("network", newTestRegistry())
	if err != nil {
		t.Fatalf("parseCheckFilters returned error: %v", err)
	}

	want := []string{"dns_resolve", "port_connectivity", "public_ip"}
	assertStringSlicesEqual(t, got, want)
}

func TestParseCheckFiltersAll(t *testing.T) {
	got, err := parseCheckFilters("all", newTestRegistry())
	if err != nil {
		t.Fatalf("parseCheckFilters returned error: %v", err)
	}
	if got != nil {
		t.Fatalf("got %v, want nil for all checks", got)
	}
}

func TestParseCheckFiltersRejectsUnknownFilter(t *testing.T) {
	if _, err := parseCheckFilters("missing", newTestRegistry()); err == nil {
		t.Fatal("expected an error for an unknown check filter")
	}
}

func newTestRegistry() *engine.CheckRegistry {
	registry := engine.NewCheckRegistry()
	registry.Register(publicip.NewPublicIPCheck())
	registry.Register(dnsresolve.NewDNSResolveCheck())
	registry.Register(tlscert.NewTLSCertCheck())
	registry.Register(portscan.NewPortScanCheck())
	return registry
}

func assertStringSlicesEqual(t *testing.T, got []string, want []string) {
	t.Helper()
	sort.Strings(got)
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
