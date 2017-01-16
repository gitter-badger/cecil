package commrouter

import (
	"fmt"
	"regexp"
)

// I is an interface
type I interface{}

// currentDefinition is the current working element
var currentDefinition interface{}

// CommRouter defines a new CommRouter
func CommRouter(Is ...I) *CommRouterStruct {
	newRouter := New()
	currentDefinition = newRouter

	for _, i := range Is {
		switch t := i.(type) {
		case func():
			{
				i.(func())()
			}
		default:
			panic(fmt.Sprintf("%#v", t))
		}
	}

	return newRouter
}

// Subject declares a new subject in the CommRouter
func Subject(name string, Is ...I) {

	previousDefinition := currentDefinition
	defer func() {
		currentDefinition = previousDefinition
	}()

	switch t := currentDefinition.(type) {
	case *CommRouterStruct:
		rr := currentDefinition.(*CommRouterStruct)
		newSubject, err := rr.AddSubject(name) // add new subject to router
		if err != nil {
			panic(fmt.Sprintf("%#v", t))
		}
		currentDefinition = newSubject
	}

	for _, i := range Is {
		switch t := i.(type) {
		case func():
			{
				i.(func())()
			}
		default:
			panic(fmt.Sprintf("%#v", t))
		}
	}

}

// Command defines a command on a subject
func Command(Is ...I) {

	previousDefinition := currentDefinition
	defer func() {
		currentDefinition = previousDefinition
	}()

	for _, i := range Is {
		switch t := i.(type) {
		case func():
			{
				i.(func())()
			}
		case []string:
			{
				switch currentDefinition.(type) {
				case *SubjectStruct:
					ss := currentDefinition.(*SubjectStruct)
					newCommand, err := ss.AddCommand(i.([]string)...) // add new subject to CommRouter
					if err != nil {
						panic(fmt.Sprintf("%#v", t))
					}
					currentDefinition = newCommand
				}
			}
		default:
			panic(fmt.Sprintf("%#v", t))
		}
	}

}

// Controller defines the controller for a command
func Controller(ctrl ControllerType) {
	switch t := currentDefinition.(type) {
	case *CommandStruct:
		currentDefinition.(*CommandStruct).Controller(ctrl)
	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// Params defines a list of parameters and their requirements
func Params(Is ...I) {
	switch t := currentDefinition.(type) {
	case *CommandStruct:
		//cmd := currentDefinition.(*CommandStruct)

		for _, i := range Is {
			switch t := i.(type) {
			case func():
				{
					i.(func())()
				}
			default:
				panic(fmt.Sprintf("%#v", t))
			}
		}

	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// Selectors defines a list of selectors and their requirements
func Selectors(Is ...I) {
	switch t := currentDefinition.(type) {
	case *CommandStruct:
		//cmd := currentDefinition.(*CommandStruct)

		for _, i := range Is {
			switch t := i.(type) {
			case func():
				{
					i.(func())()
				}
			default:
				panic(fmt.Sprintf("%#v", t))
			}
		}

	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// Param is a single param with its requirements
func Param(key string, paramType Type, Is ...I) {
	previousDefinition := currentDefinition
	defer func() {
		currentDefinition = previousDefinition
	}()

	switch t := currentDefinition.(type) {
	case *CommandStruct:
		cmd := currentDefinition.(*CommandStruct)

		req := &Requirements{}
		req.Type = paramType
		currentDefinition = req

		for _, i := range Is {
			switch t := i.(type) {
			case func():
				{
					i.(func())()
				}
			default:
				panic(fmt.Sprintf("%#v", t))
			}
		}

		cmd.AddParamRequirement(key, currentDefinition.(*Requirements))

	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// Selector is a single selector with its requirements
func Selector(key string, paramType Type, Is ...I) {
	previousDefinition := currentDefinition
	defer func() {
		currentDefinition = previousDefinition
	}()

	switch t := currentDefinition.(type) {
	case *CommandStruct:
		cmd := currentDefinition.(*CommandStruct)

		req := &Requirements{}
		req.Type = paramType
		currentDefinition = req

		for _, i := range Is {
			switch t := i.(type) {
			case func():
				{
					i.(func())()
				}
			default:
				panic(fmt.Sprintf("%#v", t))
			}
		}

		cmd.AddSelectorRequirement(key, currentDefinition.(*Requirements))

	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// Required tells that this field is required
func Required() {
	switch t := currentDefinition.(type) {
	case *Requirements:
		req := currentDefinition.(*Requirements)
		req.Required = true
	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// MinValue sets the min value for the field
func MinValue(val interface{}) {
	switch t := currentDefinition.(type) {
	case *Requirements:
		req := currentDefinition.(*Requirements)
		req.MinValue = val
	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// MinLength sets the min length
func MinLength(length int) {
	switch t := currentDefinition.(type) {
	case *Requirements:
		req := currentDefinition.(*Requirements)
		req.MinLength = &length
	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// MustRegex sets the regexp that the value must match
func MustRegex(regex *regexp.Regexp) {
	switch t := currentDefinition.(type) {
	case *Requirements:
		req := currentDefinition.(*Requirements)
		req.MustRegex = regex
	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// Description sets the description
func Description(description string) {
	switch t := currentDefinition.(type) {
	case *CommRouterStruct:
		currentDefinition.(*CommRouterStruct).Description = description
	case *SubjectStruct:
		currentDefinition.(*SubjectStruct).Description = description
	case *CommandStruct:
		currentDefinition.(*CommandStruct).Description = description
	default:
		panic(fmt.Sprintf("%#v", t))
	}
}

// Example sets the example
func Example(example string) {
	switch t := currentDefinition.(type) {
	case *CommRouterStruct:
		currentDefinition.(*CommRouterStruct).Example = example
	case *SubjectStruct:
		currentDefinition.(*SubjectStruct).Example = example
	case *CommandStruct:
		currentDefinition.(*CommandStruct).Example = example
	default:
		panic(fmt.Sprintf("%#v", t))
	}
}
