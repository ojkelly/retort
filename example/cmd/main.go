package main

import (
	"retort.dev/example/cmd/internal"
	"retort.dev/r"
)

func main() {
	r.Retort(
		r.CreateElement(
			internal.Switcher,
			r.Properties{},
			nil,
		),
		r.RetortConfiguration{},
	)
}
