package switchbot

import (
	"strings"
	"testing"
)

func TestDecodeMetrics(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
		want     Metrics
	}{
		{
			name:     "successful meter response",
			response: `{"statusCode":100,"message":"success","body":{"temperature":23.5,"humidity":48}}`,
			want:     Metrics{Temperature: 23.5, Humidity: 48},
		},
		{
			name:     "API error without body",
			response: `{"statusCode":190,"message":"System error"}`,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeMetrics(strings.NewReader(tt.response))
			if (err != nil) != tt.wantErr {
				t.Fatalf("decodeMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if got != nil {
					t.Fatalf("decodeMetrics() returned metrics on API error: %+v", got)
				}
				return
			}
			if *got != tt.want {
				t.Fatalf("decodeMetrics() = %+v, want %+v", *got, tt.want)
			}
		})
	}
}
