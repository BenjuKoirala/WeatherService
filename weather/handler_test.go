// handler_test.go
package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//----------------------------------------------------------------------------------------------------------------------

// initServer initializes the server with the provided router
func initServer(router *mux.Router) {
	// Initialize the server with the provided router
	http.Handle("/", router)
}

// MockHTTPClient is a mock implementation of the HTTPClient interface for testing
type MockHTTPClient struct {
	mock.Mock
}

// Get is a mocked implementation of the HTTP GET method
func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

var httpClient = new(MockHTTPClient) // Mock the HTTP client for testing

func init() {
	// Override the default HTTPClient with the mock
	client = httpClient
}

//----------------------------------------------------------------------------------------------------------------------

// Tests weatherHandler function
func TestWeatherHandler(t *testing.T) {
	// Mocked HTTP response for testing
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(
			bytes.NewBufferString(`{"main":{"temp":25},"weather":[{"main":"Clear","description":"clear sky"}]}`),
		),
	}

	// Set expectations for the HTTP client
	httpClient.On("Get", mock.Anything).Return(mockResponse, nil).Once()

	req, err := http.NewRequest("GET", "/weather/37.7749/-122.4194", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	handler.ServeHTTP(rr, req)

	// Assert the HTTP client was called as expected
	httpClient.AssertExpectations(t)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the response body
	expectedBody := `{"temperature":25,"weatherCondition":"Clear","description":"clear sky","weatherType":"Moderate"}`
	assert.JSONEq(t, expectedBody, rr.Body.String())
}

//----------------------------------------------------------------------------------------------------------------------

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
