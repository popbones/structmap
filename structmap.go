package structmap

import (
	"errors"
	"reflect"
	"strings"
)

var Tag = "map"

var TagPref = []string{Tag}

var (
	ErrNonStruct       = errors.New("not a struct")
	ErrNonStringKeyMap = errors.New("only supports maps with string as key")
	ErrNotPtr          = errors.New("not a pointer")
)

// check and unpack a interface for before marshalling
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

// check and inspect a interface before unmarshalling
func inspectInterface(i interface{}) (v reflect.Value, t reflect.Type, err error) {
	v = reflect.ValueOf(i)
	t = v.Type()

	if v.Kind() != reflect.Ptr {
		err = ErrNotPtr
		return
	}

	v = v.Elem()
	t = v.Type()

	if v.Kind() == reflect.Interface {
		v = v.Elem()
		t = v.Type()
	}

	if v.Kind() != reflect.Struct {
		err = ErrNonStruct
		return
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
