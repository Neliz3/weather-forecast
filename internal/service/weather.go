package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type WeatherResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TZID           string  `json:"tz_id"`
		LocaltimeEpoch int64   `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		Humidity int `json:"humidity"`
	} `json:"current"`
}

func fetchRawWeatherData(apiURL string, queryParams map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}()

	return io.ReadAll(resp.Body)
}

func FetchWeatherNow(apiURL string, queryParams map[string]string) (*map[string]any, error) {
	data, err := fetchRawWeatherData(apiURL, queryParams)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(data, &weather); err != nil {
		return nil, err
	}

	res := map[string]any{
		"temperature": weather.Current.TempC,
		"humidity":    weather.Current.Humidity,
		"description": weather.Current.Condition.Text,
	}
	return &res, nil
}
