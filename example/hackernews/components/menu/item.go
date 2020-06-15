package menu

import (
	"fmt"

	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/components/text"
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
	}
	// debug.Log("menu item loading ", loading, props.Id)

	// onClick := func(
	// 	isPrimary,
	// 	isSecondary bool,
	// 	buttonMask tcell.ButtonMask,
	// ) r.EventMouseClickRelease {
	// 	if isPrimary {
	// 		props.SetTheme(theme.White)
	// 	}
	// 	return func() {}
	// }

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
		// debug.Log("menu item err", err)
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
			box.Properties{
				Direction: box.DirectionColumn,
				Grow:      1,
				Padding: box.Padding{
					Left:  1,
					Right: 1,
				},
				Border: box.Border{
					Foreground: tcell.ColorGray,
					Style:      box.BorderStyleSingle,
				},
			},
		},
		r.Children{
			r.CreateElement(
				text.Text,
				r.Properties{
					box.Properties{Grow: 1},
					text.Properties{
						Value:      story.Title,
						Foreground: t.Foreground,
					},
				},
				nil,
			),
			r.CreateElement(
				text.Text,
				r.Properties{
					box.Properties{Grow: 1},
					text.Properties{
						Value: fmt.Sprintf(
							"Score: %d Comments: %d",
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
