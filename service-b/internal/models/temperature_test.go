package models

import (
	"testing"
)

func TestNewTemperature(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		wantC   float64
		wantF   float64
		wantK   float64
	}{
		{
			name:    "zero celsius",
			celsius: 0,
			wantC:   0,
			wantF:   32,
			wantK:   273,
		},
		{
			name:    "25 celsius",
			celsius: 25,
			wantC:   25,
			wantF:   77,
			wantK:   298,
		},
		{
			name:    "negative celsius",
			celsius: -40,
			wantC:   -40,
			wantF:   -40,
			wantK:   233,
		},
		{
			name:    "decimal celsius",
			celsius: 28.5,
			wantC:   28.5,
			wantF:   83.3,
			wantK:   301.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTemperature("Test City", tt.celsius)

			if got.C != tt.wantC {
				t.Errorf("Celsius: got %.1f, want %.1f", got.C, tt.wantC)
			}
			if got.F != tt.wantF {
				t.Errorf("Fahrenheit: got %.1f, want %.1f", got.F, tt.wantF)
			}
			if got.K != tt.wantK {
				t.Errorf("Kelvin: got %.1f, want %.1f", got.K, tt.wantK)
			}
			if got.City != "Test City" {
				t.Errorf("City: got %s, want %s", got.City, "Test City")
			}
		})
	}
}
