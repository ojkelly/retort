package retort_test

import (
	"retort.dev/component"
	"retort.dev/r"
)

type ContainerProps struct {
	Width  int
	Height int
}

func Container(p r.Properties) r.Element {
	props := p.GetProperty(
		ContainerProps{},
		"Container requires ContainerProps",
	).(ContainerProps)
	children := p.GetProperty(
		r.Children{},
		"Container requires r.Children",
	).(r.Children)

	return r.CreateFragment(r.Properties{props}, children)
}

func Example_component() {

	r.CreateElement(
		Container,
		r.Properties{
			ContainerProps{
				Width:  100,
				Height: 100,
			},
		},
		r.Children{
			r.CreateElement(
				component.Box,
				r.Properties{},
				nil,
			),
		},
	)
}
