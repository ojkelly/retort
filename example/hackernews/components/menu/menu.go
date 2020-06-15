package menu

import (
	"fmt"

	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/example/hackernews/components/common/hooks/hn"
	"retort.dev/example/hackernews/components/theme"
	"retort.dev/r"
)

type Properties struct {
	SetTheme func(t theme.Color)
}

func Menu(p r.Properties) r.Element {
	title := "Top Stories"
	props := p.GetProperty(
		Properties{},
		"Menu requires menu.Properties",
	).(Properties)

	tc := r.UseContext(theme.Context)

	t := tc.GetState(
		theme.Colors{},
	).(theme.Colors)

	onClick := func(
		isPrimary,
		isSecondary bool,
		buttonMask tcell.ButtonMask,
	) r.EventMouseClickRelease {
		if isPrimary {
			props.SetTheme(theme.White)
		}
		return func() {}
	}

	stories := hn.UseTopStories()

	// debug.Spew("stories", stories)

	var items r.Children
	if stories.Data != nil &&
		len(stories.Data) > 0 {

		for _, id := range stories.Data {
			items = append(items, r.CreateElement(
				MenuItem,
				r.Properties{
					MenuItemProps{
						Id: id,
					},
				},
				nil,
			))
		}
	}

	if stories.Loading {
		title = fmt.Sprintf("%s %s", title, "[ Loading ]")
	}

	return r.CreateElement(
		box.Box,
		r.Properties{
			box.Properties{
				Foreground: t.Foreground,
				Border: box.Border{
					Style:      box.BorderStyleSingle,
					Foreground: t.Border,
				},
				Title: box.Label{
					Value: title,
				},
				Direction: box.DirectionColumn,
				Footer: box.Label{
					Value: "Hacker News",
					Wrap:  box.LabelWrapSquareBracket,
				},
				Overflow:  box.OverflowScrollX,
				MinHeight: 5,
			},
			onClick,
		},
		items,
	)
}
