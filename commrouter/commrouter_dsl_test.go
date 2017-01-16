package commrouter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSLNew(t *testing.T) {
	newRouter := CommRouter()

	assert.NotNil(t, newRouter.mu, "mutex should not be nil")
	assert.NotNil(t, newRouter.subjects, "subjects should not be nil")
}

func TestDSLAddSubject(t *testing.T) {
	subjectName := "lease"
	subjectDescription := "A Lease defines the lease of an instance"
	subjectExample := "some example"

	newRouter := CommRouter(
		func() {
			Description("this is the description of the CommRouter")
			Subject(
				subjectName,
				func() {
					Description(subjectDescription)
					Example(subjectExample)
				},
			)
		},
	)

	gotSubject, ok := newRouter.subjects[subjectName]
	assert.True(t, ok, "there should exist a subject for that subjectName")

	assert.Equal(t, subjectDescription, gotSubject.Description, "should be equal")
	assert.Equal(t, subjectExample, gotSubject.Example, "should be equal")
}

func TestDSLAddCommand(t *testing.T) {
	subjectName := "lease"
	commandVariations := []string{"list", "show", "display"}
	subjectDescription := "A Lease defines the lease of an instance"
	subjectExample := "some example"
	commandDescription := "Display one or more leases"
	commandExample := "show lease 1"

	var ctrl ControllerType = func(ctx *Ctx) error {
		fmt.Println("This thing works: bake the pie", ctx)
		return nil
	}

	newRouter := CommRouter(
		func() {
			Description("this is the description of the CommRouter")
			Subject(
				subjectName,
				func() {
					Description(subjectDescription)
					Example(subjectExample)

					Command(
						commandVariations,
						func() {
							Description(commandDescription)
							Example(commandExample)
							Controller(ctrl)
						},
					)
				},
			)
		},
	)

	gotSubject, ok := newRouter.subjects[subjectName]
	assert.True(t, ok, "there should exist a subject for that subjectName")

	assert.Equal(t, subjectDescription, gotSubject.Description, "should be equal")
	assert.Equal(t, subjectExample, gotSubject.Example, "should be equal")

	for _, command := range commandVariations {
		gotCommand, ok := gotSubject.commands[command]
		assert.True(t, ok, "there should exist a subject for that commandVariation")

		sf1 := reflect.ValueOf(ctrl)
		sf2 := reflect.ValueOf(gotCommand.controller)
		assert.True(t, sf1.Pointer() == sf2.Pointer(), "controllers should match")

		err := newRouter.Execute(fmt.Sprintf("%v %v 1 selector=something param:else ", command, subjectName))
		assert.Nil(t, err, "error should be nil")
	}

}

func TestDSLParseRequest(t *testing.T) {
	subjectName := "lease"
	commandVariations := []string{"list", "show", "display"}
	subjectDescription := "A Lease defines the lease of an instance"
	subjectExample := "some example"
	commandDescription := "Display one or more leases"
	commandExample := "show lease 1"

	var ctrl ControllerType = func(ctx *Ctx) error {
		fmt.Println("This thing works: bake the pie", ctx)
		return nil
	}

	newRouter := CommRouter(
		func() {
			Description("this is the description of the CommRouter")
			Subject(
				subjectName,
				func() {
					Description(subjectDescription)
					Example(subjectExample)

					Command(
						commandVariations,
						func() {
							Description(commandDescription)
							Example(commandExample)
							Controller(ctrl)
						},
					)
				},
			)
		},
	)

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
		err := newRouter.Execute(request)
		assert.Nil(t, err, "error should be nil")

		subject, command, args, selectors, params, err := parseRequest(request)
		assert.Nil(t, err, "error should be nil")

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
	_, _, _, _, _, err := parseRequest(request)
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

func TestDSLExecute(t *testing.T) {
	subjectName := "lease"
	commandVariations := []string{"list", "show", "display"}
	subjectDescription := "A Lease defines the lease of an instance"
	subjectExample := "some example"
	commandDescription := "Display one or more leases"
	commandExample := "show lease 1"

	var ctrl ControllerType = func(ctx *Ctx) error {
		fmt.Println("This thing works: bake the pie", ctx)
		return nil
	}

	newRouter := CommRouter(
		func() {
			Description("this is the description of the CommRouter")
			Subject(
				subjectName,
				func() {
					Description(subjectDescription)
					Example(subjectExample)

					Command(
						commandVariations,
						func() {
							Description(commandDescription)
							Example(commandExample)
							Controller(ctrl)
						},
					)
				},
			)
		},
	)

	request := "show pie"
	err := newRouter.Execute(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = fmt.Sprintf("bake %v", subjectName)
	err = newRouter.Execute(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = fmt.Sprintf("bake %v", subjectName)
	err = newRouter.Execute(request)
	assert.NotNil(t, err, "there SHOULD be an error")

	request = fmt.Sprintf("show %s", subjectName)
	err = newRouter.Execute(request)
	assert.Nil(t, err, "error should be nil")
}
