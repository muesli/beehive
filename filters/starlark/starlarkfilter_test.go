package starlarkfilter

import (
	"testing"
)

func TestStarlarkFilter(t *testing.T) {
	t.Parallel()
	f := StarlarkFilter{}

	o := map[string]interface{}{}
	template := `
	def main(text):
		return text == 'good'
	`

	// test string convertion
	o["text"] = "good"
	if !f.Passes(o, template) {
		t.Error("must be true but it is not")
	}
	o["text"] = "not-good"
	if f.Passes(o, template) {
		t.Error("must be false but it is not")
	}

	// test int convertion
	o["text"] = 1
	if f.Passes(o, template) {
		t.Error("must be false but it is not")
	}

	// test slice convertion
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

	// test map convertion
	o["text"] = "hi"
	o["items"] = map[string]bool{"oh": true, "hi": true, "mark": true}
	if !f.Passes(o, template) {
		t.Error("[text in map] must be true but it is not")
	}
	o["items"] = map[string]bool{"oh": true, "hello": true, "mark": true}
	if f.Passes(o, template) {
		t.Error("[text in map] must be false but it is not")
	}

	// test pointer convertion
	items := map[string]bool{"oh": true, "hi": true, "mark": true}
	o["items"] = &items
	if !f.Passes(o, template) {
		t.Error("[text in map] must be true but it is not")
	}

	// test bool convertion
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
