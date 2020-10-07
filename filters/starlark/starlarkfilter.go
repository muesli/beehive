// Package templatefilter provides a starlark-based filter.
// https://github.com/google/starlark-go
package starlarkfilter

import (
	"github.com/muesli/beehive/filters"
	"go.starlark.net/starlark"
)

// StarlarkFilter is a starlark-based filter.
// https://github.com/google/starlark-go
type StarlarkFilter struct {
}

// Name returns the name of this Filter.
func (filter *StarlarkFilter) Name() string {
	return "starlark"
}

// Description returns the description of this Filter.
func (filter *StarlarkFilter) Description() string {
	return "This filter passes when `main` function returns True"
}

func (filter *StarlarkFilter) convert(v interface{}) starlark.Value {
	switch val := v.(type) {
	case string:
		return starlark.String(val)
	case bool:
		return starlark.Bool(val)
	case int:
		return starlark.MakeInt(val)
	case []interface{}:
		vals := make([]starlark.Value, len(val))
		for i, el := range val {
			vals[i] = filter.convert(el)
		}
		return starlark.NewList(vals)
	}
	panic("unknown type")
}

// Passes returns true when the Filter matched the opts.
func (filter *StarlarkFilter) Passes(opts map[string]interface{}, template string) bool {
	template = dedent(template)
	thread := &starlark.Thread{Name: "main thread"}
	globals, err := starlark.ExecFile(thread, "template.star", template, nil)
	if err != nil {
		panic(err)
	}

	kwargs := make([]starlark.Tuple, 0)
	for key, value := range opts {
		arg := starlark.Tuple([]starlark.Value{
			starlark.String(key),
			filter.convert(value),
		})
		kwargs = append(kwargs, arg)
	}

	main := globals["main"]
	result, err := starlark.Call(thread, main, nil, kwargs)
	if err != nil {
		panic(err)
	}
	return bool(result.Truth())
}

func init() {
	f := StarlarkFilter{}

	filters.RegisterFilter(&f)
}
