package sge

import (
	"github.com/klkblake/Go-SDL/sdl"
)

type KeyPresser interface {
	KeyPress(keysym *sdl.Keysym)
}

type KeyReleaser interface {
	KeyRelease(keysym *sdl.Keysym)
}

type MousePresser interface {
	MousePress(x, y uint16, button uint8)
}

type MouseReleaser interface {
	MouseRelease(x, y uint16, button uint8)
}

type MouseMover interface {
	MouseMove(x, y uint16, xrel, yrel int16)
}

type Events struct {
	KeyPressers    []KeyPresser
	KeyReleasers   []KeyReleaser
	MousePressers  []MousePresser
	MouseReleasers []MouseReleaser
	MouseMovers    []MouseMover
}

func NewEvents() *Events {
	events := new(Events)
	events.KeyPressers = make([]KeyPresser, 0)
	events.KeyReleasers = make([]KeyReleaser, 0)
	events.MousePressers = make([]MousePresser, 0)
	events.MouseReleasers = make([]MouseReleaser, 0)
	events.MouseMovers = make([]MouseMover, 0)
	return events
}

func (events *Events) AddKeyPresser(keyPresser KeyPresser) {
	events.KeyPressers = append(events.KeyPressers, keyPresser)
}

func (events *Events) AddKeyReleaser(keyReleaser KeyReleaser) {
	events.KeyReleasers = append(events.KeyReleasers, keyReleaser)
}

func (events *Events) AddMousePresser(mousePresser MousePresser) {
	events.MousePressers = append(events.MousePressers, mousePresser)
}

func (events *Events) AddMouseReleaser(mouseReleaser MouseReleaser) {
	events.MouseReleasers = append(events.MouseReleasers, mouseReleaser)
}

func (events *Events) AddMouseMover(mouseMover MouseMover) {
	events.MouseMovers = append(events.MouseMovers, mouseMover)
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
	case sdl.MouseButtonEvent:
		switch e.Type {
		case sdl.MOUSEBUTTONDOWN:
			for _, mousePresser := range events.MousePressers {
				mousePresser.MousePress(e.X, e.Y, e.Button)
			}
		case sdl.MOUSEBUTTONUP:
			for _, mouseReleaser := range events.MouseReleasers {
				mouseReleaser.MouseRelease(e.X, e.Y, e.Button)
			}
		}
	case sdl.MouseMotionEvent:
		for _, mouseMover := range events.MouseMovers {
			mouseMover.MouseMove(e.X, e.Y, e.Xrel, e.Yrel)
		}
	}
}
