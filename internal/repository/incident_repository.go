package repository

import (
	"sync"

	"netwatch-api/internal/models"
)

type IncidentRepository struct {
	mu        sync.RWMutex
	incidents []models.Incident
	nextID    int
}

func NewIncidentRepository() *IncidentRepository {
	return &IncidentRepository{
		incidents: make([]models.Incident, 0),
		nextID:    1,
	}
}

// Create registers a new incident.
func (r *IncidentRepository) Create(incident models.Incident) models.Incident {
	r.mu.Lock()
	defer r.mu.Unlock()

	incident.ID = r.nextID

	r.incidents = append(r.incidents, incident)
	r.nextID++

	return incident
}

// ListByDeviceID returns all incidents associated with a device.
func (r *IncidentRepository) ListByDeviceID(deviceID int) []models.Incident {
	r.mu.RLock()
	defer r.mu.RUnlock()

	deviceIncidents := make([]models.Incident, 0)

	for _, incident := range r.incidents {
		if incident.DeviceID == deviceID {
			deviceIncidents = append(deviceIncidents, incident)
		}
	}

	return deviceIncidents
}
