package example

import (
	"github.com/gdamore/tcell"
	"retort.dev/component"
	"retort.dev/debug"
	"retort.dev/r"
)

type MovingBoxState struct {
	Color tcell.Color
}

func ClickableBox(p r.Properties) r.Element {
	boxProps := p.GetProperty(
		component.BoxProps{},
		"Container requires ContainerProps",
	).(component.BoxProps)

	children := p.GetProperty(
		r.Children{},
		"Container requires r.Children",
	).(r.Children)

	s, setState := r.UseState(r.State{
		MovingBoxState{Color: boxProps.Border.Foreground},
	})
	state := s.GetState(
		MovingBoxState{},
	).(MovingBoxState)

	mouseEventHandler := func(e *tcell.EventMouse) {
		debug.Spew("mouseEventHandler", e, state)
		color := tcell.ColorGreen
		if state.Color == tcell.ColorGreen {
			color = tcell.ColorBlue
		}

		if state.Color == tcell.ColorBlue {
			color = tcell.ColorGreen
		}

		setState(func(s r.State) r.State {
			debug.Spew("mouseEventHandler update state", e, state)

			return r.State{MovingBoxState{
				Color: color,
			},
			}
		})
	}

	boxProps.Border.Foreground = state.Color

	return r.CreateElement(
		component.Box,
		r.Properties{
			boxProps,
			mouseEventHandler,
		},
		children,
	)
}
