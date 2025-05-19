package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const MOCK_RESPONSE = `{
			"location": {
				"name": "Kyiv",
				"region": "Kyyivs'ka Oblast'",
				"country": "Ukraine",
				"lat": 50.4333,
				"lon": 30.5167,
				"tz_id": "Europe/Kiev",
				"localtime_epoch": 1747570235,
				"localtime": "2025-05-18 15:10"
			},
			"current": {
				"temp_c": 15.9,
				"condition": {
					"text": "Patchy rain nearby"
				},
				"humidity": 52
			}
		}`

func TestFetchWeatherNow(t *testing.T) {
	// Create mock server to get fake JSON response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(MOCK_RESPONSE)); err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer mockServer.Close()

	queryParams := map[string]string{
		"q":   "Kyiv",
		"key": "fake-key",
	}

	result, err := FetchWeatherNow(mockServer.URL, "fake-key", queryParams)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	desc, ok := (*result)["description"].(string)
	if !ok {
		t.Errorf("Description is not a string. Got %T", (*result)["description"])
	} else if len(desc) == 0 {
		t.Errorf("Description is empty.")
	}

	if _, ok := (*result)["humidity"]; !ok {
		t.Errorf("Humidity not found in result")
	}

	if _, ok := (*result)["temperature"]; !ok {
		t.Errorf("Temperature not found in result")
	}
}
