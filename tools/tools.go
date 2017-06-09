// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package tools

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tleyden/cecil/emailtemplates"
)

// HMI is a map with string as key and interface as value
type HMI map[string]interface{}

// HMS is a map with string as key and string as value
type HMS map[string]string

// SchedulePeriodicJob is used to spin up a goroutine that runs
// a specific function in a cycle of specified time
func SchedulePeriodicJob(job func() error, runEvery time.Duration, errorCallback func(err error)) {
	go func() {
		for {
			// run job
			err := job()

			// if error, execute errorCallback
			if err != nil {
				errorCallback(fmt.Errorf("schedulePeriodicJob: %v", err))
			}

			// wait util next execution of job
			time.Sleep(runEvery)
		}
	}()
}

// Retry performs callback n times until error is nil
func Retry(attempts int, sleep time.Duration, callback func() error, intermediateErrorCallback func(error)) (err error) {
	for i := 1; i <= attempts; i++ {

		err = callback()
		if err == nil {
			return nil
		}
		time.Sleep(sleep)

		if err != nil {
			if intermediateErrorCallback != nil {
				intermediateErrorCallback(err)
			}
		}
	}
	return fmt.Errorf("Abandoned after %d attempts, last error: %s", attempts, err)
}

// CompileEmail compiles a template with values
func CompileEmail(tpl string, values map[string]interface{}) string {
	var emailBody bytes.Buffer // A Buffer needs no initialization.

	// TODO: check errors ???

	t := template.New("new email template")
	t, _ = t.Parse(tpl)

	_ = t.Execute(&emailBody, values)

	return emailBody.String()
}

// CompileEmailTemplate will compile a golang template from file (just filename; the folder is hardcoded here) with the provided values
func CompileEmailTemplate(name string, values map[string]interface{}) (string, error) {
	var compiledTemplate bytes.Buffer

	templateBytes, err := emailtemplates.Asset(name)
	if err != nil {
		return "", err
	}

	tpl := template.New("new email template")
	tpl, err = tpl.Parse(string(templateBytes))
	if err != nil {
		return "", err
	}

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		return "", err
	}

	return compiledTemplate.String(), nil
}

// AskForConfirmation waits for stdin input by the user
// in the cli interface. Input yes or not, then enter (newline).
func AskForConfirmation() bool {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Println("fatal: ", err)
	}
	positive := []string{"y", "Y", "yes", "Yes", "YES"}
	negative := []string{"n", "N", "no", "No", "NO"}
	if SliceContains(positive, input) {
		return true
	} else if SliceContains(negative, input) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter.")
		return AskForConfirmation()
	}
}

// SliceContains returns true if slice contains element.
func SliceContains(slice []string, element string) bool {
	for _, elem := range slice {
		if strings.EqualFold(element, elem) {
			return true
		}
	}
	return false
}

// LoudPrint is a fmt.Printf along with rows of '@' characters that make the print visible in the console
func LoudPrint(format string, a ...interface{}) (int, error) {
	m1 := `
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
`

	m2 := "####################  " + fmt.Sprintf(format, a...) + "  ####################"

	m3 := `
############################################################
############################################################
############################################################
`

	return fmt.Println(m1, m2, m3)
}

// IntPtr returns a pointer to the specified int
func IntPtr(i int) *int {
	return &i
}

// StringPtr returns a pointer to the specified string
func StringPtr(s string) *string {
	return &s
}

// TimePtr returns a pointer to the specified time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// DurationPtr returns a pointer to the specified duration
func DurationPtr(t time.Duration) *time.Duration {
	return &t
}

// UUIDPtr returns a pointer to the specified uuid.UUID
func UUIDPtr(t uuid.UUID) *uuid.UUID {
	return &t
}
