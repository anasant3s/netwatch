package repository

import (
	"sync"

	"netwatch-api/internal/models"
)

// DeviceRepository manages device data.
type DeviceRepository struct {
	mu      sync.RWMutex
	devices []models.Device
	nextID  int
}

// NewDeviceRepository creates a repository with initial devices.
func NewDeviceRepository() *DeviceRepository {
	return &DeviceRepository{
		devices: []models.Device{
			{
				ID:     1,
				Name:   "DHI-NVR5864-XI/Pro",
				IP:     "192.168.255.102",
				Type:   "nvr",
				Port:   80,
				Status: "unknown",
			},
			{
				ID:     2,
				Name:   "DH-IPC-HFW7459Z1-Z-PV-X",
				IP:     "192.168.255.105",
				Type:   "camera",
				Port:   80,
				Status: "unknown",
			},
		},
		nextID: 3,
	}
}

// List returns all registered devices.
func (r *DeviceRepository) List() []models.Device {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return append([]models.Device(nil), r.devices...)
}

// GetByID returns a device by its ID.
func (r *DeviceRepository) GetByID(id int) (models.Device, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, device := range r.devices {
		if device.ID == id {
			return device, true
		}
	}

	return models.Device{}, false
}

// Create registers a new device.
func (r *DeviceRepository) Create(device models.Device) models.Device {
	r.mu.Lock()
	defer r.mu.Unlock()

	device.ID = r.nextID
	device.Status = "unknown"

	r.devices = append(r.devices, device)
	r.nextID++

	return device
}

// Update modifies an existing device.
func (r *DeviceRepository) Update(id int, updatedDevice models.Device) (models.Device, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, device := range r.devices {
		if device.ID == id {
			updatedDevice.ID = id
			updatedDevice.Status = device.Status
			r.devices[i] = updatedDevice

			return updatedDevice, true
		}
	}

	return models.Device{}, false
}

// UpdateStatus updates the status of a device.
func (r *DeviceRepository) UpdateStatus(id int, status string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, device := range r.devices {
		if device.ID == id {
			r.devices[i].Status = status
			return true
		}
	}

	return false
}

// Delete removes a device by its ID.
func (r *DeviceRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, device := range r.devices {
		if device.ID == id {
			r.devices = append(r.devices[:i], r.devices[i+1:]...)
			return true
		}
	}

	return false
}

// UpdateStatusIfChanged updates a device status and returns its previous status.
func (r *DeviceRepository) UpdateStatusIfChanged(
	id int,
	status string,
) (string, bool, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, device := range r.devices {
		if device.ID == id {
			previousStatus := device.Status

			if previousStatus == status {
				return previousStatus, false, true
			}

			r.devices[i].Status = status

			return previousStatus, true, true
		}
	}

	return "", false, false
}
