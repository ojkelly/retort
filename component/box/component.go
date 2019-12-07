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
			return r.State{boxState{
				OffsetX:     min(intAbs(state.OffsetX+offsetX), int(float64(parentBoxLayout.Columns)/0.2)),
				OffsetY:     min(intAbs(state.OffsetY+offsetY), int(float64(parentBoxLayout.Rows)/0.2)),
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
