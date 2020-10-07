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

	template = `
	def main(text, items):
		return text in items
	`
	o["items"] = []string{"oh", "hi", "mark"}
	o["text"] = "hi"
	if !f.Passes(o, template) {
		t.Error("[text in items] must be true but it is not")
	}
	o["text"] = "hello"
	if f.Passes(o, template) {
		t.Error("[text in items] must be false but it is not")
	}

	template = `
	def main(res, **kwargs):
		return res
	`
	o["res"] = true
	if !f.Passes(o, template) {
		t.Error("[identity] must be true but it is not")
	}
	o["res"] = false
	if f.Passes(o, template) {
		t.Error("[identity] must be false but it is not")
	}
}
