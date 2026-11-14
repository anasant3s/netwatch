package repository

import (
	"testing"

	"netwatch-api/internal/models"
)

func TestCreateDevice(t *testing.T) {
	store := NewDeviceRepository()

	device := models.Device{
		Name: "Switch - NetWatch",
		IP:   "192.168.255.110",
		Type: "switch",
		Port: 80,
	}

	createdDevice := store.Create(device)

	if createdDevice.ID != 3 {
		t.Errorf("expected ID 3, got %d", createdDevice.ID)
	}

	if createdDevice.Status != "unknown" {
		t.Errorf("expected status unknown, got %s", createdDevice.Status)
	}
}

func TestGetDeviceByID(t *testing.T) {
	store := NewDeviceRepository()

	device, found := store.GetByID(1)

	if !found {
		t.Fatal("expected device to be found")
	}

	if device.ID != 1 {
		t.Errorf("expected ID 1, got %d", device.ID)
	}

	if device.Name != "DHI-NVR5864-XI/Pro" {
		t.Errorf("expected NVR name, got %s", device.Name)
	}
}

func TestGetDeviceByIDNotFound(t *testing.T) {
	store := NewDeviceRepository()

	_, found := store.GetByID(999)

	if found {
		t.Error("expected device not to be found")
	}
}

func TestUpdateDevice(t *testing.T) {
	store := NewDeviceRepository()

	originalDevice, _ := store.GetByID(1)
	originalDevice.Name = "NVR Atualizado"
	originalDevice.IP = "192.168.255.120"

	updatedDevice, found := store.Update(1, originalDevice)

	if !found {
		t.Fatal("expected device to be updated")
	}

	if updatedDevice.ID != 1 {
		t.Errorf("expected ID 1, got %d", updatedDevice.ID)
	}

	if updatedDevice.Name != "NVR Atualizado" {
		t.Errorf("expected updated name, got %s", updatedDevice.Name)
	}

	if updatedDevice.IP != "192.168.255.120" {
		t.Errorf("expected updated IP, got %s", updatedDevice.IP)
	}
}

func TestUpdateDeviceNotFound(t *testing.T) {
	store := NewDeviceRepository()

	device := models.Device{
		Name: "Dispositivo inexistente",
		IP:   "192.168.255.150",
		Type: "switch",
		Port: 80,
	}

	_, found := store.Update(999, device)

	if found {
		t.Error("expected device not to be updated")
	}
}

func TestUpdateDeviceStatus(t *testing.T) {
	store := NewDeviceRepository()
	updated := store.UpdateStatus(1, "online")
	if !updated {
		t.Fatal("expected device status to be updated")
	}

	device, found := store.GetByID(1)

	if !found {
		t.Fatal("expected device to exist")
	}

	if device.Status != "online" {
		t.Errorf("expected status online, got %s", device.Status)
	}

	if device.Name != "DHI-NVR5864-XI/Pro" {
		t.Errorf("device name should not change, got %s", device.Name)
	}
}
func TestUpdateDeviceStatusNotFound(t *testing.T) {
	store := NewDeviceRepository()
	updated := store.UpdateStatus(999, "online")
	if updated {
		t.Error("expected status not to be updated")
	}
}

func TestDeleteDevice(t *testing.T) {
	store := NewDeviceRepository()
	deleted := store.Delete(1)
	if !deleted {
		t.Fatal("expected device to be deleted")
	}

	_, found := store.GetByID(1)

	if found {
		t.Error("expected device not to exist after deletion")
	}
}
func TestDeleteDeviceNotFound(t *testing.T) {
	store := NewDeviceRepository()
	deleted := store.Delete(999)
	if deleted {
		t.Error("expected device not to be deleted")
	}
}
