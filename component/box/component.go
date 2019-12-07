package box

import (
	"time"

	"github.com/gdamore/tcell"

	"retort.dev/debug"
	"retort.dev/r"
)

type boxState struct {
	OffsetX, OffsetY int
	lastUpdated      time.Time
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
		boxState{lastUpdated: time.Now()},
	})
	state := s.GetState(
		boxState{},
	).(boxState)

	mouseEventHandler := func(ev *tcell.EventMouse) {
		offsetXDelta := 0
		offsetYDelta := 0

		switch ev.Buttons() {
		case tcell.WheelUp:
			offsetXDelta = -1
		case tcell.WheelDown:
			offsetXDelta = 1
		case tcell.WheelLeft:
			offsetYDelta = -1
		case tcell.WheelRight:
			offsetYDelta = 1
		}

		if offsetXDelta == 0 && offsetYDelta == 0 {
			// nothing to update
			return
		}

		now := time.Now()

		if now.Sub(state.lastUpdated) < 16*time.Millisecond {
			debug.Spew("throttled ", now.Sub(state.lastUpdated), state.lastUpdated, now)
			// throttle to one update a second
			return
		}

		setState(func(s r.State) r.State {
			state := s.GetState(
				boxState{},
			).(boxState)

			// BUG(ojkelly): this is a bit janky and could be better

			offsetX := state.OffsetX
			offsetY := state.OffsetY

			if boxProps.Overflow == OverflowScroll || boxProps.Overflow == OverflowScrollX {
				offsetX = min(intAbs(state.OffsetX+offsetXDelta), int(float64(parentBoxLayout.Columns)/0.2))
			}

			if boxProps.Overflow == OverflowScroll || boxProps.Overflow == OverflowScrollY {
				offsetY = min(intAbs(state.OffsetY+offsetYDelta), int(float64(parentBoxLayout.Rows)/0.2))
			}

			return r.State{boxState{
				OffsetX:     offsetX,
				OffsetY:     offsetY,
				lastUpdated: time.Now(),
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

	props := r.Properties{}

	if boxProps.Overflow != OverflowNone {
		props = append(props, mouseEventHandler)
	}

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
		props,
		childrenWithLayout,
	)
}

func intAbs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
