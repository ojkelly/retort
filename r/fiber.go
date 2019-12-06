package r

type fiberEffect int

const (
	fiberEffectNothing fiberEffect = iota
	fiberEffectUpdate
	fiberEffectPlacement
	fiberEffectDelete
)

type componentType int

const (
	nothingComponent componentType = iota
	elementComponent
	fragmentComponent
	screenComponent
)

type fiber struct {
	// dirty when there are changes to commit
	dirty bool

	componentType componentType
	component     Component
	Properties    Properties
	parent        *fiber
	sibling       *fiber
	child         *fiber
	alternate     *fiber
	effect        fiberEffect
	hooks         []*hook

	// Layout Information
	renderToScreen *RenderToScreen
	// this boxLayout is used internally, mainly to route events
	// Its different to the boxLayout that may be passed around in props
	boxLayout BoxLayout

	focus bool
}

func (f *fiber) Clone() (newFiber *fiber) {
	// Parent, sibling, and alternate are not cloned
	// as doing so will recurse forever
	newFiber = &fiber{
		componentType: f.componentType,
		component:     f.component,
		Properties:    f.Properties,
		effect:        f.effect,
		alternate:     f.alternate,
		hooks:         f.hooks,
		boxLayout:     f.boxLayout,
	}

	if f.renderToScreen != nil {
		render := *f.renderToScreen
		newFiber.renderToScreen = &render
	}

	if f.child != nil {
		newFiber.child = f.child.Clone()
	}
	if f.sibling != nil {
		newFiber.sibling = f.sibling.Clone()
	}

	return
}

func cloneElements(fibers []*fiber) (cloned []*fiber) {
	cloned = []*fiber{}

	for _, f := range fibers {
		cloned = append(cloned, f.Clone())
	}

	return
}
