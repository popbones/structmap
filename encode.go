package structmap

import (
	"fmt"
	"reflect"
)

func Marshal(i interface{}) (m map[string]interface{}, err error) {
	v := reflect.ValueOf(i)
	k := v.Kind()

	for k == reflect.Ptr || k == reflect.Interface {
		v = v.Elem()
		k = v.Kind()
	}

	if k != reflect.Map && k != reflect.Struct {
		err = ErrNonStruct
		return
	}

	if encoded, err := encode(i); err != nil {
		return nil, err
	} else {
		m, _ = encoded.(map[string]interface{})
	}
	return
}

func encode(input interface{}) (output interface{}, err error) {
	inputValue := reflect.ValueOf(input)
	inputKind := inputValue.Kind()
	inputType := inputValue.Type()

	switch inputKind {
	case reflect.Map:
		// if input is a map, output should also be a map
		keyKind := inputType.Key().Kind()
		if keyKind == reflect.String {
			inputMap, _ := input.(map[string]interface{})
			outputMap := make(map[string]interface{})
			for key, value := range inputMap {
				nested, err := encode(value)
				if err != nil {
					return nil, err
				}
				outputMap[key] = nested
			}
			return outputMap, nil
		} else {
			return input, nil
		}
	case reflect.Struct:
		// if input is a struct, we try to encode it
		outputMap := make(map[string]interface{})
		hasInline := false

		for i := 0; i < inputType.NumField(); i++ {
			field := inputType.Field(i)
			fVal := inputValue.Field(i)
			name, opts := scanTag(field)
			if name == "-" {
				continue
			}
			isInline := false
			_, isInline = opts["inline"]

			res, err := encode(fVal.Interface())
			if err != nil {
				return nil, err
			}

			if nestedMap, ok := res.(map[string]interface{}); ok && isInline {
				if !hasInline {
					hasInline = !hasInline
				} else {
					fmt.Println("warning, structmap detected multiple inline tags on", inputType.String())
				}
				updateMap(outputMap, nestedMap)
				continue
			}

			outputMap[name], err = encode(fVal.Interface())
			if err != nil {
				return nil, err
			}
		}
		return outputMap, nil
	case reflect.Ptr:
		if inputValue.IsNil() {
			return
		}
		return encode(inputValue.Elem().Interface())
	case reflect.Interface:
		return encode(inputValue.Elem().Interface())
	case reflect.Slice, reflect.Array:
		outputSlice := make([]interface{}, inputValue.Len(), inputValue.Len())
		for i := 0; i < inputValue.Len(); i++ {
			res, err := encode(inputValue.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			outputSlice[i] = res
		}
		return outputSlice, nil
	default:
		return input, nil
	}

	return
}
