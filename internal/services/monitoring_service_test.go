package services

import (
	"errors"
	"testing"

	"netwatch-api/internal/repository"
)

func TestCheckDeviceStatusOnline(t *testing.T) {
	store := repository.NewDeviceRepository()

	service := NewMonitoringService(store)

	// Simulate a successful TCP connection.
	service.checkConnection = func(ip string, port int) bool {
		return true
	}

	device, err := service.CheckDeviceStatus(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if device.Status != "online" {
		t.Errorf("expected status online, got %s", device.Status)
	}

	savedDevice, found := store.GetByID(1)
	if !found {
		t.Fatal("expected device to exist")
	}

	if savedDevice.Status != "online" {
		t.Errorf("expected saved status online, got %s", savedDevice.Status)
	}

}

func TestCheckDeviceStatusOffline(t *testing.T) {
	store := repository.NewDeviceRepository()

	service := NewMonitoringService(store)

	// Simulate a failed TCP connection.
	service.checkConnection = func(ip string, port int) bool {
		return false
	}

	device, err := service.CheckDeviceStatus(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if device.Status != "offline" {
		t.Errorf("expected status offline, got %s", device.Status)
	}

	savedDevice, found := store.GetByID(1)
	if !found {
		t.Fatal("expected device to exist")
	}

	if savedDevice.Status != "offline" {
		t.Errorf("expected saved status offline, got %s", savedDevice.Status)
	}
}

func TestCheckDeviceStatusNotFound(t *testing.T) {
	store := repository.NewDeviceRepository()
	service := NewMonitoringService(store)

	_, err := service.CheckDeviceStatus(999)

	if !errors.Is(err, ErrDeviceNotFound) {
		t.Fatalf("expected ErrDeviceNotFound, got %v", err)
	}
}

func TestCheckDeviceStatusUsesCorrectAddress(t *testing.T) {
	store := repository.NewDeviceRepository()
	service := NewMonitoringService(store)

	var receivedIP string
	var receivedPort int

	service.checkConnection = func(ip string, port int) bool {
		receivedIP = ip
		receivedPort = port
		return true
	}

	_, err := service.CheckDeviceStatus(1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if receivedIP != "192.168.255.102" {
		t.Errorf("expected IP 192.168.255.102, got %s", receivedIP)
	}

	if receivedPort != 80 {
		t.Errorf("expected port 80, got %d", receivedPort)
	}
}
