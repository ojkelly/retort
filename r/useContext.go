package r

type Context struct {
	state State
}

// CreateContext allows you to create a Context that can be used with UseContext
// It must be called from outside your Component.
func CreateContext(initial State) Context {
	if hookFiber != nil {
		panic("CreateContext was called inside a Component.")
	}

	context := Context{
		state: initial,
	}
	return context
}

// UseContext lets you subscribe to changes of Context without nesting.
func UseContext(c Context) Properties {
	hookFiberLock.Lock()
	// Walk the tree up via parents, and stop when we find a Provider
	// that matches our Context

	// var oldHook *hook

	// if hookFiber != nil &&
	// 	hookFiber.alternate != nil &&
	// 	hookFiber.alternate.hooks != nil &&
	// 	len(hookFiber.alternate.hooks) > hookIndex &&
	// 	hookFiber.alternate.hooks[hookIndex] != nil {
	// 	oldHook = hookFiber.alternate.hooks[hookIndex]
	// }

	// var h *hook
	// if oldHook != nil {
	// 	h = oldHook
	// } else {
	// 	h = &hook{
	// 		tag:   hookTagContext,
	// 		mutex: &sync.Mutex{},
	// 		state: initial,
	// 	}
	// }
	// context := Context{
	// 	f: hookFiber,
	// }
	hookFiberLock.Unlock()
	return Properties{}
}
