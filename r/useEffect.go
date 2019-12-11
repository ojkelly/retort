package r

import (
	"reflect"

	"retort.dev/debug"
)

type (
	// Effect is the function type you pass to UseEffect.
	//
	// It must return an EffectCancel, even if there is nothing to clean up.
	//
	// In the Effect you can have a routine to do something (such as fetching
	// data), and then call SetState from a UseState hook, to update your
	// Component.
	Effect func() EffectCancel
	// EffectCancel is a function that must be returned by your Effect, and is
	// called when the effect is cleaned up or canceled. This allows you to finish
	// anything you were doing such as closing channels, connections or files.
	EffectCancel func()
	// EffectDependencies lets you pass in what your Effect depends upon.
	// If they change, your Effect will be re-run.
	EffectDependencies []interface{}
)

// UseEffect is a retort hook that can be called in your Component to run side
// effects.
//
// Where UseState allows your components to re-render when their State changes,
// UseEffect allows you to change that state when you need to.
//
// Data fetching is a good example of when you would want something like
// UseEffect.
//
// Example
//
// The example below is a reasonably simple one that changes the color of the
// border of a box ever 2 seconds. The point here is to show how you can run a
// goroutine in the UseEffect callback, and clean up the channel in the
// EffectCancel return function.
//
//  import (
//   "time"
//
//   "github.com/gdamore/tcell"
//   "retort.dev/component"
//   "retort.dev/r"
//  )
//
//  type EffectExampleBoxState struct {
//    Color tcell.Color
//  }
//
//  func EffectExampleBox(p r.Properties) r.Element {
//    boxProps := p.GetProperty(
//      component.BoxProps{},
//      "Container requires ContainerProps",
//    ).(component.BoxProps)
//
//    children := p.GetProperty(
//      r.Children{},
//      "Container requires r.Children",
//    ).(r.Children)
//
//    s, setState := r.UseState(r.State{
//      EffectExampleBoxState{Color: boxProps.Border.Foreground},
//    })
//    state := s.GetState(
//      EffectExampleBoxState{},
//    ).(EffectExampleBoxState)
//
//    r.UseEffect(func() r.EffectCancel {
//      ticker := time.NewTicker(2 * time.Second)
//
//      done := make(chan bool)
//
//      go func() {
//        for {
//          select {
//          case <-done:
//            return
//          case <-ticker.C:
//            setState(func(s r.State) r.State {
//              ms := s.GetState(
//                EffectExampleBoxState{},
//              ).(EffectExampleBoxState)
//
//              color := tcell.ColorGreen
//              if ms.Color == tcell.ColorGreen {
//                color = tcell.ColorBlue
//              }
//
//              if ms.Color == tcell.ColorBlue {
//                color = tcell.ColorGreen
//              }
//
//              return r.State{EffectExampleBoxState{
//                Color: color,
//              },
//              }
//            })
//          }
//        }
//      }()
//      return func() {
//        <-done
//      }
//    }, r.EffectDependencies{})
//
//    // var mouseEventHandler r.MouseEventHandler
//    mouseEventHandler := func(e *tcell.EventMouse) {
//      color := tcell.ColorGreen
//      if state.Color == tcell.ColorGreen {
//        color = tcell.ColorBlue
//      }
//
//      if state.Color == tcell.ColorBlue {
//        color = tcell.ColorGreen
//      }
//
//      setState(func(s r.State) r.State {
//        return r.State{EffectExampleBoxState{
//          Color: color,
//        },
//        }
//      })
//    }
//
//    boxProps.Border.Foreground = state.Color
//
//    return r.CreateElement(
//      component.Box,
//      r.Properties{
//        boxProps,
//        mouseEventHandler,
//      },
//      children,
//    )
//  }
func UseEffect(effect Effect, deps EffectDependencies) {
	hookFiberLock.Lock()

	var oldHook *hook
	if hookFiber != nil &&
		hookFiber.alternate != nil &&
		hookFiber.alternate.hooks != nil &&
		len(hookFiber.alternate.hooks) > hookIndex &&
		hookFiber.alternate.hooks[hookIndex] != nil {
		oldHook = hookFiber.alternate.hooks[hookIndex]
	}

	hasChanged := true

	if oldHook != nil {
		hasChanged = hasDepsChanged(oldHook.deps, deps)
	}

	h := &hook{
		tag:    hookTagEffect,
		effect: nil,
		cancel: nil,
	}

	if hasChanged {
		h.effect = effect
		h.deps = deps
	}

	if hookFiber != nil {
		hookFiber.hooks = append(hookFiber.hooks, h)
		hookIndex++
	}
	hookFiberLock.Unlock()
}

func hasDepsChanged(
	prevDeps,
	nextDeps EffectDependencies,
) (changed bool) {
	if !reflect.DeepEqual(prevDeps, nextDeps) {
		changed = true
	}

	if len(prevDeps) == 0 {
		changed = false
	}
	if len(nextDeps) == 0 {
		changed = false
	}

	// Check the slices have the same contents, in the same order
	// for i, pd := range prevDeps {
	// 	if nextDeps[i] != pd {
	// 		changed = true
	// 	}
	// }

	debug.Log("hasDepsChanged ", changed)
	return
}

func (r *retort) processEffects(f *fiber) {
	runEffects := false
	cancelEffects := false

	if f == nil {
		return
	}

	switch f.effect {
	case fiberEffectNothing:
		runEffects = true
	case fiberEffectPlacement:
		runEffects = true
	case fiberEffectUpdate:
		runEffects = true
	case fiberEffectDelete:
		cancelEffects = true
	}

	if f.hooks != nil {
		if cancelEffects {
			for _, h := range f.hooks {
				if h.tag != hookTagEffect ||
					h.effect == nil {
					continue
				}
				h.cancel()
			}
		}

		if runEffects {
			for _, h := range f.hooks {
				if h.tag != hookTagEffect ||
					h.effect == nil {
					continue
				}
				h.cancel = h.effect()
			}
		}
	}

	r.processEffects(f.child)
	r.processEffects(f.sibling)
}
