package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockViaCEPHandler struct {
	location string
	err      error
}

func (m *MockViaCEPHandler) GetLocation(zipcode string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.location, nil
}

type MockWeatherAPIHandler struct {
	temp float64
	err  error
}

func (m *MockWeatherAPIHandler) GetTemperature(location string) (float64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.temp, nil
}

func TestGetWeather_Success(t *testing.T) {
	mockViaCEP := &MockViaCEPHandler{location: "São Paulo"}
	mockWeather := &MockWeatherAPIHandler{temp: 28.5}

	handler := New(mockViaCEP, mockWeather)

	req := httptest.NewRequest(http.MethodGet, "/?zipcode=01310100", nil)
	w := httptest.NewRecorder()

	handler.PostWeather(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result["temp_C"] != 28.5 {
		t.Errorf("expected temp_C 28.5, got %v", result["temp_C"])
	}
	if result["temp_F"] != 83.3 {
		t.Errorf("expected temp_F 83.3, got %v", result["temp_F"])
	}
	if result["temp_K"] != 301.5 {
		t.Errorf("expected temp_K 301.5, got %v", result["temp_K"])
	}
}

func TestGetWeather_InvalidZipcode(t *testing.T) {
	handler := New(&MockViaCEPHandler{}, &MockWeatherAPIHandler{})

	tests := []struct {
		name    string
		zipcode string
	}{
		{"too short", "1234567"},
		{"empty", ""},
		{"invalid format", "123a5678"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/?zipcode="+tt.zipcode, nil)
			w := httptest.NewRecorder()

			handler.PostWeather(w, req)

			if w.Code != http.StatusUnprocessableEntity {
				t.Errorf("expected status 422, got %d", w.Code)
			}

			body, _ := io.ReadAll(w.Body)
			var result ErrorResponse
			json.Unmarshal(body, &result)

			if result.Message != "invalid zipcode" {
				t.Errorf("expected 'invalid zipcode', got %s", result.Message)
			}
		})
	}
}

func TestGetWeather_ZipcodeNotFound(t *testing.T) {
	mockViaCEP := &MockViaCEPHandler{err: ErrZipcodeNotFound}
	handler := New(mockViaCEP, &MockWeatherAPIHandler{})

	req := httptest.NewRequest(http.MethodGet, "/?zipcode=99999999", nil)
	w := httptest.NewRecorder()

	handler.PostWeather(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	body, _ := io.ReadAll(w.Body)
	var result ErrorResponse
	json.Unmarshal(body, &result)

	if result.Message != "can not find zipcode" {
		t.Errorf("expected 'can not find zipcode', got %s", result.Message)
	}
}

var ErrZipcodeNotFound = fmt.Errorf("zipcode not found")
