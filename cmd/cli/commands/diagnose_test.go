package commands

import (
	"testing"
	"time"
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
			want:  defaultDiagnosisTimeout,
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
