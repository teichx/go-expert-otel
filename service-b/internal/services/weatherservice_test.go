package services

import (
	"fmt"
	"testing"
)

type MockViaCEP struct {
	location string
	err      error
}

func (m *MockViaCEP) GetLocation(zipcode string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.location, nil
}

type MockWeatherAPI struct {
	temp float64
	err  error
}

func (m *MockWeatherAPI) GetTemperature(location string) (float64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.temp, nil
}

func TestGetWeatherByZipcode_Success(t *testing.T) {
	mockViaCEP := &MockViaCEP{location: "São Paulo"}
	mockWeather := &MockWeatherAPI{temp: 25}

	service := New(mockViaCEP, mockWeather)
	result, err := service.GetWeatherByZipcode(t.Context(), "01310100")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.C != 25 {
		t.Errorf("expected 25°C, got %.1f", result.C)
	}
	if result.F != 77 {
		t.Errorf("expected 77°F, got %.1f", result.F)
	}
	if result.K != 298 {
		t.Errorf("expected 298K, got %.1f", result.K)
	}
}

func TestGetWeatherByZipcode_InvalidZipcode(t *testing.T) {
	mockViaCEP := &MockViaCEP{}
	mockWeather := &MockWeatherAPI{}

	service := New(mockViaCEP, mockWeather)

	tests := []struct {
		name    string
		zipcode string
	}{
		{"too short", "1234567"},
		{"too long", "123456789"},
		{"with letters", "1234567a"},
		{"empty", ""},
		{"with space", "12345 678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetWeatherByZipcode(t.Context(), tt.zipcode)
			if err == nil {
				t.Error("expected error for invalid zipcode")
			}
			if err.Error() != "invalid zipcode" {
				t.Errorf("expected 'invalid zipcode', got %v", err)
			}
		})
	}
}

func TestGetWeatherByZipcode_ZipcodeNotFound(t *testing.T) {
	mockViaCEP := &MockViaCEP{err: fmt.Errorf("not found")}
	mockWeather := &MockWeatherAPI{}

	service := New(mockViaCEP, mockWeather)
	_, err := service.GetWeatherByZipcode(t.Context(), "12345678")

	if err == nil {
		t.Fatal("expected error for zipcode not found")
	}
	if err.Error() != "can not find zipcode" {
		t.Errorf("expected 'can not find zipcode', got %v", err)
	}
}

func TestGetWeatherByZipcode_WeatherAPIError(t *testing.T) {
	mockViaCEP := &MockViaCEP{location: "São Paulo"}
	mockWeather := &MockWeatherAPI{err: fmt.Errorf("API error")}

	service := New(mockViaCEP, mockWeather)
	_, err := service.GetWeatherByZipcode(t.Context(), "01310100")

	if err == nil {
		t.Fatal("expected error from weather API")
	}
}
