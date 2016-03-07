package structmap

import "reflect"

// Unmarshal a map to interface i. i should be a pointer to a struct.
func Unmarshal(m map[string]interface{}, i interface{}) (err error) {

	_, _, err = inspectInterface(i)
	if err != nil {
		return
	}

	err = unmarshal(m, i)
	return

}

func unmarshal(j interface{}, i interface{}) (err error) {

	v, t := unpackInterface(j)
	vi, _, err := inspectInterface(i)
	if err != nil {
		return
	}

	switch {
	case v.Kind() == reflect.Map:
		if v.Len() == 0 {
			return
		}
		if v.MapKeys()[0].Kind() != reflect.String {
			err = ErrNonStringKeyMap
			return
		}

		fields := make(map[string]reflect.Value)
		var inlineMap map[string]interface{}

		for i := 0; i < t.NumField(); i++ {

			field := t.Field(i)
			value := v.Field(i)
			value, _ = unpackInterface(value.Interface())
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
		for _, key := range v.MapKeys() {
			name := key.Interface().(string)
			if _, ok := fields[name]; ok {
				fields[name].Set(v.MapIndex(key))
			} else if inlineMap != nil {
				inlineMap[name] = v.MapIndex(key).Interface
			}
		}
	default:
		vi.Set(v)
	}

	return
}
