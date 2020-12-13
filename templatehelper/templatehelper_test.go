package templatehelper

import (
	"bytes"
	"testing"
	"text/template"
)

func executeTemplate(text string, data interface{}) (string, error) {
	var res bytes.Buffer
	tmpl, err := template.New("_test").Funcs(FuncMap).Parse(text)
	if err == nil {
		err = tmpl.Execute(&res, data)
	}
	return res.String(), err
}

func Test_FuncMap_Positive(t *testing.T) {
	t.Parallel()
	cases := []struct {
		text     string
		expected string
	}{
		// sanity checks

		{`{{if true}}ok{{end}}`, "ok"},
		{`{{if false}}ok{{end}}`, ""},

		// boolean filters

		{`{{if Matches "123" "\\d+"}}ok{{end}}`, "ok"},
		{`{{if Matches "hello" "\\d+"}}ok{{end}}`, ""},

		{`{{if Contains "123" "2"}}ok{{end}}`, "ok"},
		{`{{if Contains "123" "4"}}ok{{end}}`, ""},

		{`{{if ContainsAny "123" "24"}}ok{{end}}`, "ok"},
		{`{{if ContainsAny "123" "45"}}ok{{end}}`, ""},

		{`{{if EqualFold "HellO" "hello"}}ok{{end}}`, "ok"},
		{`{{if EqualFold "ПривеТ" "привет"}}ok{{end}}`, "ok"},
		{`{{if EqualFold "good" "goed"}}ok{{end}}`, ""},

		{`{{if HasPrefix "hello" "he"}}ok{{end}}`, "ok"},
		{`{{if HasPrefix "hello" "lo"}}ok{{end}}`, ""},

		{`{{if HasSuffix "hello" "lo"}}ok{{end}}`, "ok"},
		{`{{if HasSuffix "hello" "he"}}ok{{end}}`, ""},

		// filters returning a string

		{`{{JSON 123}}`, "[123]"},

		{`{{Left "123456" 2}}`, "12"},
		{`{{Left "123456" 10}}`, "123456"},

		{`{{Right "123456" 2}}`, "56"},
		{`{{Right "123456" 10}}`, "123456"},

		{`{{Last (Split "12,34,56" ",")}}`, "56"},

		{`{{Mid "123456" 2}}`, "3456"},
		{`{{Mid "123456" 10}}`, ""},
		{`{{Mid "123456" 2 4}}`, "34"},
		{`{{Mid "123456" 2 10}}`, "3456"},

		{`{{Join (Split "12,34,56" ",") "|"}}`, "12|34|56"},
		{`{{Join (Split "12" ",") "|"}}`, "12"},

		{`{{Repeat "12" 3}}`, "121212"},
		{`{{Repeat "12" 0}}`, ""},
		{`{{Repeat "" 10}}`, ""},

		{`{{Replace "1234" "23" "56" -1}}`, "1564"},
		{`{{Replace "12223" "2" "5" 1}}`, "15223"},
		// ...
	}

	for _, tcase := range cases {
		tcase := tcase // important, needed for running in parallel
		t.Run(tcase.text, func(t *testing.T) {
			t.Parallel()
			result, err := executeTemplate(tcase.text, nil)
			if err != nil {
				t.Errorf("error executing template: %s", err.Error())
				t.FailNow()
			}
			if result != tcase.expected {
				t.Errorf("expected `%s` but actually `%s`", tcase.expected, result)
				t.FailNow()
			}
		})
	}
}
