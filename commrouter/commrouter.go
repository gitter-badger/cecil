package commrouter

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// CommRouterStruct is the main object to which you can add routes and controllers
type CommRouterStruct struct {
	mu       *sync.RWMutex
	subjects map[string]*SubjectStruct
	Common
}

// Common defines common descriptive elements like description, example
type Common struct {
	Description string
	Example     string
}

// SubjectStruct is an entity type on which you issue commands;
type SubjectStruct struct {
	mu       *sync.RWMutex
	commands map[string]*CommandStruct
	Common
}

// CommandStruct is an action to be performed
type CommandStruct struct {
	controller         ControllerType
	selectorValidators requirementsMap
	paramValidators    requirementsMap
	Common
}

// Ctx is the context that the controller receives
type Ctx struct {
	Args      []string
	Selectors H
	Params    H
}

// H is a map with string as key and interface as value
type H map[string]string

// ControllerType is a function that runs the command defined on a subject
type ControllerType func(ctx *Ctx) error

// New returns a pointer to an initialized new Router
func New() *CommRouterStruct {
	return &CommRouterStruct{
		mu:       &sync.RWMutex{},
		subjects: make(map[string]*SubjectStruct),
	}
}

// Execute finds the controller and executes it
func (r *CommRouterStruct) Execute(request string) error {
	subject, command, args, selectors, params, err := parseRequest(request)
	if err != nil {
		return err
	}

	sub, ok := r.subjects[subject]
	if !ok {
		return fmt.Errorf("subject %q not found", subject)
	}

	comm, ok := sub.commands[command]
	if !ok {
		return fmt.Errorf("command %q for subject %q not found", command, subject)
	}

	if comm.controller == nil {
		return fmt.Errorf("controller is nil for command %q for subject %q", command, subject)
	}

	// TODO: validate here the selectors and params
	// check for required elements that are not there
	// check if the existing elements do have some validation to be made

	err = comm.RequiredSelectorsAreSet(selectors)
	if err != nil {
		return err
	}
	err = comm.RequiredParamsAreSet(params)
	if err != nil {
		return err
	}

	err = comm.ValidateSelectors(selectors)
	if err != nil {
		return err
	}
	err = comm.ValidateParams(params)
	if err != nil {
		return err
	}

	ctrl := comm.controller

	ctx := &Ctx{
		Args:      args,
		Params:    params,
		Selectors: selectors,
	}

	return ctrl(ctx)
}

func parseRequest(request string) (
	Subject string,
	Command string,
	Args []string,
	Selectors H,
	Params H,
	err error,
) {
	Selectors = make(H)
	Params = make(H)

	defer func() {
		Subject = strings.TrimSpace(Subject)
		Command = strings.TrimSpace(Command)
	}()

	request = strings.TrimSpace(request)
	indexOfFirstSpace := strings.Index(request, " ")

	for {
		if strings.Contains(request, "  ") {
			request = strings.Replace(request, "  ", " ", -1)
		} else {
			break
		}
	}

	if strings.Count(request, " ") > 0 {
		Command = request[:indexOfFirstSpace]
	} else {
		err = errors.New("subject not specified")
		return
	}

	if strings.Count(request, " ") == 1 {
		request = request[indexOfFirstSpace:]
		Subject = request
	} else if strings.Count(request, " ") > 1 {
		request = request[indexOfFirstSpace:]
		request = strings.TrimSpace(request)
		indexOfSecondSpace := strings.Index(request, " ")
		Subject = request[:indexOfSecondSpace]
		request = request[indexOfSecondSpace:]
		request = strings.TrimSpace(request)
	} else {
		err = errors.New("subject not specified")
		return
	}

	segments := strings.Fields(request)
	indexIsPastArgs := false

	for _, seg := range segments {
		if strings.Contains(seg, ":") && strings.Contains(seg, "=") {
			err = fmt.Errorf("segment %q contains both separators", seg)
			break
		}
		if !strings.Contains(seg, ":") && !strings.Contains(seg, "=") {
			if !indexIsPastArgs {
				Args = append(Args, seg)
			} else {
				err = fmt.Errorf("segment %q has no separator; args can be put ONLY after the subject; selectors and params do not support args", seg)
				break
			}
		}

		if strings.Contains(seg, "=") {
			indexIsPastArgs = true
			key, value, e := Selectors.parseSegment(seg, "=")
			if e != nil {
				err = e
				break
			}
			Selectors[key] = value
		}

		if strings.Contains(seg, ":") {
			indexIsPastArgs = true
			key, value, e := Params.parseSegment(seg, ":")
			if e != nil {
				err = e
				break
			}
			Params[key] = value
		}
	}

	// TODO: check each return value
	return
}

func (h H) parseSegment(segment string, separator string) (key, value string, err error) {
	if strings.Count(segment, separator) > 1 {
		err = fmt.Errorf("segment %q contains repeating separators", segment)
		return
	}
	kv := strings.Split(segment, separator)
	if len(kv) > 2 {
		err = fmt.Errorf("segment %q contains too many separators", segment)
		return
	}
	key = kv[0]
	key = strings.TrimSpace(key)
	if len(key) == 0 {
		err = fmt.Errorf("segment %q contains empty key", segment)
		return
	}
	value = kv[1]
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		err = fmt.Errorf("segment %q contains empty value", segment)
		return
	}
	if _, ok := h[key]; ok {
		err = fmt.Errorf("key %q already declared", key)
		return
	}
	return
}

// AddSubject adds a Subject to the router and returns it with eventual error
func (r *CommRouterStruct) AddSubject(name string) (*SubjectStruct, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	newSubject := SubjectStruct{
		mu:       r.mu,
		commands: make(map[string]*CommandStruct),
	}

	// TODO: validate name
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, errors.New("name is empty")
	}

	// check whether the subject already exists
	_, ok := r.subjects[name]
	if ok {
		return nil, fmt.Errorf("Subject %v already exists in router", name)
	}

	// add to router
	r.subjects[name] = &newSubject

	// return pointer
	return &newSubject, nil
}

// AddCommand adds a command on a subject and returns it with eventual error
func (s *SubjectStruct) AddCommand(nameVariations ...string) (*CommandStruct, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newCommand := CommandStruct{
		selectorValidators: make(requirementsMap),
		paramValidators:    make(requirementsMap),
	}

	// TODO: validate name

	for _, name := range nameVariations {
		// validate name
		name = strings.TrimSpace(name)
		if len(name) == 0 {
			return nil, errors.New("name is empty")
		}
		// check whether the command already exists
		_, ok := s.commands[name]
		if ok {
			return nil, fmt.Errorf("Command %v already exists in subject", name)
		}

		// add to subject
		s.commands[name] = &newCommand
	}

	// return pointer
	return &newCommand, nil
}

// Controller lets you define the controller function of the command on the subject
func (c *CommandStruct) Controller(ctrl ControllerType) error {
	if ctrl == nil {
		return errors.New("controller is nil")
	}
	c.controller = ctrl
	return nil
}
