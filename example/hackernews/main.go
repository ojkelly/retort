package main

import (
	"retort.dev/example/hackernews/components/app"
	"retort.dev/r"
)

// TODO: https://github.com/munrocape/hn
func main() {
	r.Retort(
		r.CreateElement(
			app.App,
			nil,
			nil,
		),
		r.RetortConfiguration{},
	)
}
