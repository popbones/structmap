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

// unpack an interface and get the underlying value and type. Enters infinite loop if i is a pointer to itself.
func unpackInterface(i interface{}) (v reflect.Value, t reflect.Type) {
	v = reflect.ValueOf(i)

	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	t = v.Type()

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
