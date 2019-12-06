package main

import (
	_ "net/http/pprof"

	"github.com/gdamore/tcell"

	"retort.dev/component"
	"retort.dev/example"
	"retort.dev/r"
)

func main() {

	group := r.Children{
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				component.BoxProps{
					FlexGrow:   3,
					Foreground: tcell.ColorBeige,
					Border: component.Border{
						Style:      component.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			nil,
		),
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				component.BoxProps{
					FlexGrow:   1,
					Foreground: tcell.ColorCadetBlue,
					Border: component.Border{
						Style:      component.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			nil,
		),
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				component.BoxProps{
					FlexGrow:   1,
					Foreground: tcell.ColorLawnGreen,
					Border: component.Border{
						Style:      component.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			nil,
		),
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				component.BoxProps{
					FlexGrow:   1,
					Foreground: tcell.ColorLightCyan,
					Border: component.Border{
						Style:      component.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			nil,
		),
	}

	r.Retort(
		r.CreateElement(
			example.ClickableBox,
			r.Properties{
				component.BoxProps{
					Width:      100,
					Height:     100,
					Foreground: tcell.ColorGold,
					Border: component.Border{
						Style:      component.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
				},
			},
			r.Children{
				r.CreateElement(
					example.EffectExampleBox,
					r.Properties{
						component.BoxProps{
							FlexGrow:      3,
							Foreground:    tcell.ColorBeige,
							FlexDirection: component.FlexDirectionColumn,
							Border: component.Border{
								Style:      component.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					group,
				),
				r.CreateElement(
					example.ClickableBox,
					r.Properties{
						component.BoxProps{
							FlexGrow:   1,
							Foreground: tcell.ColorCadetBlue,
							Border: component.Border{
								Style:      component.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					nil,
				),
				r.CreateElement(
					example.ClickableBox,
					r.Properties{
						component.BoxProps{
							FlexGrow:   1,
							Foreground: tcell.ColorLawnGreen,
							Border: component.Border{
								Style:      component.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
						},
					},
					nil,
				),
				r.CreateElement(
					example.EffectExampleBox,
					r.Properties{
						component.BoxProps{
							FlexGrow:   1,
							Foreground: tcell.ColorLightCyan,
							Border: component.Border{
								Style:      component.BorderStyleSingle,
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
