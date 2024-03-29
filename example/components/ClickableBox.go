package components

import (
	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/r"
)

type MovingBoxState struct {
	Color tcell.Color
}

func ClickableBox(p r.Properties) r.Element {
	boxProps := p.GetProperty(
		box.Properties{},
		"Container requires ContainerProps",
	).(box.Properties)

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

	mouseEventHandler := func(
		isPrimary,
		isSecondary bool,
		buttonMask tcell.ButtonMask,
	) r.EventMouseClickRelease {
		color := tcell.ColorGreen
		if state.Color == tcell.ColorGreen {
			color = tcell.ColorBlue
		}

		if state.Color == tcell.ColorBlue {
			color = tcell.ColorGreen
		}

		setState(func(s r.State) r.State {
			return r.State{MovingBoxState{
				Color: color,
			},
			}
		})
		return func() {}
	}

	boxProps.Border.Foreground = state.Color

	return r.CreateElement(
		box.Box,
		r.Properties{
			boxProps,
			mouseEventHandler,
		},
		children,
	)
}
