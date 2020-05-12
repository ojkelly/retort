package app

import (
	"retort.dev/components/box"
	"retort.dev/r"

	"retort.dev/example/hackernews/components/cache"
	"retort.dev/example/hackernews/components/menu"
	"retort.dev/example/hackernews/components/story"
	"retort.dev/example/hackernews/components/theme"
)

type State struct {
	Color theme.Color
}

var defaultState = State{Color: theme.Orange}

func App(p r.Properties) r.Element {

	s, setState := r.UseState(r.State{defaultState})

	state := s.GetState(
		defaultState,
	).(State)

	setTheme := func(t theme.Color) {
		setState(func(s r.State) r.State {
			return r.State{
				State{
					Color: t,
				},
			}
		})
	}

	return r.CreateElement(
		theme.Theme,
		r.Properties{
			theme.Properties{
				Color: state.Color,
			},
		},
		r.Children{
			r.CreateElement(
				cache.Cache,
				r.Properties{},
				r.Children{
					// Wrapper
					r.CreateElement(
						box.Box,
						r.Properties{
							box.Properties{
								Width:  100,
								Height: 100,
							},
						},
						r.Children{
							// Menu
							r.CreateElement(
								menu.Menu,
								r.Properties{
									box.Properties{
										Width: 100,
									},
									menu.Properties{
										SetTheme: setTheme,
									},
								},
								nil,
							),
							// Story view
							r.CreateElement(
								story.Story,
								r.Properties{
									box.Properties{},
								},
								nil,
							),
						},
					),
				},
			),
		},
	)
}
