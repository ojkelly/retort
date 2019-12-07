package main

import (
	"github.com/gdamore/tcell"

	"retort.dev/component/box"
	"retort.dev/component/text"
	"retort.dev/example"
	"retort.dev/r"
)

func main() {

	loremIpsum := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut nec
metus id tellus iaculis porttitor. Maecenas malesuada vitae ex et pharetra.
Suspendisse potenti. In hac habitasse platea dictumst. Mauris malesuada nunc
id placerat euismod. Etiam consectetur nisl in dolor pulvinar bibendum. Cras
molestie ornare placerat. Donec in varius sapien, et mattis augue. Aliquam
viverra nisl at turpis fringilla faucibus. Pellentesque congue viverra
pharetra. Fusce hendrerit bibendum bibendum. Curabitur non tincidunt
nulla. Vestibulum id eros at ex venenatis sollicitudin. Sed in ante quis quam
finibus tristique at sed tortor. Ut maximus molestie ante et elementum. Cras
eget purus eget ante maximus dictum.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut nec
metus id tellus iaculis porttitor. Maecenas malesuada vitae ex et pharetra.
Suspendisse potenti. In hac habitasse platea dictumst. Mauris malesuada nunc
id placerat euismod. Etiam consectetur nisl in dolor pulvinar bibendum. Cras
molestie ornare placerat. Donec in varius sapien, et mattis augue. Aliquam
viverra nisl at turpis fringilla faucibus. Pellentesque congue viverra
pharetra. Fusce hendrerit bibendum bibendum. Curabitur non tincidunt
nulla. Vestibulum id eros at ex venenatis sollicitudin. Sed in ante quis quam
finibus tristique at sed tortor. Ut maximus molestie ante et elementum. Cras
eget purus eget ante maximus dictum.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut nec
metus id tellus iaculis porttitor. Maecenas malesuada vitae ex et pharetra.
Suspendisse potenti. In hac habitasse platea dictumst. Mauris malesuada nunc
id placerat euismod. Etiam consectetur nisl in dolor pulvinar bibendum. Cras
molestie ornare placerat. Donec in varius sapien, et mattis augue. Aliquam
viverra nisl at turpis fringilla faucibus. Pellentesque congue viverra
pharetra. Fusce hendrerit bibendum bibendum. Curabitur non tincidunt
nulla. Vestibulum id eros at ex venenatis sollicitudin. Sed in ante quis quam
finibus tristique at sed tortor. Ut maximus molestie ante et elementum. Cras
eget purus eget ante maximus dictum.`

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
					box.Box,
					r.Properties{
						box.Properties{
							Foreground:    tcell.ColorBeige,
							FlexDirection: box.FlexDirectionColumn,
							Border: box.Border{
								Style:      box.BorderStyleSingle,
								Foreground: tcell.ColorWhite,
							},
							// BUG(ojkelly): padding doesnt work right for text
							// Padding: box.Padding{
							// 	Top:    1,
							// 	Right:  1,
							// 	Bottom: 1,
							// 	Left:   1,
							// },
							Title: box.Label{
								Value: "Example",
							},
						},
					},
					r.Children{
						r.CreateElement(
							text.Text,
							r.Properties{
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
