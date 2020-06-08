package r

import (
	"github.com/gdamore/tcell"
	"retort.dev/debug"
	"retort.dev/r/internal/quadtree"
)

type EdgeSizes struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// BlockLayout is used by ScreenElements to determine the exact location to
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
//
// You shouldn't use this except for a call to r.CreateScreenElement
type BlockLayout struct {
	X, Y, Rows, Columns int

	Padding, Border, Margin EdgeSizes

	// Grow, like flex-grow
	// TODO: better docs
	Grow int

	// ZIndex is the layer this Box is printed on.
	// Specifically, it determines the order of painting on the screen, with
	// higher numbers being painted later, and appearing on top.
	// This is also used to direct some events, where the highest zindex is used.
	ZIndex int

	// Order is set to control the display order of a group of children
	Order int

	// Fixed if the Rows, Columns are an explicit fixed size, else they're fluid
	FixedRows, FixedColumns bool

	// Valid is set when the BlockLayout has been initialised somewhere
	// if it's false, it means we've got a default
	Valid bool
}

type BlockLayouts = []BlockLayout

type CalculateLayoutStage int

const (
	// Initial Pass
	// Calculate implicit or explicit absolute bounds
	CalculateLayoutStageInitial CalculateLayoutStage = iota

	// After this Blocks children have calculated their layouts
	// we recalculate this blocks layou
	CalculateLayoutStageWithChildren

	// Final Pass
	CalculateLayoutStageFinal
)

// CalculateLayout
//
// childrenBlockLayouts will be empty until at
// least CalculateLayoutStageWithChildren
//
// innerBlockLayout is the draw area for children blocks, and will
// be smaller due to padding or border effects
type CalculateLayout func(
	s tcell.Screen,
	stage CalculateLayoutStage,
	parentBlockLayout BlockLayout,
	children BlockLayouts,
) (
	blockLayout BlockLayout,
	innerBlockLayout BlockLayout,
	childrenBlockLayouts BlockLayouts,
)

// reconcileQuadTree updates the quadtree with our new layout, and provides
// the default box layout (from the parent) if none is available on the element
func (r *retort) reconcileQuadTree(f *fiber) {
	if f == nil {
		return
	}

	skip := false

	BlockLayout := f.Properties.GetOptionalProperty(
		BlockLayout{},
	).(BlockLayout)

	if BlockLayout.X == 0 &&
		BlockLayout.Y == 0 &&
		BlockLayout.Rows == 0 &&
		BlockLayout.Columns == 0 {
		skip = true
	}

	if !skip {
		r.quadtree.Insert(quadtree.Bounds{
			X:      BlockLayout.X,
			Y:      BlockLayout.Y,
			Width:  BlockLayout.Columns,
			Height: BlockLayout.Rows,

			// Store a pointer to our fiber for retrieval
			// We will need to cast this on the way out
			Value: f,
		})
	}

	r.reconcileQuadTree(f.child)
	r.reconcileQuadTree(f.sibling)
}

func (r *retort) calculateLayout(f *fiber, inheritedBlockLayout BlockLayout) {
	if f == nil {
		return
	}

	screen := UseScreen()

	// default own BlockLayout
	// parentBlockLayout := findParentLayout(f)

	parentBlockLayout := inheritedBlockLayout

	// ?
	f.BlockLayout = parentBlockLayout
	f.InnerBlockLayout = parentBlockLayout

	if f.calculateLayout != nil {
		// if parentFiber := f.parent; parentFiber != nil {
		// 	parentBlockLayout = parentFiber.InnerBlockLayout
		// 	debug.Log("parent", f.parent.BlockLayout)
		// 	debug.Log("parent ib", f.parent.InnerBlockLayout)

		// }
		debug.Log("parentBlockLayout", parentBlockLayout)

		calcLayout := *f.calculateLayout

		blockLayout, innerBlockLayout, _ := calcLayout(
			screen,
			CalculateLayoutStageInitial,
			parentBlockLayout,
			nil,
		)
		debug.Log("clac", blockLayout, innerBlockLayout)
		f.BlockLayout = blockLayout
		f.InnerBlockLayout = innerBlockLayout

		parentBlockLayout = innerBlockLayout
	}

	// r.calculateLayout(f.child)
	// r.calculateLayout(f.sibling)

	r.calculateLayout(f.child, parentBlockLayout)
	r.calculateLayout(f.sibling, parentBlockLayout)

	if f.calculateLayout != nil {

		children := f.Properties.GetOptionalProperty(
			Children{},
		).(Children)
		debug.Log("children", children)

		childrenBlockLayouts := []BlockLayout{}

		for _, c := range children {
			cbl := BlockLayout{}
			if c != nil {
				cbl = c.BlockLayout
			}
			childrenBlockLayouts = append(childrenBlockLayouts, cbl)
		}

		calcLayout := *f.calculateLayout

		_, _, childrenBlockLayouts = calcLayout(
			screen,
			CalculateLayoutStageWithChildren,
			f.InnerBlockLayout,
			childrenBlockLayouts,
		)
		// Put the updated blockLayouts back onto the children
		for i, c := range children {
			if c == nil {
				continue
			}

			c.BlockLayout = childrenBlockLayouts[i]
			c.InnerBlockLayout = childrenBlockLayouts[i]
		}
		// debug.Spew("CDCCCCC", children)

		// debug.Spew("BEFORE", f.Properties)
		f.Properties = ReplaceProps(f.Properties, f.BlockLayout)
		f.Properties = ReplaceProps(f.Properties, children)
		// debug.Spew("AFTER", f.Properties)

		// debug.Spew("end layout", f)
	}

	// r.calculateLayout(f.child, parentBlockLayout)
	// r.calculateLayout(f.sibling, parentBlockLayout)

	// debug.Spew("end of layout", f)
}

// func findParentLayout(f *fiber) (parentBlockLayout BlockLayout) {
// 	if f == nil {
// 		// TODO: this shouldn't happen, it means we've lost the root block
// 		return
// 	}

// 	// Default to own layout
// 	parentBlockLayout = f.BlockLayout

// 	if parentFiber := f.parent; parentFiber != nil {
// 		// If all four of these values are 0, we cannot render to it, so we
// 		// will search higher to find a more valid block
// 		if parentFiber.InnerBlockLayout.X != 0 &&
// 			parentFiber.InnerBlockLayout.Y != 0 &&
// 			parentFiber.InnerBlockLayout.Rows != 0 &&
// 			parentFiber.InnerBlockLayout.Columns != 0 {
// 			return parentFiber.InnerBlockLayout
// 		}
// 	}

// 	if f.parent != nil {
// 		return findParentLayout(f.parent)
// 	}

// 	return
// }
