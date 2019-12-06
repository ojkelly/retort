package r

import "retort.dev/r/internal/quadtree"

// BoxLayout is used by ScreenElements to determine the exact location to
// calculate/render from.
// It represents the concrete positioning information specific
// to the size of the terminal screen.
//
// You never set this directly, it's calculated via a component like
// Box. Which allows for more expressive objects, with padding, and margin.
//
// It is recalculated when the screen is resized.
//
// This layout information is also used to calculate which elements mouse events
// effect.

// You shouldn't use this except for a call to r.CreateScreenElement
type BoxLayout struct {
	X, Y, Rows, Columns int

	// ZIndex is the layer this Box is printed on.
	// Specifically, it determines the order of painting on the screen, with
	// higher numbers being painted later, and appearing on top.
	// This is also used to direct some events, where the highest zindex is used.
	ZIndex int

	// Order is set to control the display order of a group of flex box children
	Order int
}

// reconcileQuadTree updates the quadtree with our new layout, and provides
// the default box layout (from the parent) if none is available on the element
func (r *retort) reconcileQuadTree(f *fiber) {
	if f == nil {
		return
	}

	skip := false

	boxLayout := f.Properties.GetOptionalProperty(
		BoxLayout{},
	).(BoxLayout)

	if boxLayout.X == 0 &&
		boxLayout.Y == 0 &&
		boxLayout.Rows == 0 &&
		boxLayout.Columns == 0 {
		skip = true
	}

	if !skip {
		r.quadtree.Insert(quadtree.Bounds{
			X:      boxLayout.X,
			Y:      boxLayout.Y,
			Width:  boxLayout.Columns,
			Height: boxLayout.Rows,

			// Store a pointer to our fiber for retrieval
			// We will need to cast this on the way out
			Value: f,
		})
	}

	r.reconcileQuadTree(f.child)
	r.reconcileQuadTree(f.sibling)
}
