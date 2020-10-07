package starlarkfilter

import (
	"testing"
)

func TestStarlarkFilter(t *testing.T) {
	f := StarlarkFilter{}

	o := map[string]interface{}{}
	template := `
	def main(text):
		return text == 'good'
	`

	o["text"] = "good"
	if !f.Passes(o, template) {
		t.Error("must be true but it is not")
	}

	o["text"] = "not-good"
	if f.Passes(o, template) {
		t.Error("must be false but it is not")
	}

	o["text"] = 1
	if f.Passes(o, template) {
		t.Error("must be false but it is not")
	}
}
