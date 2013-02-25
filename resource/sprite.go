package resource

import "github.com/FinnStokes/huge/sprite"

type spriteSpec struct {
	Image      string
	Width      int
	Height     int
	Animations map[string]animationSpec
	Playing    string
}

type animationSpec struct {
	Frames []int
	Fps    int
	Next   string
}

func (s *spriteSpec) New(m *Manager) (*sprite.Sprite, error) {
	animations := make(map[string]*sprite.Animation, len(s.Animations))
	for k, a := range s.Animations {
		animations[k] = &sprite.Animation{
			a.Frames,
			a.Fps,
			nil,
		}
	}
	for k, a := range s.Animations {
		animations[k].Next = animations[a.Next]
	}
	img, err := m.GetImage(s.Image)
	if err != nil {
		return nil, err
	}
	return &sprite.Sprite{
		img,
		animations,
		s.Width, s.Height,
		animations[s.Playing],
		0,
		0,
	}, nil
}
