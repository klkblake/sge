package sge

import (
	"atom/sdl"
)

type KeyPresser interface {
	KeyPress(keysym *sdl.Keysym)
}

type KeyReleaser interface {
	KeyRelease(keysym *sdl.Keysym)
}

type Events struct {
	KeyPressers []KeyPresser
	KeyReleasers []KeyReleaser
}

func NewEvents() *Events {
	events := new(Events)
	events.KeyPressers = make([]KeyPresser, 0)
	events.KeyReleasers = make([]KeyReleaser, 0)
	return events
}

func (events *Events) AddKeyPresser(keyPresser KeyPresser) {
	events.KeyPressers = append(events.KeyPressers, keyPresser)
}

func (events *Events) AddKeyReleaser(keyReleaser KeyReleaser) {
	events.KeyReleasers = append(events.KeyReleasers, keyReleaser)
}

func (events *Events) Dispatch(event interface{}) {
	switch e := event.(type) {
	case sdl.KeyboardEvent:
		switch e.Type {
		case sdl.KEYDOWN:
			for _, keyPresser := range events.KeyPressers {
				keyPresser.KeyPress(&e.Keysym)
			}
		case sdl.KEYUP:
			for _, keyReleaser := range events.KeyReleasers {
				keyReleaser.KeyRelease(&e.Keysym)
			}
		}
	}
}
