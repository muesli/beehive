// Package starlarkfilter provides a starlark-based filter.
// https://github.com/google/starlark-go
package starlarkfilter

import (
	"fmt"
	"reflect"

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

func (filter *StarlarkFilter) convert(reflected reflect.Value) starlark.Value {
	switch reflected.Kind() {
	case reflect.Bool:
		return starlark.Bool(reflected.Bool())
	case reflect.String:
		return starlark.String(reflected.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return starlark.MakeInt64(reflected.Int())
	case reflect.Slice:
		count := reflected.Len()
		vals := make([]starlark.Value, count)
		for i := 0; i < count; i++ {
			el := reflected.Index(i)
			vals[i] = filter.convert(el)
		}
		return starlark.NewList(vals)
	default:
		panic(fmt.Sprintf("unknown type: %s", reflected.Kind()))
	}
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
			filter.convert(reflect.ValueOf(value)),
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
