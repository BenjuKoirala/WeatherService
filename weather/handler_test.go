package weather

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//----------------------------------------------------------------------------------------------------------------------

// MockHTTPClient is a mock implementation of the HTTPClient interface for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Tests weatherHandler
func TestWeatherHandler(t *testing.T) {
	// Mocked HTTP response for testing
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"main":{"temp":25},"weather":[{"main":"Clear","description":"clear sky"}]}`))),
	}

	// Set up the mocked HTTP client
	mockHTTPClient := new(MockHTTPClient)
	mockHTTPClient.On("Get", mock.AnythingOfType("string")).Return(mockResponse, nil)

	// Override the default HTTPClient with the mock
	client = mockHTTPClient

	// Mock internal functions
	GetCoordinates = func(vars map[string]string) (float64, float64) {
		return 37.7749, -122.4194
	}
	GetWeatherData = func(lat, lon float64) (*WeatherResponse, error) {
		return &WeatherResponse{
			Weather: []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
			}{
				{
					Main:        "Clear",
					Description: "clear sky",
				},
			},
			Main: struct {
				Temp      float64 `json:"temp"`
				FeelsLike float64 `json:"feels_like"`
			}{
				Temp:      25.0,
				FeelsLike: 0.0,
			},
		}, nil
	}

	// Set up a request
	req, err := http.NewRequest("GET", "/weather/37.7749/-122.4194", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create a handler function and serve the HTTP request
	handler := http.HandlerFunc(WeatherHandler)
	handler.ServeHTTP(rr, req)

	// Assert the HTTP client was called as expected
	mockHTTPClient.AssertExpectations(t)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body
	expectedBody := `{"temperature":25,"weatherCondition":"Clear","description":"clear sky","weatherType":"Hot"}`
	assert.JSONEq(t, expectedBody, rr.Body.String())
}

//-----------------------------------------------------------------------------------------------------------------

// Tests categorizeWeather function
func TestCategorizeWeather(t *testing.T) {
	testCases := []struct {
		temperature float64
		expected    string
	}{
		{temperature: -5.0, expected: "Freezing"},
		{temperature: 5.0, expected: "Cold"},
		{temperature: 15.0, expected: "Moderate"},
		{temperature: 30.0, expected: "Hot"},
	}

	for _, tc := range testCases {
		t.Run(
			fmt.Sprintf("Temperature %f", tc.temperature),
			func(t *testing.T) {
				result := categorizeWeather(tc.temperature)
				assert.Equal(t, tc.expected, result)
			},
		)
	}
}

//----------------------------------------------------------------------------------------------------------------------
