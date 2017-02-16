package commrouter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	newRouter := New()

	assert.NotNil(t, newRouter.mu, "mutex should not be nil")
	assert.NotNil(t, newRouter.subjects, "subjects should not be nil")
}

func TestAddSubject(t *testing.T) {
	newRouter := New()

	subjectName := "lease"
	leaseSubject, err := newRouter.AddSubject(subjectName) // add new subject to router
	assert.Nil(t, err, "error should be nil")

	description := "A Lease defines the lease of an instance"
	leaseSubject.Description = description
	example := "some example"
	leaseSubject.Examples = []string{example}

	assert.Equal(t, description, leaseSubject.Description, "should be equal")
	assert.Equal(t, []string{example}, leaseSubject.Examples, "should be equal")

	gotSubject, ok := newRouter.subjects[subjectName]
	assert.True(t, ok, "there should exist a subject for that subjectName")
	assert.Equal(t, leaseSubject, gotSubject, "they should be equal")
}

func TestAddCommand(t *testing.T) {
	newRouter := New()

	subjectName := "lease"
	leaseSubject, err := newRouter.AddSubject(subjectName) // add new subject to router
	assert.Nil(t, err, "error should be nil")

	commandVariations := []string{"list", "show", "display"}
	// declare list command
	listCommand, err := leaseSubject.AddCommand(commandVariations...) // add variations of the spelling of a command to a subject
	assert.Nil(t, err, "error should be nil")

	description := "Display one or more leases"
	listCommand.Description = description
	example := "show lease 1"
	listCommand.Examples = []string{example}

	assert.Equal(t, description, listCommand.Description, "should be equal")
	assert.Equal(t, []string{example}, listCommand.Examples, "should be equal")

	for _, command := range commandVariations {
		gotCommand, ok := leaseSubject.commands[command]
		assert.True(t, ok, "there should exist a subject for that commandVariation")
		assert.Equal(t, listCommand, gotCommand, "they should be equal")
	}

	ctrl := func(ctx interface{}) error {
		fmt.Println("hey! this is a new request to list!", ctx)
		return nil
	}

	listCommand.Controller(ctrl)

	for _, command := range commandVariations {
		err = newRouter.Execute(fmt.Sprintf("%v %v 1 selector=something param:else ", command, subjectName), nil)
		assert.Nil(t, err, "error should be nil")
	}
}

func TestParseRequest(t *testing.T) {
	newRouter := New()

	subjectName := "lease"
	leaseSubject, err := newRouter.AddSubject(subjectName) // add new subject to router
	assert.Nil(t, err, "error should be nil")

	commandVariations := []string{"list", "show", "display"}
	// declare list command
	listCommand, err := leaseSubject.AddCommand(commandVariations...) // add variations of the spelling of a command to a subject
	assert.Nil(t, err, "error should be nil")

	description := "Display one or more leases"
	listCommand.Description = description
	example := "show lease 1"
	listCommand.Examples = []string{example}

	assert.Equal(t, description, listCommand.Description, "should be equal")
	assert.Equal(t, []string{example}, listCommand.Examples, "should be equal")

	for _, command := range commandVariations {
		gotCommand, ok := leaseSubject.commands[command]
		assert.True(t, ok, "there should exist a subject for that commandVariation")
		assert.Equal(t, listCommand, gotCommand, "they should be equal")
	}

	ctrl := func(ctx interface{}) error {
		fmt.Println("hey! this is a new request to list!", ctx)
		return nil
	}

	listCommand.Controller(ctrl)

	for _, commandVariation := range commandVariations {
		selectorKey, selectorValue := "selector", "something"
		paramKey, paramValue := "param", "else"
		arg := "i-12121212"

		request := fmt.Sprintf(
			"%v %v %v %v=%v %v:%v ",
			commandVariation,
			subjectName,
			arg,
			selectorKey,
			selectorValue,
			paramKey,
			paramValue,
		)
		err = newRouter.Execute(request, nil)
		assert.Nil(t, err, "error should be nil")

		subject, command, args, selectors, params, err1 := parseRequest(request)
		assert.Nil(t, err1, "error should be nil")

		assert.Equal(t, subjectName, subject, "subjects should match")
		assert.Equal(t, commandVariation, command, "commands should match")
		assert.Len(t, args, 1, "args should contain one element")
		assert.Len(t, selectors, 1, "selectors should contain one element")
		assert.Len(t, params, 1, "params should contain one element")

		assert.Equal(t, arg, args[0], "args should match")
		assert.Equal(t, selectorValue, selectors[selectorKey], "selector should match")
		assert.Equal(t, paramValue, params[paramKey], "param should match")
	}

	request := "command"
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command     "
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command subject key::value "
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command subject key==value "
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command subject key="
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command subject key:"
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command subject :value"
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command subject =value"
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = "command subject key=value arg"
	_, _, _, _, _, err = parseRequest(request)
	assert.NotNil(t, err, "there SHOULD be an error")
}

func TestExecute(t *testing.T) {
	newRouter := New()

	subjectName := "lease"
	leaseSubject, err := newRouter.AddSubject(subjectName) // add new subject to router
	assert.Nil(t, err, "error should be nil")

	commandVariations := []string{"list", "show", "display"}
	// declare list command
	listCommand, err := leaseSubject.AddCommand(commandVariations...) // add variations of the spelling of a command to a subject
	assert.Nil(t, err, "error should be nil")

	for _, command := range commandVariations {
		gotCommand, ok := leaseSubject.commands[command]
		assert.True(t, ok, "there should exist a subject for that commandVariation")
		assert.Equal(t, listCommand, gotCommand, "they should be equal")
	}

	ctrl := func(ctx interface{}) error {
		fmt.Println("hey! this is a new request to list!", ctx)
		return nil
	}

	request := "show pie"
	err = newRouter.Execute(request, nil)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = fmt.Sprintf("bake %v", subjectName)
	err = newRouter.Execute(request, nil)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = fmt.Sprintf("bake %v", subjectName)
	err = newRouter.Execute(request, nil)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = fmt.Sprintf("show %v", subjectName)
	err = newRouter.Execute(request, nil)
	assert.NotNil(t, err, "there SHOULD be an error")

	listCommand.Controller(ctrl)

	request = fmt.Sprintf("show %s", subjectName)
	err = newRouter.Execute(request, nil)
	assert.Nil(t, err, "error should be nil")
}
