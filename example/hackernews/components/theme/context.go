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

	color := Colors{}
	switch props.Color {
	case Orange:
		color = orange
	case White:
		color = white
	}

	// TODO: double check this
	s, _ := r.UseState(r.State{color})

	state := s.GetState(
		Colors{},
	).(Colors)

	Context.Mount(r.State{state})

	return r.CreateFragment(children)
}
