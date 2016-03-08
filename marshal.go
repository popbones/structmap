package structmap

import "reflect"

// Marshal a interface object to a map. The interace object should be a struct or a pointer to a struct otherwise it will return ErrNotStruct.
// i should not contain maps other than map[string]interface{} maps, otherwise it will return ErrNonStringKeyMap
//
// By default, you can use 'map' tag to specifiy how struct field should be marsheld, in the following syntax
// map:"<key_name>[,inline]"
//
// Inline here means everthing in this field should be flated to the parent map. Marshal doesn't complain if more than one inline fields are specified, but generally it should be avoided as it could cause unexpected output.
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
		if t.Key().Kind() != reflect.String {
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
