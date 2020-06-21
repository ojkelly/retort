package internal

import (
	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/r"
)

type SelectedView int

const (
	ViewOneSelected SelectedView = iota
	ViewTwoSelected
)

type SelectedViewState struct {
	Selected SelectedView
}

func Switcher(p r.Properties) r.Element {
	s, setState := r.UseState(r.State{
		SelectedViewState{Selected: ViewOneSelected},
	})
	state := s.GetState(
		SelectedViewState{},
	).(SelectedViewState)

	children := r.Children{}

	switch state.Selected {
	case ViewOneSelected:
		children = r.Children{r.CreateElement(
			ViewOne,
			r.Properties{},
			nil,
		)}
	case ViewTwoSelected:
		children = r.Children{r.CreateElement(
			ViewTwo,
			r.Properties{},
			nil,
		)}
	}

	keyEventHandler := func(e *tcell.EventKey, meta r.EventMeta) r.EventMeta {
		switch e.Rune() {
		case 49: // 1
			setState(func(s r.State) r.State {
				return r.State{
					SelectedViewState{
						Selected: ViewOneSelected,
					},
				}
			})

		case 50: // 2
			setState(func(s r.State) r.State {
				return r.State{
					SelectedViewState{
						Selected: ViewTwoSelected,
					},
				}
			})
		}

		meta.StopPropegation = true
		return meta
	}

	return r.CreateElement(
		box.Box,
		r.Properties{
			keyEventHandler,
			box.Properties{
				Direction: box.DirectionColumn,
				Border: box.Border{
					Style:      box.BorderStyleSingle,
					Foreground: tcell.ColorGray,
				},
				Title: box.Label{
					Value: "Switcher",
				},
			},
		},
		children,
	)
}
