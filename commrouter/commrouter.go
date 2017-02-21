package commrouter

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// CommRouterStruct is the main object to which you can add routes and controllers
type CommRouterStruct struct {
	mu       *sync.RWMutex
	subjects map[string]*SubjectStruct
	Common
}

// Usage outputs the usage of the command router
func (rt *CommRouterStruct) Usage() string {
	var response bytes.Buffer
	for subName, sub := range rt.subjects {
		response.WriteString(fmt.Sprintf("Subject *%v*:\n", subName))
		coms := map[uintptr]Common{}
		for comName, com := range sub.commands {
			comPtr := reflect.ValueOf(com.controller)

			_, ok := coms[comPtr.Pointer()]
			if ok {
				info := coms[comPtr.Pointer()]
				info.Description = com.Description
				info.Examples = com.Examples
				info.names = append(info.names, comName)
				coms[comPtr.Pointer()] = info
			} else {
				newInfo := Common{}
				newInfo.Description = com.Description
				newInfo.Examples = com.Examples
				newInfo.names = []string{comName}
				coms[comPtr.Pointer()] = newInfo
			}
		}

		for _, com := range coms {
			comNames := strings.Join(com.names, "|")
			response.WriteString(fmt.Sprintf("\t*%v*: %v\n", comNames, com.Description))
			if len(com.Examples) > 0 {
				response.WriteString("\t\t_Examples_:")
				for _, example := range com.Examples {
					response.WriteString(fmt.Sprintf(" `%v`", example))
				}
				response.WriteString("\n")
			}
		}
	}

	return response.String()
}

// Common defines common descriptive elements like description, example
type Common struct {
	names       []string
	Description string
	Examples    []string
}

// SubjectStruct is an entity type on which you issue commands;
type SubjectStruct struct {
	mu       *sync.RWMutex
	commands map[string]*CommandStruct
	Common
}

// CommandStruct is an action to be performed
type CommandStruct struct {
	controller         func(interface{}) error
	selectorValidators requirementsMap
	paramValidators    requirementsMap
	Common
}

// Ctx is the context that the controller receives
type Ctx struct {
	args      []string
	selectors H
	params    H
	extra     interface{}
	usage     string
}

func NewCtx() CtxType {
	return &Ctx{}
}

type CtxType interface {
	Args() []string
	Selectors() H
	Params() H
	Extra() interface{}
	RouterUsage() string

	setArgs([]string)
	setSelectors(H)
	setParams(H)
	setExtra(interface{})
	setUsage(string)
}

func (ctx Ctx) Args() []string {
	return ctx.args
}
func (ctx Ctx) Selectors() H {
	return ctx.selectors
}
func (ctx Ctx) Params() H {
	return ctx.params
}
func (ctx Ctx) Extra() interface{} {
	return ctx.extra
}
func (ctx Ctx) RouterUsage() string {
	return ctx.usage
}

func (ctx *Ctx) setArgs(v []string) {
	ctx.args = v
}
func (ctx *Ctx) setSelectors(v H) {
	ctx.selectors = v
}
func (ctx *Ctx) setParams(v H) {
	ctx.params = v
}
func (ctx *Ctx) setExtra(v interface{}) {
	ctx.extra = v
}
func (ctx *Ctx) setUsage(v string) {
	ctx.usage = v
}

// H is a map with string as key and string as value
type H map[string]string

func (h H) GetString(key string) (string, error) {
	raw, ok := h[key]
	if !ok || raw == "" { // TODO: empty == error???
		return "", ErrorKeyNotSet
	}
	return raw, nil
}

func (h H) GetInt(key string) (int, error) {
	raw, ok := h[key]
	if !ok {
		return 0, ErrorKeyNotSet
	}
	return strconv.Atoi(raw)
}

var ErrorKeyNotSet = errors.New("key not set")

// ControllerType is a function that runs the command defined on a subject
type ControllerType func(ctx interface{}) error

// New returns a pointer to an initialized new Router
func New() *CommRouterStruct {
	return &CommRouterStruct{
		mu:       &sync.RWMutex{},
		subjects: make(map[string]*SubjectStruct),
	}
}

// Execute finds the controller and executes it
func (r *CommRouterStruct) Execute(request string, customCtx CtxType) error {
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

	var context CtxType

	if customCtx != nil {
		context = customCtx
	} else {
		context = NewCtx()
	}

	if context == nil {
		return errors.New("context is nil")
	}

	context.setArgs(args)
	context.setParams(params)
	context.setSelectors(selectors)
	context.setUsage(r.Usage())

	return ctrl(context)
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
func (c *CommandStruct) Controller(ctrl func(interface{}) error) error {
	if ctrl == nil {
		return errors.New("controller is nil")
	}
	c.controller = ctrl
	return nil
}
