package r

import (
	"github.com/gdamore/tcell"
	"retort.dev/debug"
)

// RenderToScreen is the callback passed to create a Screen Element
type RenderToScreen func(
	s tcell.Screen,
) BoxLayout

func (r *retort) commitRoot() {
	screen := UseScreen()

	for _, deletion := range r.deletions {
		r.commitWork(deletion)
	}

	// Draw
	r.commitWork(r.wipRoot)
	// debug.Spew("committed work", r)
	screen.Show()

	// Update effects
	r.processEffects(r.wipRoot)
	// debug.Spew("processed effects")

	// Update our quadtree for collisions
	r.quadtree.Clear()
	// debug.Spew("clear QuadTree", r.quadtree)
	r.reconcileQuadTree(r.wipRoot)
	debug.Spew("reconciled QuadTree", r.quadtree.Total)

	r.currentRoot = r.wipRoot
	r.wipRoot = nil
	r.hasChangesToRender = false
	debug.Spew("committed root")
}

func (r *retort) commitWork(f *fiber) {
	if f == nil {
		return
	}

	screen := UseScreen()

	// TODO: collect all the renderToScreen paired with their zIndex
	// render all from lowest to highest index
	switch f.effect {
	case fiberEffectNothing:
	case fiberEffectPlacement:
		// debug.Spew("fiberEffectPlacement", f)
		// TODO: extract all renderToScreen's and execute them in ZIndex order lowest to highest
		// this should allow layered things
		// for _, el := range f.elements {
		if f.renderToScreen == nil {
			break
		}

		// if el.renderToScreen == nil {
		// 	continue
		// }

		// need to keep track of previous location of this element
		// so when it's called we can clear the screen it used to be in before redrawing
		render := *f.renderToScreen
		f.boxLayout = render(screen)
		// }
	case fiberEffectUpdate:
		// cancelEffects(f)
		// for _, el := range f.elements {
		// if el == nil || el.renderToScreen == nil {
		// 	continue
		// }

		if f.renderToScreen == nil {
			break
		}

		// need to keep track of previous location of this element
		// so when it's called we can clear the screen it used to be in before redrawing
		render := *f.renderToScreen
		f.boxLayout = render(screen)
		// }
	case fiberEffectDelete:
	}

	r.commitWork(f.child)
	r.commitWork(f.sibling)

	f.dirty = false
}

func (r *retort) commitDeletion(f *fiber) {
	// if (fiber.dom) {
	//   domParent.removeChild(fiber.dom);
	// } else {
	//   commitDeletion(fiber.child, domParent);
	// }
}
