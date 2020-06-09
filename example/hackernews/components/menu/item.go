package menu

import (
	"fmt"

	"retort.dev/components/box"
	"retort.dev/components/text"
	"retort.dev/example/hackernews/components/common/hooks/hn"
	"retort.dev/example/hackernews/components/theme"
	"retort.dev/r"
	"retort.dev/r/debug"
)

type MenuItemProps struct {
	Id int
}

func MenuItem(p r.Properties) r.Element {
	props := p.GetProperty(
		MenuItemProps{},
		"MenuItem requires MenuItemProps",
	).(MenuItemProps)

	story, loading, err := hn.UseStory(props.Id)

	tc := r.UseContext(theme.Context)

	t := tc.GetState(
		theme.Colors{},
	).(theme.Colors)

	boxProps := box.Properties{
		Margin: box.Margin{
			Bottom: 1,
		},
		Padding: box.Padding{
			Left:  1,
			Right: 1,
		},
		MinHeight: 10,
		Border: box.Border{
			Style:      box.BorderStyleSingle,
			Foreground: t.Border,
		},
	}
	debug.Log("menu item loading", loading)

	if loading {
		return r.CreateElement(
			text.Text,
			r.Properties{
				boxProps,
				text.Properties{
					Value:      "Loading",
					Foreground: t.Subtle,
				},
			},
			nil,
		)
	}

	if err != nil {
		debug.Log("menu item err", err)
		return r.CreateElement(
			text.Text,
			r.Properties{
				boxProps,
				text.Properties{
					Value:      fmt.Sprintf("%s", err),
					Foreground: t.Foreground,
				},
			},
			nil,
		)
	}

	if story == nil {
		return nil
	}
	// return nil
	return r.CreateElement(
		box.Box,
		r.Properties{
			boxProps,
		},
		r.Children{
			r.CreateElement(
				text.Text,
				r.Properties{
					box.Properties{},
					text.Properties{
						Value:      story.Title,
						Foreground: t.Foreground,
						// WordBreak:  text.BreakAll,
					},
				},
				nil,
			),
			r.CreateElement(
				text.Text,
				r.Properties{
					box.Properties{},
					text.Properties{
						Value: fmt.Sprintf(
							"Score: %d\nComments: %d",
							story.Score,
							story.Descendants,
						),
						Foreground: t.Subtle,
					},
				},
				nil,
			),
		},
	)
}
