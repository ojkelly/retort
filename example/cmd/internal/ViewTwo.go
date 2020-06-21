package internal

import (
	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/components/text"
	"retort.dev/example/components"
	"retort.dev/r"
)

func ViewTwo(p r.Properties) r.Element {
	return r.CreateElement(
		box.Box,
		r.Properties{
			box.Properties{
				Direction: box.DirectionRow,
				Border: box.Border{
					Style:      box.BorderStyleSingle,
					Foreground: tcell.ColorWhite,
				},
				Title: box.Label{
					Value: "View Two",
				},
			},
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
						// Padding: box.Padding{
						// 	Top:    0,
						// 	Right:  1,
						// 	Bottom: 1,
						// 	Left:   1,
						// },
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
								Value: loremIpsum,
								// WordBreak:  text.BreakAll,
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
