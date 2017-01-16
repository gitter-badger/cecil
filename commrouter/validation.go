package commrouter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Requirements defines the requirements for a specific element (selector or param)
type Requirements struct {
	Type      Type
	Required  bool
	MinValue  interface{}
	MinLength *int
	MustRegex *regexp.Regexp
}

// Type is a type
type Type string

const (
	// String defines the string type for input
	String Type = "string"

	// Int defines the int type for input
	Int Type = "int"

	// Duration defines the time.Duration type for input
	Duration Type = "duration"

	// Any defines the interface{} type for input
	Any Type = "any" // this is string also
)

type requirementsMap map[string]*Requirements

// RequiredParamsAreSet checks whether params contains all the required params
func (c *CommandStruct) RequiredParamsAreSet(params H) error {
	return c.paramValidators.requiredAreSet(params, "param")
}

// RequiredSelectorsAreSet checks whether selectors contains all the required selectors
func (c *CommandStruct) RequiredSelectorsAreSet(selectors H) error {
	return c.selectorValidators.requiredAreSet(selectors, "selector")
}

// ValidateParams validates the provided params
func (c *CommandStruct) ValidateParams(params H) error {
	return c.paramValidators.validate(params, "param")
}

// ValidateSelectors validates the provided selectors
func (c *CommandStruct) ValidateSelectors(selectors H) error {
	return c.selectorValidators.validate(selectors, "selector")
}

// AddParamRequirement adds a requirement for a param
func (c *CommandStruct) AddParamRequirement(key string, r *Requirements) error {
	return c.paramValidators.add(key, r)
}

// AddSelectorRequirement adds a requirement for a selector
func (c *CommandStruct) AddSelectorRequirement(key string, r *Requirements) error {
	return c.selectorValidators.add(key, r)
}

func (rm requirementsMap) add(key string, r *Requirements) error {
	if r.Type == "" {
		return errors.New("Type not specified")
	}

	// if a MinValue is specified, make sure the Type is right
	if r.MinValue != nil {
		switch r.Type {
		case String:
			{
				_, ok := r.MinValue.(string)
				if !ok {
					return errors.New("MinValue is not a string")
				}
			}

		case Int:
			{
				_, ok := r.MinValue.(int)
				if !ok {
					return errors.New("MinValue is not an int")
				}
			}

		case Duration:
			{
				_, ok := r.MinValue.(time.Duration)
				if !ok {
					return errors.New("MinValue is not time.Duration")
				}
			}

		}
	}

	if r.MinLength != nil && r.Type != String {
		return errors.New("MinLength is only for strings")
	}
	if r.MustRegex != nil && r.Type != String {
		return errors.New("MustRegex is only for strings")
	}

	if rm == nil {
		rm = make(requirementsMap)
	}

	// TODO: should it just overwrite?
	_, ok := rm[key]
	if ok {
		return fmt.Errorf("requirements for key %v already set", key)
	}

	rm[key] = r

	return nil
}

// validate parsed elements
func (rm requirementsMap) validate(elements H, elementName string) error {
	for key, val := range elements {
		validator, ok := rm[key]
		if !ok {
			continue
		}

		// validate type by trying to assert

		switch validator.Type {
		case String:
			{
				if validator.MustRegex != nil {
					ok := validator.MustRegex.MatchString(val)
					if !ok {
						return fmt.Errorf("value of %v %q (%q) does not match regexp %v", elementName, key, val, validator.MustRegex.String())
					}
				}
			}
		case Any:
			{
				// ok
			}

		case Int:
			{
				integer, err := strconv.Atoi(val)
				if err != nil {
					return fmt.Errorf("%v %q must be an int", elementName, key)
				}
				if validator.MinValue != nil {
					expected := validator.MinValue.(int)
					got := integer
					if got < expected {
						return fmt.Errorf("%v %q min value is %v, but got %v", elementName, key, expected, got)
					}
				}
			}

		case Duration:
			{
				duration, err := time.ParseDuration(val)
				if err != nil {
					return fmt.Errorf("%v %q must be of time.Duration type", elementName, key)
				}

				if validator.MinValue != nil {
					expected := validator.MinValue.(time.Duration)
					got := duration
					if got < expected {
						return fmt.Errorf("%v %q min value is %v, but got %v", elementName, key, expected, got)
					}
				}
			}

		}

	}
	return nil
}

func (rm requirementsMap) requiredAreSet(elements H, elementName string) error {
	for key, req := range rm {
		if req.Required {
			_, ok := elements[key]
			if !ok {
				return fmt.Errorf("%v %q (%v type) is not set", elementName, key, req.Type)
			}
		}
	}
	return nil
}
