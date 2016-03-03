package structmap

import (
	"errors"
	"reflect"
	"strings"
)

var Tag = "structmap"

var TagPref = []string{Tag, "json"}

var (
	ErrNonStruct       = errors.New("not a struct")
	ErrNonStringKeyMap = errors.New("only supports maps with string as key")
)

func Marshal(i interface{}) (m map[string]interface{}, err error) {

	v, _ := unpackInterface(i)

	if v.Kind() != reflect.Struct {
		err = ErrNonStruct
		return
	}

	if mi, err := marshal(i); err != nil {
		return m, err
	} else {
		m, _ = mi.(map[string]interface{})
	}
	return
}

func marshal(i interface{}) (m interface{}, err error) {

	v, t := unpackInterface(i)

	var output = make(map[string]interface{})

	switch {

	case v.Kind() == reflect.Map:
		if v.Len() == 0 {
			m = output
			return
		}
		if v.MapKeys()[0].Kind() != reflect.String {
			err = ErrNonStringKeyMap
			return
		}
		for _, key := range v.MapKeys() {
			var mi interface{}
			if mi, err = marshal(v.MapIndex(key).Interface()); err != nil {
				return
			}
			output[key.String()] = mi
		}
		m = output

	case v.Kind() == reflect.Struct:
		inlineMap := make(map[string]interface{})
		inline := false

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)
			name, opts := scanTag(field)

			if name == "-" {
				continue
			}

			if name == "" {
				name = field.Name
			}

			var (
				mi  interface{}
				mim map[string]interface{}
				ok  bool
			)

			if mi, err = marshal(value.Interface()); err != nil {
				return
			}

			if mim, ok = mi.(map[string]interface{}); ok {
				if _, ok = opts["inline"]; ok {
					inline = true
					updateMap(inlineMap, mim)
				}
			} else {
				output[name] = mi
			}

		}

		if inline {
			updateMap(output, inlineMap)
		}

		m = output

	case v.Kind() == reflect.Slice || v.Kind() == reflect.Array:
		outputL := make([]interface{}, v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			var mi interface{}
			if mi, err = marshal(v.Index(i).Interface()); err != nil {
				return
			}
			outputL[i] = mi
		}
		m = outputL
	default:
		m = i
	}

	return
}

func unpackInterface(i interface{}) (v reflect.Value, t reflect.Type) {
	v = reflect.ValueOf(i)
	t = v.Type()

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
		t = v.Type()
	}

	if v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
		t = v.Type()
	}

	return
}

func updateMap(dst map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}

func scanTag(f reflect.StructField) (name string, opts map[string]interface{}) {
	tags := f.Tag
	opts = make(map[string]interface{})

	for _, tagName := range TagPref {
		_name, _opts := parseTag(tags.Get(tagName))
		if name == "" && _name != "" {
			name = _name
		}
		updateMap(opts, _opts)
	}

	return
}

func parseTag(tag string) (name string, opts map[string]interface{}) {
	opts = make(map[string]interface{})
	if tag == "" {
		return
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, part := range parts[1:] {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 1 {
			opts[kv[0]] = true
		}
		if len(kv) == 2 {
			opts[kv[0]] = kv[1]
		}
	}
	return
}
