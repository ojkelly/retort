package components

import (
	"time"

	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/r"
)

type EffectExampleBoxState struct {
	Color tcell.Color
}

func EffectExampleBox(p r.Properties) r.Element {
	boxProps := p.GetProperty(
		box.Properties{},
		"EffectExampleBox requires ContainerProps",
	).(box.Properties)

	children := p.GetProperty(
		r.Children{},
		"EffectExampleBox requires r.Children",
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
		box.Box,
		r.Properties{
			boxProps,
			mouseEventHandler,
		},
		children,
	)
}
