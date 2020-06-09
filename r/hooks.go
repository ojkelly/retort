package r

import (
	"sync"
)

// hook is a struct containing the fields needed for all core hooks.
// The hookTag determines which fields are in use.
type hook struct {
	tag   hookTag
	mutex *sync.Mutex

	// UseState
	state State
	queue []Action

	// UseEffect
	deps   EffectDependencies
	effect Effect
	cancel EffectCancel

	// UseContext
	context *Context
}

type hookTag int

const (
	hookTagNothing hookTag = iota
	hookTagState
	hookTagEffect
	hookTagReducer
	hookTagContext
)

// Clone safely makes a copy of a hook for use with fiber updates
func (h *hook) Clone() *hook {
	return &hook{
		tag:   h.tag,
		mutex: h.mutex,

		state: h.state,
		queue: h.queue,

		deps:   h.deps,
		effect: h.effect,
		cancel: h.cancel,

		context: h.context,
	}
}
