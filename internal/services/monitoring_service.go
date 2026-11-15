package services

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	"netwatch-api/internal/models"
	"netwatch-api/internal/repository"
)

var ErrDeviceNotFound = errors.New("device not found")

const defaultCheckTimeout = 3 * time.Second

type MonitoringService struct {
	store           *repository.DeviceRepository
	checkConnection func(ip string, port int) bool
}

func NewMonitoringService(store *repository.DeviceRepository) *MonitoringService {
	return NewMonitoringServiceWithTimeout(store, defaultCheckTimeout)
}

func NewMonitoringServiceWithTimeout(
	store *repository.DeviceRepository,
	timeout time.Duration,
) *MonitoringService {
	if timeout <= 0 {
		timeout = defaultCheckTimeout
	}

	return NewMonitoringServiceWithChecker(store, func(ip string, port int) bool {
		return checkTCPConnection(ip, port, timeout)
	})
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

// CheckAllDevices checks the connectivity of every registered device.
func (service *MonitoringService) CheckAllDevices() {
	devices := service.store.List()

	for _, device := range devices {
		_, err := service.CheckDeviceStatus(device.ID)
		if err != nil {
			log.Printf(
				"failed to check device %d: %v",
				device.ID,
				err,
			)
		}
	}
}

// StartAutomaticMonitoring periodically checks all registered devices.
func (service *MonitoringService) StartAutomaticMonitoring(
	ctx context.Context,
	interval time.Duration,
) error {
	if interval <= 0 {
		return errors.New("monitoring interval must be greater than zero")
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Perform the first check immediately.
		service.CheckAllDevices()

		for {
			select {
			case <-ctx.Done():
				log.Println("automatic monitoring stopped")
				return

			case <-ticker.C:
				service.CheckAllDevices()
			}
		}
	}()

	return nil
}

// checkTCPConnection tests whether a TCP port accepts connections.
func checkTCPConnection(ip string, port int, timeout time.Duration) bool {
	address := net.JoinHostPort(ip, strconv.Itoa(port))

	connection, err := net.DialTimeout(
		"tcp",
		address,
		timeout,
	)
	if err != nil {
		return false
	}

	_ = connection.Close()
	return true
}
