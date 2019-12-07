package box

import (
	"github.com/gdamore/tcell"

	"retort.dev/debug"
	"retort.dev/r"
)

type boxState struct {
	OffsetX, OffsetY int
}

// Box is the basic building block for a retort app.
// Box implements the Box Model, see Properties
func Box(p r.Properties) r.Element {
	screen := r.UseScreen()

	// Get our Properties
	boxProps := p.GetProperty(
		Properties{},
		"Box requires Properties",
	).(Properties)

	// Get our BoxLayout
	parentBoxLayout := p.GetProperty(
		r.BoxLayout{},
		"Box requires a parent BoxLayout.",
	).(r.BoxLayout)

	// Get any children
	children := p.GetOptionalProperty(
		r.Children{},
	).(r.Children)

	s, setState := r.UseState(r.State{
		boxState{},
	})
	state := s.GetState(
		boxState{},
	).(boxState)

	mouseEventHandler := func(ev *tcell.EventMouse) {
		offsetX := 0
		offsetY := 0

		switch ev.Buttons() {
		case tcell.WheelUp:
			offsetX = -1
		case tcell.WheelDown:
			offsetX = 1
		case tcell.WheelLeft:
			offsetY = -1
		case tcell.WheelRight:
			offsetY = 1
		}

		if offsetX == 0 && offsetY == 0 {
			// nothing to update
			return
		}

		setState(func(s r.State) r.State {
			state := s.GetState(
				boxState{},
			).(boxState)

			return r.State{boxState{
				OffsetX: state.OffsetX + offsetX,
				OffsetY: state.OffsetY + offsetY,
			},
			}
		})
	}

	// Calculate the BoxLayout of this Box
	boxLayout, innerBoxLayout := calculateBoxLayout(
		screen,
		parentBoxLayout,
		boxProps,
	)

	innerBoxLayout.OffsetX = state.OffsetX
	innerBoxLayout.OffsetY = state.OffsetY

	// Calculate the BoxLayout of any children
	childrenWithLayout := calculateBoxLayoutForChildren(
		screen,
		boxProps,
		innerBoxLayout,
		children,
	)
	debug.Spew("innerBoxLayout.OffsetX", innerBoxLayout.OffsetX)
	debug.Spew("state.OffsetX", state.OffsetX)
	return r.CreateScreenElement(
		func(s tcell.Screen) r.BoxLayout {
			if s == nil {
				panic("Box can't render no screen")
			}

			w, h := s.Size()

			if w == 0 || h == 0 {
				panic("Box can't render on a zero size screen")
			}

			render(
				screen,
				boxProps,
				boxLayout,
			)

			return boxLayout
		},
		r.Properties{mouseEventHandler},
		childrenWithLayout,
	)
}
