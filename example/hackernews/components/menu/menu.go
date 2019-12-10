package menu

import (
	"github.com/gdamore/tcell"
	"retort.dev/component/box"
	"retort.dev/example/hackernews/components/theme"
	"retort.dev/r"
)

type Properties struct {
	SetTheme func(t theme.Color)
}

func Menu(p r.Properties) r.Element {
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
					Value: "Hacker News",
				},
			},
			onClick,
		},
		nil,
	)
}
