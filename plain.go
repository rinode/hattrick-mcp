package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/lucianoq/hattrick/chpp"
)

var (
	moneyType = reflect.TypeFor[chpp.Money]()
	tsiType   = reflect.TypeFor[chpp.TSI]()
)

// plain converts CHPP values into JSON-friendly data for a language model:
// enums become "name (n)", times become strings, and zero values are dropped.
// chpp types carry only XML tags, and HattrickTime has no JSON encoding.
func plain(v any) any {
	return convert(reflect.ValueOf(v))
}

func convert(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}
		return convert(v.Elem())
	}
	// Keep booleans even when false: a false flag is real information, unlike
	// an empty string or a zero number.
	if v.Kind() == reflect.Bool {
		return v.Bool()
	}
	if v.IsZero() {
		return nil
	}
	if v.Type() == moneyType || v.Type() == tsiType {
		return v.Interface()
	}
	if s, ok := v.Interface().(fmt.Stringer); ok {
		return stringer(v, s.String())
	}

	switch v.Kind() {
	case reflect.Struct:
		m := map[string]any{}
		for i := range v.NumField() {
			f := v.Type().Field(i)
			if !f.IsExported() {
				continue
			}
			if x := convert(v.Field(i)); x != nil {
				m[f.Name] = x
			}
		}
		if len(m) == 0 {
			return nil
		}
		return m
	case reflect.Slice, reflect.Array:
		l := make([]any, 0, v.Len())
		for i := range v.Len() {
			l = append(l, convert(v.Index(i)))
		}
		return l
	case reflect.Map:
		m := map[string]any{}
		for _, k := range v.MapKeys() {
			m[fmt.Sprint(k.Interface())] = convert(v.MapIndex(k))
		}
		return m
	}
	return v.Interface()
}

// stringer keeps the number of an integer enum next to its name; IDs stay numbers.
func stringer(v reflect.Value, s string) any {
	var n string
	switch {
	case v.CanInt():
		n = strconv.FormatInt(v.Int(), 10)
	case v.CanUint():
		n = strconv.FormatUint(v.Uint(), 10)
	default:
		return s
	}
	if s == "" || s == n {
		return v.Interface()
	}
	return s + " (" + n + ")"
}
