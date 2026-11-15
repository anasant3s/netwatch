package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"netwatch-api/internal/config"
	"netwatch-api/internal/models"
	"netwatch-api/internal/repository"
	"netwatch-api/internal/services"
)

type API struct {
	store      *repository.DeviceRepository
	incidents  *repository.IncidentRepository
	monitoring *services.MonitoringService
}

// writeJSON sends a JSON response with the given HTTP status.
func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// writeError sends an error message in JSON format.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{
		"error": message,
	})
}

// decodeJSON reads a single JSON object from the request body.
func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	// Reject additional JSON values after the first object.
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("request must contain one JSON object")
	}

	return nil
}

// validateDevice checks the fields required to register a device.
func validateDevice(device models.Device) error {
	if strings.TrimSpace(device.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(device.Type) == "" {
		return errors.New("type is required")
	}

	if net.ParseIP(device.IP) == nil {
		return errors.New("a valid IP address is required")
	}

	if device.Port < 1 || device.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}

	return nil
}

// GET /api/health
func (api *API) healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (api *API) listDevices(w http.ResponseWriter, r *http.Request) {
	deviceList := api.store.List()

	writeJSON(w, http.StatusOK, deviceList)
}

func (api *API) listDeviceIncidents(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || deviceID < 1 {
		writeError(w, http.StatusBadRequest, "invalid device ID")
		return
	}

	_, found := api.store.GetByID(deviceID)
	if !found {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	incidents := api.incidents.ListByDeviceID(deviceID)

	writeJSON(w, http.StatusOK, incidents)
}

func (api *API) getDevice(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || deviceID < 1 {
		writeError(w, http.StatusBadRequest, "invalid device ID")
		return
	}

	device, found := api.store.GetByID(deviceID)
	if !found {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	writeJSON(w, http.StatusOK, device)
}

func (api *API) createDevice(w http.ResponseWriter, r *http.Request) {
	var newDevice models.Device

	if err := decodeJSON(r, &newDevice); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := validateDevice(newDevice); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	newDevice.Name = strings.TrimSpace(newDevice.Name)
	newDevice.Type = strings.ToLower(strings.TrimSpace(newDevice.Type))

	createdDevice := api.store.Create(newDevice)

	writeJSON(w, http.StatusCreated, createdDevice)
}

func (api *API) updateDevice(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || deviceID < 1 {
		writeError(w, http.StatusBadRequest, "invalid device ID")
		return
	}

	var updatedDevice models.Device

	if err := decodeJSON(r, &updatedDevice); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := validateDevice(updatedDevice); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedDevice.Name = strings.TrimSpace(updatedDevice.Name)
	updatedDevice.Type = strings.ToLower(strings.TrimSpace(updatedDevice.Type))

	savedDevice, found := api.store.Update(deviceID, updatedDevice)
	if !found {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	writeJSON(w, http.StatusOK, savedDevice)
}

// DELETE /api/devices/{id}
func (api *API) deleteDevice(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || deviceID < 1 {
		writeError(w, http.StatusBadRequest, "invalid device ID")
		return
	}

	deleted := api.store.Delete(deviceID)
	if !deleted {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// checkTCPConnection tests whether a TCP port accepts connections.
func checkTCPConnection(ip string, port int) bool {
	address := net.JoinHostPort(ip, strconv.Itoa(port))

	connection, err := net.DialTimeout(
		"tcp",
		address,
		3*time.Second,
	)
	if err != nil {
		return false
	}

	_ = connection.Close()
	return true
}

func (api *API) checkDeviceStatus(w http.ResponseWriter, r *http.Request) {
	deviceID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || deviceID < 1 {
		writeError(w, http.StatusBadRequest, "invalid device ID")
		return
	}

	device, err := api.monitoring.CheckDeviceStatus(deviceID)
	if errors.Is(err, services.ErrDeviceNotFound) {
		writeError(w, http.StatusNotFound, "device not found")
		return
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to check device status")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":     device.ID,
		"ip":     device.IP,
		"port":   device.Port,
		"status": device.Status,
	})
}

// registerRoutes maps HTTP methods and paths to their handlers.
func (api *API) registerRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", api.healthCheck)
	mux.HandleFunc("GET /api/devices", api.listDevices)
	mux.HandleFunc("GET /api/devices/{id}/incidents", api.listDeviceIncidents)
	mux.HandleFunc("POST /api/devices", api.createDevice)

	mux.HandleFunc("GET /api/devices/{id}", api.getDevice)
	mux.HandleFunc("PUT /api/devices/{id}", api.updateDevice)
	mux.HandleFunc("DELETE /api/devices/{id}", api.deleteDevice)

	mux.HandleFunc("GET /api/devices/{id}/status", api.checkDeviceStatus)

	return mux
}

func main() {
	appConfig, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	deviceStore := repository.NewDeviceRepository()
	incidentStore := repository.NewIncidentRepository()

	monitoringService := services.NewMonitoringServiceWithDependencies(
		deviceStore,
		incidentStore,
		appConfig.CheckTimeout,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := monitoringService.StartAutomaticMonitoring(
		ctx,
		appConfig.CheckInterval,
	); err != nil {
		log.Fatalf("failed to start monitoring: %v", err)
	}

	api := &API{
		store:      deviceStore,
		incidents:  incidentStore,
		monitoring: monitoringService,
	}

	server := &http.Server{
		Addr:              ":8080",
		Handler:           api.registerRoutes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Println("NetWatch API running at http://localhost:8080")
	fmt.Printf(
		"Automatic monitoring started (interval: %s, timeout: %s)\n",
		appConfig.CheckInterval,
		appConfig.CheckTimeout,
	)

	log.Fatal(server.ListenAndServe())
}
