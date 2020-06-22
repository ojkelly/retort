package internal

import (
	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/components/text"
	"retort.dev/example/components"
	"retort.dev/r"
)

func ViewOne(p r.Properties) r.Element {
	boxProps := p.GetOptionalProperty(
		box.Properties{},
	).(box.Properties)

	boxProps.Direction = box.DirectionColumn
	boxProps.Border = box.Border{
		Style:      box.BorderStyleSingle,
		Foreground: tcell.ColorWhite,
	}

	boxProps.Title = box.Label{
		Value: "View One",
	}

	return r.CreateElement(
		box.Box,
		r.Properties{
			boxProps,
		},
		r.Children{
			r.CreateElement(
				box.Box,
				r.Properties{
					box.Properties{
						Foreground: tcell.ColorBeige,
						Grow:       3,
						Border: box.Border{
							Style:      box.BorderStyleSingle,
							Foreground: tcell.ColorWhite,
						},
						Title: box.Label{
							Value: "Grow 3 - with text",
						},
					},
				},
				r.Children{
					r.CreateElement(
						text.Text,
						r.Properties{
							box.Properties{
								Overflow: box.OverflowScroll,
							},
							text.Properties{
								Value:      loremIpsum,
								Foreground: tcell.ColorWhite,
							},
						},
						nil,
					),
				},
			),
			r.CreateElement(
				components.ClickableBox,
				r.Properties{
					box.Properties{
						Grow:       2,
						Foreground: tcell.ColorCadetBlue,
						Border: box.Border{
							Style:      box.BorderStyleSingle,
							Foreground: tcell.ColorWhite,
						},
						Title: box.Label{
							Value: "Grow 2",
						},
					},
				},
				nil,
			),
			r.CreateElement(
				components.ClickableBox,
				r.Properties{
					box.Properties{
						Grow:       1,
						Foreground: tcell.ColorLawnGreen,
						Border: box.Border{
							Style:      box.BorderStyleSingle,
							Foreground: tcell.ColorWhite,
						},
						Title: box.Label{
							Value: "Grow 1",
						},
					},
				},
				nil,
			),
		},
	)
}
