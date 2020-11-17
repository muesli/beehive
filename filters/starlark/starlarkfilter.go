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

func (filter *StarlarkFilter) convert(reflected reflect.Value) (starlark.Value, error) {
	var err error
	switch reflected.Kind() {
	case reflect.Bool:
		return starlark.Bool(reflected.Bool()), nil
	case reflect.String:
		return starlark.String(reflected.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return starlark.MakeInt64(reflected.Int()), nil
	case reflect.Slice:
		count := reflected.Len()
		vals := make([]starlark.Value, count)
		for i := 0; i < count; i++ {
			el := reflected.Index(i)
			vals[i], err = filter.convert(el)
			if err != nil {
				return nil, err
			}
		}
		return starlark.NewList(vals), nil
	case reflect.Map:
		result := starlark.NewDict(reflected.Len())
		iter := reflected.MapRange()
		for iter.Next() {
			key, err := filter.convert(iter.Key())
			if err != nil {
				return nil, err
			}
			value, err := filter.convert(iter.Value())
			if err != nil {
				return nil, err
			}
			err = result.SetKey(key, value)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case reflect.Ptr:
		return filter.convert(reflected.Elem())
	}
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("unknown type: %s", reflected.Kind())
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
		converted, err := filter.convert(reflect.ValueOf(value))
		if err != nil {
			panic(err)
		}
		arg := starlark.Tuple([]starlark.Value{starlark.String(key), converted})
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
