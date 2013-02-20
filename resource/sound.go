package resource

import (
	"log"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/mkb218/gosndfile/sndfile"
)

type playState int

const (
	stopped playState = iota
	playing
	paused
)

// A Sound provides access to an audio stream. A sound is initially stopped.
// If the Looping field is true, the sound plays again from the beginning on completion.
// The Volume field sets an amplification factor for the sound.
type Sound struct {
	Looping bool
	Volume  float32
	state   playState
	file    *sndfile.File
	stream  *portaudio.Stream
}

func newSound(name string) (s *Sound, err error) {
	s = new(Sound)
	s.Volume = 1.0
	s.file, err = sndfile.Open(name, sndfile.Read, &sndfile.Info{})
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Play starts a stopped or paused source playing.
func (s *Sound) Play() {
	if s.state != playing {
		var err error
		if s.stream == nil {
			s.stream, err = portaudio.OpenDefaultStream(0, (int)(s.file.Format.Channels),
				(float64)(s.file.Format.Samplerate), 4096, s)
			if err != nil {
				log.Println("portaudio output creation failed", err)
				return
			}
		}
		err = s.stream.Start()
		if err != nil {
			log.Println("portaudio output start failed", err)
			return
		}
		s.state = playing
	}
}

func (s *Sound) ProcessAudio(_, out []float32) {
	buff := out
	for len(buff) > 0 {
		n, err := s.file.ReadItems(buff)
		if err != nil {
			log.Println("reading sound file failed", err)
		}
		buff = buff[n:]

		if s.Looping && len(buff) > 0 {
			s.Rewind()
		} else {
			for i := range buff {
				buff[i] = 0
			}
			buff = buff[len(buff):]
		}
	}
	for i, v := range out {
		out[i] = s.Volume * v
	}
}

// Pause pauses a currently playing sound.
func (s *Sound) Pause() {
	if s.state == playing {
		s.stream.Stop()
		s.state = paused
	}
}

// Stop stops and rewinds a playing or paused sound.
func (s *Sound) Stop() {
	if s.state != stopped {
		s.stream.Stop()
		s.Rewind()
		s.state = stopped
	}
}

// Resume starts a paused sound playing.
func (s *Sound) Resume() {
	if s.state == paused {
		s.Play()
	}
}

// Rewind sets the surrently playing position to the start of the sound.
func (s *Sound) Rewind() {
	s.file.Seek(0, sndfile.Set)
}

// Seek sets the currently playing position of the sound.
func (s *Sound) Seek(offset time.Duration) {
	frame := (int64)(time.Duration(s.file.Format.Samplerate) * offset / time.Second)
	s.file.Seek(frame, sndfile.Set)
}

// Tell gets the currently playing position of the sound.
func (s *Sound) Tell() (offset time.Duration) {
	frame, _ := s.file.Seek(0, sndfile.Current)
	offset = time.Duration(frame) * time.Second / time.Duration(s.file.Format.Samplerate)
	return
}

// Len gets the entire duration of the sound.
func (s *Sound) Len() (offset time.Duration) {
	frame, _ := s.file.Seek(0, sndfile.End)
	offset = time.Duration(frame) * time.Second / time.Duration(s.file.Format.Samplerate)
	return
}
