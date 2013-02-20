// Package resource manages loading and accessing game resources.
package resource

// Manager is a type that stores the loaded resources and allows access with automatic loading.
type Manager struct {
	music  map[string]*Sound
	sounds map[string]*Sound
}

// NewManager returns an initialised resource manager.
func NewManager() *Manager {
	m := new(Manager)
	m.music = make(map[string]*Sound)
	m.sounds = make(map[string]*Sound)
	return m
}

// GetSound fetches an audio file in the .wav format
func (m *Manager) GetSound(name string) (s *Sound, err error) {
	s, ok := m.sounds[name]
	if !ok {
		s, err = newSound(name + ".wav")
		if err != nil {
			return nil, err
		}
		m.sounds[name] = s
	}
	return s, nil
}

// GetMusic fetches an audio file in the .ogg format and sets it to loop
func (m *Manager) GetMusic(name string) (s *Sound, err error) {
	s, ok := m.music[name]
	if !ok {
		s, err = newSound(name + ".ogg")
		if err != nil {
			return nil, err
		}
		s.Looping = true
		m.music[name] = s
	}
	return s, nil
}
