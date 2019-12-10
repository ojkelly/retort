package component

import (
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
