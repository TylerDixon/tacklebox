package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"
)

type RenderMap struct {
	Render         func() string
	setting        interface{}
	stringToRender string
}

type Template struct {
	Name     string
	Location string
}

// Retrieves a templateToRender by name given a list of templates
func getTemplateByName(templates []Template, name string) (Template, error) {
	for _, foundTemplate := range templates {
		if foundTemplate.Name == name {
			return foundTemplate, nil
		}
	}
	return Template{}, fmt.Errorf("Couldn't find templateToRender with name %s", name)

}

// Render takes in a map of settings, and with the read in templateToRender from the Location property, renders
// the templateToRender with the available directives (if, switch, render), returning a final byte array of the
// rendered templateToRender
func (templateToRender *Template) Render(settings map[string]interface{}) ([]byte, error) {
	// Read in templateToRender
	readTemplate, readTemplateError := ioutil.ReadFile(templateToRender.Location)
	if readTemplateError != nil {
		errorToReturn := fmt.Errorf("Failed to read in templateToRender of name %s at location %s due to error: %s", templateToRender.Name, templateToRender.Location, readTemplateError)
		fmt.Println(errorToReturn)
		return nil, errorToReturn
	}
	parsedTemplate := string(readTemplate)
	var renderedTemplate bytes.Buffer
	definedTemplate, err := template.New(templateToRender.Name).Parse(parsedTemplate)
	if err != nil {
		errorToReturn := fmt.Errorf("Failed to parse template of name %s at location %s due to error: %s", templateToRender.Name, templateToRender.Location, readTemplateError)
		fmt.Println(errorToReturn)
		return nil, errorToReturn
	}
	err = definedTemplate.Execute(&renderedTemplate, settings)
	if err != nil {
		errorToReturn := fmt.Errorf("Failed to execute template of name %s at location %s due to error: %s", templateToRender.Name, templateToRender.Location, readTemplateError)
		fmt.Println(errorToReturn)
		return nil, errorToReturn
	}
	return renderedTemplate.Bytes(), nil
}
