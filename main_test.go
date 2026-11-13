package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"netwatch-api/internal/models"
)

// createTestAPI creates a fresh API for each test.
func createTestAPI() *API {
	return &API{
		store: NewDeviceStore(),
	}
}

// performRequest simulates an HTTP request.
func performRequest(handler http.Handler, method, url, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, url, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	return response
}

// TestHealthCheck checks if the API is running.
func TestHealthCheck(t *testing.T) {
	api := createTestAPI()

	response := performRequest(
		api.registerRoutes(),
		http.MethodGet,
		"/api/health",
		"",
	)

	if response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", response.Code)
	}
}

// TestListDevices checks if all devices are returned.
func TestListDevices(t *testing.T) {
	api := createTestAPI()

	response := performRequest(
		api.registerRoutes(),
		http.MethodGet,
		"/api/devices",
		"",
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var devices []models.Device

	if err := json.Unmarshal(response.Body.Bytes(), &devices); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(devices) != 2 {
		t.Errorf("expected 2 devices, got %d", len(devices))
	}
}

// TestCreateDevice checks if a new device can be registered.
func TestCreateDevice(t *testing.T) {
	api := createTestAPI()

	requestBody := `{
		"name": "Network Switch",
		"ip": "192.168.255.110",
		"type": "switch",
		"port": 80
	}`

	response := performRequest(
		api.registerRoutes(),
		http.MethodPost,
		"/api/devices",
		requestBody,
	)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}

	var createdDevice models.Device

	if err := json.Unmarshal(response.Body.Bytes(), &createdDevice); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdDevice.ID != 3 {
		t.Errorf("expected ID 3, got %d", createdDevice.ID)
	}

	if createdDevice.Name != "Network Switch" {
		t.Errorf("unexpected device name: %s", createdDevice.Name)
	}

	if createdDevice.Status != "unknown" {
		t.Errorf("expected status unknown, got %s", createdDevice.Status)
	}
}

// TestCreateDeviceWithInvalidIP checks IP validation.
func TestCreateDeviceWithInvalidIP(t *testing.T) {
	api := createTestAPI()

	requestBody := `{
		"name": "Invalid Device",
		"ip": "999.999.999.999",
		"type": "camera",
		"port": 80
	}`

	response := performRequest(
		api.registerRoutes(),
		http.MethodPost,
		"/api/devices",
		requestBody,
	)

	if response.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", response.Code)
	}
}

// TestGetDevice checks if a device can be found by ID.
func TestGetDevice(t *testing.T) {
	api := createTestAPI()

	response := performRequest(
		api.registerRoutes(),
		http.MethodGet,
		"/api/devices/1",
		"",
	)

	if response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", response.Code)
	}
}

// TestGetNonexistentDevice checks the response for an unknown ID.
func TestGetNonexistentDevice(t *testing.T) {
	api := createTestAPI()

	response := performRequest(
		api.registerRoutes(),
		http.MethodGet,
		"/api/devices/999",
		"",
	)

	if response.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", response.Code)
	}
}

// TestUpdateDevice checks if an existing device can be updated.
func TestUpdateDevice(t *testing.T) {
	api := createTestAPI()

	requestBody := `{
		"name": "Updated NVR",
		"ip": "192.168.255.102",
		"type": "nvr",
		"port": 8080
	}`

	response := performRequest(
		api.registerRoutes(),
		http.MethodPut,
		"/api/devices/1",
		requestBody,
	)

	if response.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", response.Code)
	}
}

// TestDeleteDevice checks if a device can be deleted.
func TestDeleteDevice(t *testing.T) {
	api := createTestAPI()

	response := performRequest(
		api.registerRoutes(),
		http.MethodDelete,
		"/api/devices/1",
		"",
	)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}

	// Confirm that the device was actually deleted.
	getResponse := performRequest(
		api.registerRoutes(),
		http.MethodGet,
		"/api/devices/1",
		"",
	)

	if getResponse.Code != http.StatusNotFound {
		t.Errorf("expected status 404 after deletion, got %d", getResponse.Code)
	}
}
