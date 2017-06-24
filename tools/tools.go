// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

// IsNumeric returns true if the string is numeric
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Stringify converts a value to a string
func Stringify(v interface{}) string {
	switch v.(type) {
	case time.Time:
		{
			return v.(time.Time).Format(time.RFC3339)
		}
	}
	return fmt.Sprint(v)
}
