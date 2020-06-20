package main

import (
	"github.com/gdamore/tcell"

	"retort.dev/components/box"
	"retort.dev/example/components"
	"retort.dev/r"
)

const loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Enim praesent elementum facilisis leo. Et odio pellentesque diam volutpat commodo sed egestas. Risus viverra adipiscing at in tellus. Ornare suspendisse sed nisi lacus sed. Malesuada nunc vel risus commodo viverra maecenas accumsan lacus vel. Sit amet facilisis magna etiam. Bibendum neque egestas congue quisque egestas. Praesent tristique magna sit amet purus. Auctor eu augue ut lectus arcu bibendum at. Urna cursus eget nunc scelerisque viverra mauris in aliquam. Elit at imperdiet dui accumsan sit amet nulla. Sed euismod nisi porta lorem mollis aliquam ut porttitor. Volutpat diam ut venenatis tellus in metus vulputate eu scelerisque. Pharetra pharetra massa massa ultricies mi quis. Porta non pulvinar neque laoreet suspendisse interdum consectetur. Suspendisse in est ante in nibh mauris cursus mattis. Velit ut tortor pretium viverra suspendisse. Interdum varius sit amet mattis vulputate enim nulla.

Venenatis urna cursus eget nunc scelerisque viverra. Libero enim sed faucibus turpis in eu mi bibendum neque. Mi in nulla posuere sollicitudin aliquam ultrices sagittis. Sagittis purus sit amet volutpat. Maecenas ultricies mi eget mauris pharetra et. Ac tortor vitae purus faucibus ornare. Sollicitudin ac orci phasellus egestas tellus rutrum tellus pellentesque. Imperdiet dui accumsan sit amet nulla. Semper feugiat nibh sed pulvinar proin gravida hendrerit lectus. Cras semper auctor neque vitae tempus quam pellentesque nec nam. Cursus sit amet dictum sit amet justo. Aenean vel elit scelerisque mauris pellentesque pulvinar pellentesque habitant. Lacinia at quis risus sed vulputate odio ut.

Sed turpis tincidunt id aliquet. In aliquam sem fringilla ut morbi tincidunt. Pharetra convallis posuere morbi leo urna. Velit euismod in pellentesque massa. Pellentesque massa placerat duis ultricies lacus sed turpis tincidunt id. Risus quis varius quam quisque id diam. Urna condimentum mattis pellentesque id. Id interdum velit laoreet id donec ultrices tincidunt. At auctor urna nunc id cursus metus. Adipiscing diam donec adipiscing tristique risus nec. Ut porttitor leo a diam sollicitudin tempor. Est sit amet facilisis magna etiam. Tellus mauris a diam maecenas sed enim ut.

Et malesuada fames ac turpis egestas. Egestas dui id ornare arcu odio ut sem nulla. Gravida cum sociis natoque penatibus et magnis dis. Tellus in hac habitasse platea. Ultrices tincidunt arcu non sodales. Lorem ipsum dolor sit amet consectetur. Egestas tellus rutrum tellus pellentesque. Ac auctor augue mauris augue neque gravida in fermentum et. Iaculis at erat pellentesque adipiscing commodo. Malesuada fames ac turpis egestas integer eget aliquet nibh praesent. Sit amet consectetur adipiscing elit ut aliquam purus sit amet. Vitae tortor condimentum lacinia quis vel eros donec ac. Purus faucibus ornare suspendisse sed nisi. Mi ipsum faucibus vitae aliquet nec ullamcorper sit amet. Ac turpis egestas sed tempus urna. Nibh venenatis cras sed felis eget velit. Sit amet purus gravida quis blandit turpis cursus in.

Amet est placerat in egestas erat imperdiet sed euismod. Eget felis eget nunc lobortis. Ac auctor augue mauris augue neque. Ac tortor vitae purus faucibus ornare suspendisse. Placerat duis ultricies lacus sed. Tortor vitae purus faucibus ornare suspendisse sed nisi. Vulputate dignissim suspendisse in est ante in nibh. Elit duis tristique sollicitudin nibh sit. Tellus at urna condimentum mattis pellentesque id nibh tortor. Proin fermentum leo vel orci porta non pulvinar neque. Eu ultrices vitae auctor eu augue ut. Erat pellentesque adipiscing commodo elit at imperdiet. Auctor elit sed vulputate mi sit amet mauris. Tellus orci ac auctor augue mauris.`

func main() {
	r.Retort(
		r.CreateElement(
			box.Box,
			r.Properties{
				box.Properties{
					Direction: box.DirectionColumn,
					Border: box.Border{
						Style:      box.BorderStyleSingle,
						Foreground: tcell.ColorWhite,
					},
					Title: box.Label{
						Value: "Wrapper",
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
							// 	Right:  0,
							// 	Bottom: 0,
							// 	Left:   0,
							// },
							Title: box.Label{
								Value: "Grow 3 - with text",
							},
							Overflow: box.OverflowScrollX,
						},
					},
					// r.Children{
					// 	r.CreateElement(
					// 		text.Text,
					// 		r.Properties{
					// 			text.Properties{
					// 				Value:      loremIpsum,
					// 				WordBreak:  text.BreakAll,
					// 				Foreground: tcell.ColorWhite,
					// 			},
					// 		},
					// 		nil,
					// 	),
					// },
					nil,
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
		),
		r.RetortConfiguration{},
	)
}
