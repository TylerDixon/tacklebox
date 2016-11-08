package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestRender(t *testing.T) {
	var cases = []struct {
		template Template
		values   map[string]interface{}
		rendered string
	}{
		{
			template: Template{
				Name:     "Test",
				Location: "test/templates/ifrender.txt",
			},
			rendered: "test/templates/ifrender_out.txt",
			values: map[string]interface{}{
				"TestIf": true,
			},
		},
		{
			template: Template{
				Name:     "Test",
				Location: "test/templates/literalrender.txt",
			},
			rendered: "test/templates/literalrender_string_out.txt",
			values: map[string]interface{}{
				"TestRender": "Rad",
			},
		},
		{
			template: Template{
				Name:     "Test",
				Location: "test/templates/literalrender.txt",
			},
			rendered: "test/templates/literalrender_number_out.txt",
			values: map[string]interface{}{
				"TestRender": 42.3,
			},
		},
		{
			template: Template{
				Name:     "Test",
				Location: "test/templates/literalrender.txt",
			},
			rendered: "test/templates/literalrender_boolean_out.txt",
			values: map[string]interface{}{
				"TestRender": true,
			},
		},
	}

	for _, c := range cases {
		renderedTemplate, _ := c.template.Render(c.values)
		readTemplate, _ := ioutil.ReadFile(c.rendered)
		if bytes.Compare(readTemplate, renderedTemplate) != 0 {
			t.Errorf("Expected the following templates to equal: \n\n%s\n============\n%s", renderedTemplate, readTemplate)
		}

	}
}
