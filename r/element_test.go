package r

import (
	"reflect"
	"testing"
)

func TestReplaceProps(t *testing.T) {
	type ExampleProperty struct {
		Value int
	}
	type ExampleTextProperty struct {
		Value string
	}

	var cases = []struct {
		Props         Properties
		NewProp       interface{}
		ExpectedProps Properties
	}{
		{
			Props: Properties{
				ExampleProperty{
					Value: 1,
				},
			},
			NewProp: ExampleProperty{
				Value: 2,
			},
			ExpectedProps: Properties{
				ExampleProperty{
					Value: 2,
				},
			},
		},
		{
			Props: Properties{
				ExampleProperty{
					Value: 425,
				},
				ExampleTextProperty{
					Value: "test",
				},
			},
			NewProp: ExampleProperty{
				Value: 234567,
			},
			ExpectedProps: Properties{
				ExampleProperty{
					Value: 234567,
				},
				ExampleTextProperty{
					Value: "test",
				},
			},
		},
	}

	for i, c := range cases {
		actual := ReplaceProps(c.Props, c.NewProp)
		if !reflect.DeepEqual(actual, c.ExpectedProps) {
			t.Errorf("Fib(%d): expected %d, actual %d", i, c.ExpectedProps, actual)
		}
	}
}
