package example

import (
	"time"

	"github.com/gdamore/tcell"
	"retort.dev/component"
	"retort.dev/r"
)

type EffectExampleBoxState struct {
	Color tcell.Color
}

func EffectExampleBox(p r.Properties) r.Element {
	boxProps := p.GetProperty(
		component.BoxProps{},
		"Container requires ContainerProps",
	).(component.BoxProps)

	children := p.GetProperty(
		r.Children{},
		"Container requires r.Children",
	).(r.Children)

	s, setState := r.UseState(r.State{
		EffectExampleBoxState{Color: boxProps.Border.Foreground},
	})
	state := s.GetState(
		EffectExampleBoxState{},
	).(EffectExampleBoxState)

	r.UseEffect(func() r.EffectCancel {
		ticker := time.NewTicker(2 * time.Second)

		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					setState(func(s r.State) r.State {
						ms := s.GetState(
							EffectExampleBoxState{},
						).(EffectExampleBoxState)

						color := tcell.ColorGreen
						if ms.Color == tcell.ColorGreen {
							color = tcell.ColorBlue
						}

						if ms.Color == tcell.ColorBlue {
							color = tcell.ColorGreen
						}

						return r.State{EffectExampleBoxState{
							Color: color,
						},
						}
					})
				}
			}
		}()
		return func() {
			<-done
		}
	}, r.EffectDependencies{})

	// var mouseEventHandler r.MouseEventHandler
	mouseEventHandler := func(e *tcell.EventMouse) {
		color := tcell.ColorGreen
		if state.Color == tcell.ColorGreen {
			color = tcell.ColorBlue
		}

		if state.Color == tcell.ColorBlue {
			color = tcell.ColorGreen
		}

		setState(func(s r.State) r.State {
			return r.State{EffectExampleBoxState{
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
