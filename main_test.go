package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"netwatch-api/internal/models"
	"netwatch-api/internal/repository"
	"netwatch-api/internal/services"
	"strings"
	"testing"
)

func createTestAPI() *API {
	deviceStore := repository.NewDeviceRepository()

	monitoringService := services.NewMonitoringServiceWithChecker(
		deviceStore,
		func(ip string, port int) bool {
			return true
		},
	)

	return &API{
		store:      deviceStore,
		monitoring: monitoringService,
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

func TestCheckDeviceStatus(t *testing.T) {
	api := createTestAPI()

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/devices/1/status",
		nil,
	)

	response := httptest.NewRecorder()

	api.registerRoutes().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	var result map[string]any

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result["status"] != "online" {
		t.Errorf("expected device status online, got %v", result["status"])
	}
}
