package story

import (
	"github.com/gdamore/tcell"
	"retort.dev/component/box"
	"retort.dev/r"
)

func Story(p r.Properties) r.Element {

	return r.CreateElement(
		box.Box,
		r.Properties{
			box.Properties{

				Border: box.Border{
					Style:      box.BorderStyleSingle,
					Foreground: tcell.ColorOrange,
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
