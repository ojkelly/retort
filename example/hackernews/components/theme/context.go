package theme

import (
	"retort.dev/r"
)

var Context = r.CreateContext(r.State{orange})

type Properties struct {
	Color Color
}

func Theme(p r.Properties) r.Element {
	children := p.GetProperty(
		r.Children{},
		"Theme requires r.Children",
	).(r.Children)

	props := p.GetProperty(
		Properties{
			Color: Orange,
		},
		"Theme requires Properties",
	).(Properties)

	state := Colors{}
	switch props.Color {
	case Orange:
		state = orange
	case White:
		state = white
	}

	Context.SetState(r.State{state})

	return r.CreateFragment(children)
}
