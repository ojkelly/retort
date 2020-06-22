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
	ViewThreeSelected
	ViewFourSelected
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

	children := r.Children{
		r.CreateElement(
			ViewOne,
			r.Properties{
				box.Properties{
					Hide: state.Selected != ViewOneSelected,
				},
			},
			nil,
		),
		// r.CreateElement(
		// 	ViewTwo,
		// 	r.Properties{
		// 		box.Properties{
		// 			Hide: state.Selected != ViewTwoSelected,
		// 		},
		// 	},
		// 	nil,
		// ),
		// r.CreateElement(
		// 	ViewThree,
		// 	r.Properties{
		// 		box.Properties{
		// 			Hide: state.Selected != ViewThreeSelected,
		// 		},
		// 	},
		// 	nil,
		// ),
		// r.CreateElement(
		// 	ViewFour,
		// 	r.Properties{
		// 		box.Properties{
		// 			Hide: state.Selected != ViewFourSelected,
		// 		},
		// 	},
		// 	nil,
		// ),
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

		case 51: // 3
			setState(func(s r.State) r.State {
				return r.State{
					SelectedViewState{
						Selected: ViewThreeSelected,
					},
				}
			})
		case 52: // 4
			setState(func(s r.State) r.State {
				return r.State{
					SelectedViewState{
						Selected: ViewFourSelected,
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
