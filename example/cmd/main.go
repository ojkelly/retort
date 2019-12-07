package main

import (
	_ "net/http/pprof"

	"github.com/gdamore/tcell"

	"retort.dev/component/box"
	"retort.dev/example"
	"retort.dev/r"
)

func main() {

	group := r.Children{
		r.CreateElement(
			box.Box,
			r.Properties{
				box.Properties{
					FlexGrow:   3,
					Foreground: tcell.ColorBeige,
					Border: box.Border{
						Style:      box.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			nil,
		),
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				box.Properties{
					FlexGrow:   1,
					Foreground: tcell.ColorCadetBlue,
					Border: box.Border{
						Style:      box.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			nil,
		),
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				box.Properties{
					FlexGrow:   1,
					Foreground: tcell.ColorLawnGreen,
					Border: box.Border{
						Style:      box.BorderStyleSingle,
						Foreground: tcell.ColorLavender,
					},
				},
			},
			nil,
		),
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				box.Properties{
					FlexGrow:   1,
					Foreground: tcell.ColorLightCyan,
					Border: box.Border{
						Style:      box.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			nil,
		),
	}

	r.Retort(
		r.CreateElement(
			box.Box,
			r.Properties{
				box.Properties{
					Width:  100,
					Height: 100,
				},
			},
			r.Children{
				r.CreateElement(
					example.EffectExampleBox,
					r.Properties{
						box.Properties{
							FlexGrow:      3,
							Foreground:    tcell.ColorBeige,
							FlexDirection: box.FlexDirectionColumn,
							Border: box.Border{
								Style:      box.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
							Padding: box.Padding{
								Top:    1,
								Right:  1,
								Bottom: 1,
								Left:   1,
							},
							Title: box.Label{
								Value: "Example",
							},
						},
					},
					group,
				),
				r.CreateElement(
					example.ClickableBox,
					r.Properties{
						box.Properties{
							FlexGrow:   1,
							Foreground: tcell.ColorCadetBlue,
							Border: box.Border{
								Style:      box.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					nil,
				),
				r.CreateElement(
					example.ClickableBox,
					r.Properties{
						box.Properties{
							FlexGrow:   1,
							Foreground: tcell.ColorLawnGreen,
							Border: box.Border{
								Style:      box.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					nil,
				),
				r.CreateElement(
					example.EffectExampleBox,
					r.Properties{
						box.Properties{
							FlexGrow:   1,
							Foreground: tcell.ColorLightCyan,
							Border: box.Border{
								Style:      box.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					nil,
				),
			},
			// nil,
		),
		r.RetortConfiguration{},
	)
}
