package main

import (
	"testing"
	"io/ioutil"
	"bytes"
)

func TestBootstrap(t *testing.T) {
	var cases = []struct {
		template Template
		values   map[string]interface{}
		rendered string
	}{
		{
			template: Template{
				Name: "Test",
				Location: "test/templates/simple.txt",
			},
			rendered: "test/templates/simple_out.txt",
			values: map[string]interface{}{
				"TestRender": "Rad",
				"TestIf": true,
				"TestSwitch": "IsThis",
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