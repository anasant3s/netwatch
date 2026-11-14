package services

import (
	"errors"
	"net"
	"strconv"
	"time"

	"netwatch-api/internal/models"
	"netwatch-api/internal/repository"
)

var ErrDeviceNotFound = errors.New("device not found")

type MonitoringService struct {
	store           *repository.DeviceRepository
	checkConnection func(ip string, port int) bool
}

func NewMonitoringService(store *repository.DeviceRepository) *MonitoringService {
	return NewMonitoringServiceWithChecker(store, checkTCPConnection)
}

func NewMonitoringServiceWithChecker(
	store *repository.DeviceRepository,
	checker func(ip string, port int) bool,
) *MonitoringService {
	return &MonitoringService{
		store:           store,
		checkConnection: checker,
	}
}

// CheckDeviceStatus checks TCP connectivity and updates the device status.
func (service *MonitoringService) CheckDeviceStatus(deviceID int) (models.Device, error) {
	device, found := service.store.GetByID(deviceID)
	if !found {
		return models.Device{}, ErrDeviceNotFound
	}

	status := "offline"

	if service.checkConnection(device.IP, device.Port) {
		status = "online"
	}

	updated := service.store.UpdateStatus(deviceID, status)
	if !updated {
		return models.Device{}, ErrDeviceNotFound
	}

	device.Status = status

	return device, nil
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
