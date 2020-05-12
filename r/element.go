package r

import (
	"reflect"

	"retort.dev/debug"
)

type (
	// Element is the smallest building block of retort.
	// You create them with r.CreateElement or r.CreateFragment
	//
	// Internally they are a pointer to a fiber, which is used
	// to keep track of the render tree.
	Element *fiber

	// Component is main thing you will be making to create a
	// retort app. Your component must match this function
	// signature.
	Component func(props Properties) Element

	// Children is a slice of Elements (or pointers to a fiber)
	// It's used in r.CreateElement or r.CreateFragment to
	// specify the child nodes of the Element.
	// It can also be extracted from props with GetProperty or
	// GetOptionalProperty, if you want to pass children on.
	//
	// In general, unless you're creating a ScreenElement,
	// you should pass any children passed into props
	// on to the return Element.
	Children []*fiber

	// Properties are immutable state that is passed into a component, and pass
	// down to components to share data.
	//
	// Properties is ultimately a slice of interfaces, which lets you and retort
	// and any components your using add any concrete structs to it. Because of
	// this, there are some helper methods to retrieve props. These are
	// GetProperty and GetOptionalProperty.
	//
	// Properties can only contain one struct of a given type. In this sense the
	// type of the struct is a key.
	//
	// Sometimes called props.
	Properties []interface{}

	// State is local to a component.
	// It is mutable via the setState function from UseState. Don't edit State
	// directly, as retort will not know that you have, and will not trigger an
	// update and re-render.
	// It can be used to create new props to pass down to other components.
	State []interface{}
)

// CreateElement is used to create the building blocks of a retort application,
// and the thing that Components are ultimately made up of, Elements.
//
//  import (
//    "github.com/gdamore/tcell"
//    "retort.dev/r"
//    "retort.dev/r/component"
//  )
//
//  // Wrapper is a simple component that wraps the
//  // children Components in a box with a white border.
//  func Wrapper(p r.Properties) r.Element {
//    children := r.GetProperty(
//      p,
//      r.Children{},
//      "Container requires r.Children",
//    ).(r.Children)
//
//     return r.CreateElement(
//      component.Box,
//      r.Properties{
//          component.BoxProps{
//            Border: component.Border{
//              Style:      component.BorderStyleSingle,
//              Foreground: tcell.ColorWhite,
//            },
//          },
//        },
//      children,
//    )
//  }
//
// By creating an Element and passing Properties and Children seperately, retort
// can keep track of the entire tree of Components, and decide when to compute
// which parts, and in turn when to render those to the screen.
func CreateElement(
	component Component,
	props Properties,
	children Children,
) *fiber {
	// debug.Log("CreateElement", component, props, children)
	if !checkPropTypesAreUnique(props) {
		panic("props are not unique")
	}
	return &fiber{
		componentType: elementComponent,
		component:     component,
		Properties: append(
			props,
			children,
		),
	}
}

// CreateFragment is like CreateElement except you do not need a Component
// or Properties. This is useful when you need to make Higher Order Components,
// or other Components that wrap or compose yet more Components.
func CreateFragment(children Children) *fiber {
	return &fiber{
		componentType: fragmentComponent,
		component:     nil,
		Properties: Properties{
			children,
		},
	}
}

// CreateScreenElement is like a normal Element except it has the
// ability to render output to the screen.
//
// Once retort has finished calculating which components have changed
// all those with changes are passed to a render function.
// This walks the tree and finds ScreenElements and calls their
// RenderToScreen function, passing in the current Screen.
//
// RenderToScreen needs to return a BlockLayout, which is used among
// other things to direct Mouse Events to the right Component.
//
//	func Box(p r.Properties) r.Element {
//		return r.CreateScreenElement(
//			func(s tcell.Screen) r.BlockLayout {
//				return BlockLayout
//			},
//			nil,
//		)
//	}
func CreateScreenElement(
	calculateLayout CalculateLayout,
	render RenderToScreen,
	props Properties,
	children Children,
) *fiber {
	// debug.Log("CreateScreenElement", render)

	// TODO: maybe this should be multi-step
	// ie pass in some functions
	// - calculateLayout BlockLayout
	// 		- this is called more than once
	//		-	if the inputs are the same, or the result is the same as before,
	//		-	then it should be done
	// - renderToScreen(BlockLayout)
	return &fiber{
		componentType:   screenComponent,
		calculateLayout: &calculateLayout,
		renderToScreen:  &render,
		Properties: append(
			props,
			children,
		),
		component: nil,
	}
}

