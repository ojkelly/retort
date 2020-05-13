package story

import (
	"retort.dev/components/box"
	"retort.dev/example/hackernews/components/theme"
	"retort.dev/r"
)

func Story(p r.Properties) r.Element {

	tc := r.UseContext(theme.Context)

	t := tc.GetState(
		theme.Colors{},
	).(theme.Colors)

	return r.CreateElement(
		box.Box,
		r.Properties{
			box.Properties{
				Width:      100,
				Height:     100,
				Foreground: t.Foreground,
				Border: box.Border{
					Style:      box.BorderStyleSingle,
					Foreground: t.Border,
				},

				Title: box.Label{
					Value: "Loading Story",
				},
				Overflow: box.OverflowScrollX,
			},
		},
		nil,
	)
}
