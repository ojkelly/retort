package r

import (
	"sync"
)

var hookIndex int

type hook struct {
	tag   hookTag
	mutex *sync.Mutex

	// State
	state State
	queue []Action

	// Effect
	deps   EffectDependencies
	effect Effect
	cancel EffectCancel
}

type hookTag int

const (
	hookTagNothing hookTag = iota
	hookTagState
	hookTagEffect
	hookTagReducer
	hookTagContext
)

func (h *hook) Clone() *hook {
	return &hook{
		tag:   h.tag,
		mutex: h.mutex,

		state: h.state,
		queue: h.queue,

		deps:   h.deps,
		effect: h.effect,
		cancel: h.cancel,
	}
}
