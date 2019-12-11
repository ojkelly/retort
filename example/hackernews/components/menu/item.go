package menu

import (
	"fmt"

	"retort.dev/component/box"
	"retort.dev/component/text"
	"retort.dev/debug"
	"retort.dev/example/hackernews/components/common/hooks/hn"
	"retort.dev/example/hackernews/components/theme"
	"retort.dev/r"
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
		MinHeight:     10,
		FlexDirection: box.FlexDirectionColumn,
		Border: box.Border{
			Style:      box.BorderStyleSingle,
			Foreground: t.Border,
		},
	}

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
		debug.Log("menut item err", err)
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
					box.Properties{
						FlexGrow: 1,
					},
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
					box.Properties{
						FlexGrow: 1,
					},
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
