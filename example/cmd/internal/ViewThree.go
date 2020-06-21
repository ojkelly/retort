package internal

import (
	"fmt"

	"github.com/gdamore/tcell"
	"retort.dev/components/box"
	"retort.dev/components/text"
	"retort.dev/example/components"
	"retort.dev/r"
)

func ViewThree(p r.Properties) r.Element {
	itemState, _ := UseItems(50) // setSelectedItem

	menuItems := r.Children{}

	for i, item := range itemState.Items {
		item := r.CreateElement(
			text.Text,
			r.Properties{
				text.Properties{
					Value:      fmt.Sprintf("[%d] %s", i, item.Title),
					Foreground: tcell.ColorWhite,
				},
			},
			nil,
		)

		menuItems = append(menuItems, item)
	}

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
					Value: "View Three [50 Items]",
				},
			},
		},
		r.Children{
			r.CreateElement(
				components.ClickableBox,
				r.Properties{
					box.Properties{
						Direction:  box.DirectionColumn,
						Grow:       1,
						Foreground: tcell.ColorCadetBlue,
						Border: box.Border{
							Style:      box.BorderStyleSingle,
							Foreground: tcell.ColorWhite,
						},
						Title: box.Label{
							Value: "Items",
						},
					},
				},
				menuItems,
			),
			r.CreateElement(
				components.ClickableBox,
				r.Properties{
					box.Properties{
						Grow:       2,
						Foreground: tcell.ColorLawnGreen,
						Border: box.Border{
							Style:      box.BorderStyleSingle,
							Foreground: tcell.ColorWhite,
						},
						Title: box.Label{
							Value: "Details",
						},
					},
				},
				nil,
			),
		},
	)
}
