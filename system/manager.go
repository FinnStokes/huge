package system

import (
	"time"

	"../entity"
)

// Manager tracks all active systems and their update rates.
type Manager struct {
	systems      []System
	systemSpeeds [NumSpeeds][]System
}

// NewManager returns an initialised system manager.
func NewManager() *Manager {
	m := new(Manager)
	m.systemSpeeds[Slow] = make([]System, 0)
	m.systemSpeeds[Normal] = make([]System, 0)
	m.systemSpeeds[Fast] = make([]System, 0)
	return m
}

// AddSystem adds the specified system to the list of active systems at the specified speed.
func (m *Manager) AddSystem(speed Speed, system System) {
	m.systems = append(m.systems, system)
	m.systemSpeeds[speed] = append(m.systemSpeeds[speed], system)
}

// Update calls the Update method on all of the active systems at the specified speed.
func (m *Manager) Update(speed Speed, dt time.Duration, entities *entity.Manager) {
	for _, s := range m.systemSpeeds[speed] {
		s.Update(dt, entities)
	}
}

// Draw calls the Draw method on all of the active systems.
func (m *Manager) Draw(entities *entity.Manager) {
	for _, s := range m.systems {
		s.Draw(entities)
	}
}
