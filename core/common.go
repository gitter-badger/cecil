package core

import (
	"bytes"
	"text/template"
	"time"
)

func runForever(f func() error, sleepDuration time.Duration) {
	for {
		err := f()
		if err != nil {
			logger.Error("runForever", err)
		}
		time.Sleep(sleepDuration)
	}
}

func compileEmail(tpl string, values map[string]interface{}) string {
	var emailBody bytes.Buffer // A Buffer needs no initialization.

	// TODO: check errors ???

	t := template.New("new email template")
	t, _ = t.Parse(tpl)

	_ = t.Execute(&emailBody, values)

	return emailBody.String()
}
