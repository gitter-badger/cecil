package tools

import (
	"bytes"
	"html/template"

	"github.com/tleyden/cecil/emailtemplates"
	"github.com/tleyden/cecil/gotemplates"
)

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

// CompileGoTemplate will compile a golang template from go-bindata with the provided values
func CompileGoTemplate(name string, values map[string]interface{}) (*bytes.Buffer, error) {
	var compiledTemplate bytes.Buffer

	templateBytes, err := gotemplates.Asset(name)
	if err != nil {
		return nil, err
	}

	tpl := template.New("new go template")
	tpl, err = tpl.Parse(string(templateBytes))
	if err != nil {
		return nil, err
	}

	err = tpl.Execute(&compiledTemplate, values)
	if err != nil {
		return nil, err
	}

	return &compiledTemplate, nil
}
