package r

import (
	"sync"
)

type Context struct {
	defaultState State
}

func (c *Context) Mount(state State) {
	hookFiberLock.Lock()

	if hookFiber == nil {
		panic("UseContext was not called inside a Component.")
	}
	var oldHook *hook

	if hookFiber != nil &&
		hookFiber.alternate != nil &&
		hookFiber.alternate.hooks != nil &&
		len(hookFiber.alternate.hooks) > hookIndex &&
		hookFiber.alternate.hooks[hookIndex] != nil {
		oldHook = hookFiber.alternate.hooks[hookIndex]
	}

	var h *hook
	if oldHook != nil {
		h = oldHook
		// h.state = s ??
	} else {
		h = &hook{
			tag:     hookTagContext,
			mutex:   &sync.Mutex{},
			state:   state,
			context: c,
		}
	}

	if hookFiber != nil {
		hookFiber.hooks = append(hookFiber.hooks, h)
		hookIndex++
	}

	hookFiberLock.Unlock()
}

// CreateContext allows you to create a Context that can be used with UseContext
// It must be called from outside your Component.
func CreateContext(initial State) *Context {
	if hookFiber != nil {
		panic("CreateContext was called inside a Component.")
	}

	context := &Context{
		defaultState: initial,
	}

	return context
}

// UseContext lets you subscribe to changes of Context without nesting.
func UseContext(c *Context) State {
	hookFiberLock.Lock()
	// Walk the tree up via parents, and stop when we find a Provider
	// that matches our Context
	if hookFiber == nil {
		panic("UseContext was not called inside a Component.")
	}

	state := findContext(c, hookFiber)

	var oldHook *hook

	if hookFiber != nil &&
		hookFiber.alternate != nil &&
		hookFiber.alternate.hooks != nil &&
		len(hookFiber.alternate.hooks) > hookIndex &&
		hookFiber.alternate.hooks[hookIndex] != nil {
		oldHook = hookFiber.alternate.hooks[hookIndex]
	}

	var h *hook
	if oldHook != nil {
		h = oldHook
		h.state = state
	} else {
		h = &hook{
			tag:   hookTagState,
			mutex: &sync.Mutex{},
			state: state,
		}
	}

	var actions []Action
	if oldHook != nil {
		actions = oldHook.queue

		for _, action := range actions {
			h.state = action(h.state)
		}
	}

	if hookFiber != nil {
		hookFiber.hooks = append(hookFiber.hooks, h)
		hookIndex++
	}
	// debug.Spew("useContext", h, state)
	hookFiberLock.Unlock()
	return state
}

func findContext(c *Context, f *fiber) State {
	if f == nil {
		return nil
	}

	foundContext := false
	var matchingHook *hook

	for _, h := range f.hooks {
		if h == nil {
			continue
		}
		if h.tag == hookTagContext && h.context == c {
			foundContext = true
			matchingHook = h
		}
	}

	if !foundContext || matchingHook == nil {
		if f.parent == nil {
			return nil
		}
		return findContext(c, f.parent)
	}
	if matchingHook == nil || len(matchingHook.state) != 1 {
		return nil
	}

	return matchingHook.state
}
