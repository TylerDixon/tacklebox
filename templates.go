package main

import (
	"io/ioutil"
	"fmt"
	"regexp"
	"strings"
	"strconv"
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

var templateEncapsulationRegexp = regexp.MustCompile(`{%[\sa-zA-Z\(\):?",]*%}`)
var directiveVariableMatcher = regexp.MustCompile(`(if|switch|render)\(([a-zA-Z-_]+)\)`)

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
	fmt.Println(replacementMap)
	for n, replacement := range replacementMap {
		var segmentBefore string
		if n == 0 {
			segmentBefore = parsedTemplate[:directives[n][0]]
		} else if n > 0 && n < len(replacementMap) {
			segmentBefore = parsedTemplate[directives[n - 1][1]:directives[n][0]]
		}
		renderedTemplateToJoin = append(renderedTemplateToJoin, segmentBefore)

		renderedTemplateToJoin = append(renderedTemplateToJoin, replacement)

		if n == len(replacementMap)-1 {
			segmentAfter := parsedTemplate[directives[n][1]:]
			renderedTemplateToJoin = append(renderedTemplateToJoin, segmentAfter)
		}
		fmt.Println(renderedTemplateToJoin)
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
			switch directive {
			case `if`: renderedTemplates = append(renderedTemplates, RenderIf(variable, templatePortion)); break
			case `switch`: renderedTemplates = append(renderedTemplates, RenderSwitch(variable, templatePortion)); break
			case `render`: renderedTemplates = append(renderedTemplates, RenderLiteralRender(variable, templatePortion)); break
			}
		} else {
			error := fmt.Errorf("No directive matched for the string %s", templatePortion)
			fmt.Println(error)
			return nil, error
		}
	}
	return renderedTemplates, nil
}

// RenderLiteralRender takes in an empty interface, and if it's a string, boolean, or integer, return it's stringified
// value.
func RenderLiteralRender(value interface{}, template string) string {
	if str, ok := value.(string); ok {
		return str
	} else if boolean, ok := value.(bool); ok {
		return strconv.FormatBool(boolean)

	} else if integer, ok := value.(int64); ok {
		return strconv.FormatInt(integer, 10)
	}
	//TODO: implement error returning?
	return "???"
}

//TODO: Implement switch and if rendering
func RenderSwitch(value interface{}, template string) string {
	return RenderLiteralRender(value, template)
}
func RenderIf(value interface{}, template string) string {
	return RenderLiteralRender(value, template)
}