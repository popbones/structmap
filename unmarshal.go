package structmap

import (
	"fmt"
	"reflect"
)

// Unmarshal a map to interface i. i should be a pointer to a struct.
func Unmarshal(m map[string]interface{}, i interface{}) (err error) {

	err = unmarshal(m, i)
	return

}

func unmarshal(j interface{}, i interface{}) (err error) {
	_ = "breakpoint"
	vj := reflect.ValueOf(j)
	tj := vj.Type()

	if vj.Kind() == reflect.Ptr {
		vj = vj.Elem()
		tj = vj.Type()
	}

	if vj.Kind() == reflect.Interface {
		vj = vj.Elem()
		tj = vj.Type()
	}

	vi := reflect.ValueOf(i)
	ti := vi.Type()

	if vi.Kind() != reflect.Ptr {
		err = ErrNotPtr
		return
	}

	vi = vi.Elem()
	ti = vi.Type()

	if vi.Kind() == reflect.Interface {
		vi = vi.Elem()
		ti = vi.Type()
	}

	fmt.Println(tj.String())
	_ = "breakpoint"
	switch {
	case vj.Kind() == reflect.Map:

		// if vj is a map, vi should be a struct or a map of the same key type
		if vj.Len() == 0 {
			return
		}
		if vj.MapKeys()[0].Kind() != reflect.String {
			err = ErrNonStringKeyMap
			return
		}

		if vi.Kind() == reflect.Struct {
			fields := make(map[string]reflect.Value)
			var inlineMap map[string]interface{}

			for i := 0; i < ti.NumField(); i++ {

				field := ti.Field(i)
				value := vi.Field(i)

				name, opts := scanTag(field)

				if name == "-" {
					continue
				}

				if name == "" {
					name = field.Name
				}

				if _, ok := opts["inline"]; ok {
					inlineMap, _ = value.Interface().(map[string]interface{})
				} else {
					fields[name] = value
				}
			}
			_ = "breakpoint"
			for _, key := range vj.MapKeys() {
				name := key.Interface().(string)
				if _, ok := fields[name]; ok {
					value := vj.MapIndex(key)
					if value.Kind() == reflect.Interface {
						value = value.Elem()
					}
					fields[name].Set(value)
				} else if inlineMap != nil {
					inlineMap[name] = vj.MapIndex(key).Interface()
				}
			}
		}

	default:
		vi.Set(vj)
	}

	return
}
