package text

import (
	"time"

	"github.com/gdamore/tcell"
	"retort.dev/component/box"
	"retort.dev/debug"
	"retort.dev/r"
)

type boxState struct {
	OffsetX, OffsetY int
	lastUpdated      time.Time
}

// Text is the basic building block for a retort app.
// Text implements the Text Model, see Properties
func Text(p r.Properties) r.Element {
	screen := r.UseScreen()

	// Get our Properties
	textProps := p.GetProperty(
		Properties{},
		"Text requires Properties",
	).(Properties)

	// Get our Properties
	boxProps := p.GetOptionalProperty(
		box.Properties{},
	).(box.Properties)

	// Get our BoxLayout
	parentBoxLayout := p.GetProperty(
		r.BoxLayout{},
		"Text requires a parent BoxLayout.",
	).(r.BoxLayout)

	// Get any children
	children := p.GetOptionalProperty(
		r.Children{},
	).(r.Children)
	if len(children) != 0 {
		panic("Text cannot have children")
	}

	s, setState := r.UseState(r.State{
		boxState{lastUpdated: time.Now()},
	})
	state := s.GetState(
		boxState{},
	).(boxState)

	// Calculate the BoxLayout of this Text
	boxLayout := calculateBoxLayout(
		screen,
		parentBoxLayout,
		textProps,
	)

	mouseEventHandler := func(up, down, left, right bool) {
		offsetXDelta := 0
		offsetYDelta := 0

		switch {
		case up:
			offsetXDelta = -1
		case down:
			offsetXDelta = 1
		case left:
			offsetYDelta = -1
		case right:
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

			if boxProps.Overflow == box.OverflowScroll ||
				boxProps.Overflow == box.OverflowScrollX {
				offsetX = min(
					intAbs(state.OffsetX+offsetXDelta),
					int(float64(parentBoxLayout.Columns)/0.2),
				)
			}

			if boxProps.Overflow == box.OverflowScroll ||
				boxProps.Overflow == box.OverflowScrollY {
				offsetY = min(
					intAbs(state.OffsetY+offsetYDelta),
					int(float64(parentBoxLayout.Rows)/0.2),
				)
			}

			return r.State{boxState{
				OffsetX:     offsetX,
				OffsetY:     offsetY,
				lastUpdated: time.Now(),
			},
			}
		})
	}

	props := r.Properties{}

	if boxProps.Overflow != box.OverflowNone {
		props = append(props, mouseEventHandler)
	}

	return r.CreateElement(
		box.Box,
		r.Properties{
			boxProps,
			mouseEventHandler,
		},
		r.Children{
			r.CreateScreenElement(
				func(s tcell.Screen) r.BoxLayout {
					if s == nil {
						panic("Text can't render no screen")
					}

					w, h := s.Size()

					if w == 0 || h == 0 {
						panic("Text can't render on a zero size screen")
					}

					renderText(
						s,
						textProps,
						boxLayout,
						state.OffsetX, state.OffsetY,
					)

					return boxLayout
				},
				r.Properties{},
				nil,
			),
		},
	)
}
