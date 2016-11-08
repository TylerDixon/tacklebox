package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
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

var allowedVariableCharacters = `[a-zA-Z-_0-9]+`
var templateEncapsulationRegexp = regexp.MustCompile(`{%[\sa-zA-Z()\\:?",]*%}`)
var directiveVariableMatcher = regexp.MustCompile(`(if|switch|render)\((` + allowedVariableCharacters + `)\)`)

// Retrieves a template by name given a list of templates
func getTemplateByName(templates []Template, name string) (Template, error) {
	for _, template := range templates {
		if template.Name == name {
			return template, nil
		}
	}
	return Template{}, fmt.Errorf("Couldn't find template with name %s", name)

}

// Render takes in a map of settings, and with the read in template from the Location property, renders
// the template with the available directives (if, switch, render), returning a final byte array of the
// rendered template
func (template *Template) Render(settings map[string]interface{}) ([]byte, error) {
	// Read in template
	readTemplate, readTemplateError := ioutil.ReadFile(template.Location)
	if readTemplateError != nil {
		error := fmt.Errorf("Failed to read in template of name %s at location %s due to error: %s", template.Name, template.Location, readTemplateError)
		fmt.Println(error)
		return nil, error
	}
	parsedTemplate := string(readTemplate)

	// Find all
	directives := templateEncapsulationRegexp.FindAllStringIndex(parsedTemplate, -1)
	replacementMap, replacementMapError := GetReplacementMap(parsedTemplate, directives, settings)
	if replacementMapError != nil {
		return nil, replacementMapError
	}
	var renderedTemplateToJoin []string
	for n, replacement := range replacementMap {
		var segmentBefore string
		if n == 0 {
			segmentBefore = parsedTemplate[:directives[n][0]]
		} else if n > 0 && n < len(replacementMap) {
			segmentBefore = parsedTemplate[directives[n-1][1]:directives[n][0]]
		}
		renderedTemplateToJoin = append(renderedTemplateToJoin, segmentBefore)

		renderedTemplateToJoin = append(renderedTemplateToJoin, replacement)

		if n == len(replacementMap)-1 {
			segmentAfter := parsedTemplate[directives[n][1]:]
			renderedTemplateToJoin = append(renderedTemplateToJoin, segmentAfter)
		}
	}
	return []byte(strings.Join(renderedTemplateToJoin, "")), nil
}

// GetReplacementMap takes in a template to parse, a set of integer pairs to reference the directive locations in the
// string, and a map of settings. It returns an array of rendered directives, each render indexed the same as the range
// pair where it should be replaced.
func GetReplacementMap(parsedTemplate string, encapsulationPairs [][]int, settings map[string]interface{}) ([]string, error) {
	var renderedTemplates []string
	for _, pair := range encapsulationPairs {
		templatePortion := parsedTemplate[pair[0]:pair[1]]
		if directiveVariableMatcher.MatchString(templatePortion) {
			directiveAndVariable := directiveVariableMatcher.FindStringSubmatch(templatePortion)
			directive := directiveAndVariable[1]
			variable := settings[directiveAndVariable[2]]
			var renderedString string
			var renderError error
			switch directive {
			case `if`:
				renderedString, renderError = RenderIf(variable, templatePortion)
				break
			case `switch`:
				renderedString, renderError = RenderSwitch(variable, templatePortion)
				break
			case `render`:
				renderedString, renderError = RenderLiteralRender(variable, templatePortion)
				break
			}
			if renderError != nil {
				return nil, renderError
			}
			renderedTemplates = append(renderedTemplates, renderedString)
		} else {
			error := fmt.Errorf("No directive matched for the string %s", templatePortion)
			fmt.Println(error)
			return nil, error
		}
	}
	return renderedTemplates, nil
}

// RenderLiteralRender takes in an empty interface, and if it's a string, boolean, or float64, return it's stringified
// value.
func RenderLiteralRender(value interface{}, template string) (string, error) {
	if str, ok := value.(string); ok {
		return str, nil
	} else if boolean, ok := value.(bool); ok {
		return strconv.FormatBool(boolean), nil

	} else if float, ok := value.(float64); ok {
		return strconv.FormatFloat(float, 'f', -1, 64), nil
	}
	renderError := errors.New("Literal render() variable must be one of type: string, bool, float64")
	fmt.Println(renderError)
	return "", renderError
}

// Render if takes a boolean value, and given a directive such as "IfTrue" ? "IfFalse"
// "IfTrue" would be rendered if the value was true, "IfFalse" were it false.
func RenderIf(value interface{}, template string) (string, error) {
	if boolean, ok := value.(bool); ok {
		caseMatching := regexp.MustCompile(`"((?:[^"]|(?:\\"))*)"\s*\?\s*"((?:[^"]|(?:\\"))*)*"`)
		trueAndFalse := caseMatching.FindStringSubmatch(template)
		if boolean {
			return trueAndFalse[1], nil
		} else {
			return trueAndFalse[2], nil
		}
	} else {
		renderError := errors.New("if() variable must be of type bool")
		return "", renderError
	}
}

//TODO: Implement switch rendering
func RenderSwitch(value interface{}, template string) (string, error) {
	//const allowedCaseString = `[^)]*`
	//caseMatching := regexp.MustCompile(`{%
	//\s*switch(` + allowedVariableCharacters + `)\s*
	//(case\("?` + allowedCaseString + `"?\):
	//}`)
	return RenderLiteralRender(value, template)
}
