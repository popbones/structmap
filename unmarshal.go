package structmap

import "reflect"

// Unmarshal a map to interface i. i should be a pointer to a struct.
func Unmarshal(m map[string]interface{}, i interface{}) (err error) {

	err = unmarshal(m, i)
	return

}

func unmarshal(j interface{}, i interface{}) (err error) {
	// _ = "breakpoint"

	// Unpack all the interfaces
	vi, ti := unpackInterface(i)
	vj, tj := unpackInterface(j)

	// _ = "breakpoint"

	switch {

	// When the input is a map we further unmarshal it
	case vj.Kind() == reflect.Map && vi.Kind() == reflect.Struct:

		// If j is a map but not a map[string]interface{}
		if tj.Key().Kind() != reflect.String {
			err = ErrNonStringKeyMap
			return
		}

		// If j is empty
		if vj.Len() == 0 {
			return
		}

		// Inspect the output sturct
		fields := make(map[string]reflect.Value)
		var inlineMap map[string]interface{}

		for i := 0; i < ti.NumField(); i++ {
			_ = "breakpoint"
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

				value, _ = unpackInterface(value.Addr().Interface())

				if value.Kind() == reflect.Map && value.Type().Key().Kind() == reflect.String {
					inlineMap = value.Interface().(map[string]interface{})
					if inlineMap == nil {
						inlineMap = make(map[string]interface{})
						value.Set(reflect.ValueOf(inlineMap))
					}
				}
			} else {
				fields[name] = value
			}
		}

		// Unmarshalling
		for _, key := range vj.MapKeys() {
			_ = "breakpoint"
			name := key.Interface().(string)

			if _, ok := fields[name]; ok {
				if err = unmarshal(vj.MapIndex(key).Interface(), fields[name].Addr().Interface()); err != nil {
					return
				}
			} else if inlineMap != nil {
				inlineMap[name] = vj.MapIndex(key).Interface()
			}
		}

	default:
		// skip
		vi.Set(vj)
	}

	return
}
