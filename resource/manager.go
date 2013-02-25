// Package resource manages loading and accessing game resources.
package resource

import (
	"encoding/json"
	"image"
	_ "image/png"
	"os"

	"github.com/FinnStokes/huge/sprite"
)

// Manager is a type that stores the loaded resources and allows access with automatic loading.
type Manager struct {
	music   map[string]*Sound
	sounds  map[string]*Sound
	images  map[string]image.Image
	json    map[string]*os.File
	sprites map[string]*spriteSpec
}

// NewManager returns an initialised resource manager.
func NewManager() *Manager {
	m := new(Manager)
	m.music = make(map[string]*Sound)
	m.sounds = make(map[string]*Sound)
	m.images = make(map[string]image.Image)
	m.json = make(map[string]*os.File)
	m.sprites = make(map[string]*spriteSpec)
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

// GetImage fetches and decodes an image file in the .png format
func (m *Manager) GetImage(name string) (img image.Image, err error) {
	img, ok := m.images[name]
	if !ok {
		file, err := os.Open(name + ".png")
		if err != nil {
			return nil, err
		}
		img, _, err = image.Decode(file)
		if err != nil {
			return nil, err
		}
		m.images[name] = img
	}
	return img, nil
}

// GetJson fetches a json file and loads it into the given struct
func (m *Manager) GetJson(name string, target interface{}) (err error) {
	file, ok := m.json[name]
	if !ok {
		file, err = os.Open(name + ".json")
		if err != nil {
			return err
		}
		m.json[name] = file
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(target)
	if err != nil {
		return err
	}
	return nil
}

// GetSprite creates a sprite based on the specification in a .json file
func (m *Manager) GetSprite(name string) (s *sprite.Sprite, err error) {
	sprite, ok := m.sprites[name]
	if !ok {
		sprite = new(spriteSpec)
		err = m.GetJson(name, sprite)
		if err != nil {
			return nil, err
		}
		m.sprites[name] = sprite
	}
	return sprite.New(m)
}
