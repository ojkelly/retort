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
		box.Properties{},
	).(box.Properties)

	boxProps := t

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

	boxProps.Title.Value = "Hacker News"
	return r.CreateElement(
		box.Box,
		r.Properties{
			boxProps,
			onClick,
		},
		nil,
	)
}