func checkPropTypesAreUnique(props Properties) bool {
	seenPropTypes := make(map[reflect.Type]bool)

	for _, p := range props {
		if seen := seenPropTypes[reflect.TypeOf(p)]; seen {
			return false
		}
		seenPropTypes[reflect.TypeOf(p)] = true
	}
	return true
}

// GetProperty will search props for the Property matching the type of the
// struct you passed in, and will throw an Error with the message provided
// if not found.
//
// This is useful when your component will not work without the provided
// Property. However it is very unforgiving, and generally you will want to use
// GetOptionalProperty which allows you to provide a default Property to use.
//
// Because this uses reflection, you must pass in a concrete struct not just the
// type. For example r.Children is the type but r.Children{} is a struct of that
// type. Only the latter will work.
//
//  func Wrapper(p r.Properties) r.Element {
//    children := p.GetProperty(
//      r.Children{},
//      "Container requires r.Children",
//    ).(r.Children)
//
//     return r.CreateElement(
//      component.Box,
//      r.Properties{
//          component.BoxProps{

//            Border: component.Border{
//              Style:      component.BorderStyleSingle,
//              Foreground: tcell.ColorWhite,
//            },
//          },
//        },
//      children,
//    )
//  }
func (props Properties) GetProperty(
	propType interface{},
	errorMessage string,
) interface{} {
	for _, p := range props {
		if reflect.TypeOf(p) == reflect.TypeOf(propType) {
			return p
		}
	}
	debug.Spew(props)
	panic(errorMessage)
}

// GetOptionalProperty will search props for the Property matching the type of
// struct you passed in. If it was not in props, the struct passed into propType
// will be returned.
//
// You need to cast the return type of the function exactly the same as the
// struct you pass in.
//
// This allows you to specify a defaults for a property.
//
// In the following example if Wrapper is not passed a Property of the type
// component.BoxProps, the default values provided will be used.
//
//  func Wrapper(p r.Properties) r.Element {
//    boxProps := p.GetOptionalProperty(
//      component.BoxProps{
//        Border: component.Border{
//          Style:      component.BorderStyleSingle,
//          Foreground: tcell.ColorWhite,
//        },
//      },
//    ).(component.BoxProps)
//
//     return r.CreateElement(
//      component.Box,
//      r.Properties{
//          boxProps
//        },
//      children,
//    )
//  }
func (props Properties) GetOptionalProperty(
	propType interface{},
) interface{} {
	for _, p := range props {
		if reflect.TypeOf(p) == reflect.TypeOf(propType) {
			return p
		}
	}
	return propType
}

// filterProps returns all props except the type you pass.
func filterProps(props Properties, prop interface{}) Properties {
	newProps := Properties{}

	for _, p := range props {
		if reflect.TypeOf(p) != reflect.TypeOf(prop) {
			newProps = append(newProps, p)
		}
	}
	return newProps
}

// ReplaceProps by replacing with the same type you passed.
func ReplaceProps(props Properties, prop interface{}) Properties {
	newProps := Properties{prop}

	for _, p := range props {
		if reflect.TypeOf(p) != reflect.TypeOf(prop) {
			newProps = append(newProps, p)
		}
	}
	return newProps
}

// AddPropsIfNone will add the prop to props is no existing prop of that type
// is found.
func AddPropsIfNone(props Properties, prop interface{}) Properties {
	foundProp := false

	for _, p := range props {
		if reflect.TypeOf(p) == reflect.TypeOf(prop) {
			foundProp = true
		}
	}

	if !foundProp {
		return append(props, prop)
	}

	return props
}
