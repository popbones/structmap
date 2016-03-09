package structmap

import (
	"fmt"
	"reflect"
)

// decode a map to a struct
func Unmarshal(src map[string]interface{}, dst interface{}) (err error) {
	return decode(src, dst)
}

func decode(src interface{}, dst interface{}) (err error) {
	var (
		srcValue, dstValue reflect.Value
		srcType, dstType   reflect.Type
		srcKind, dstKind   reflect.Kind
	)
	srcValue = reflect.ValueOf(src)
	srcKind = srcValue.Kind()
	if srcValue.IsValid() {
		srcType = srcValue.Type()
	}
	dstValue = reflect.ValueOf(dst)
	dstKind = dstValue.Kind()
	if dstValue.IsValid() {
		dstType = dstValue.Type()
	}

	switch dstKind {
	case reflect.Ptr, reflect.Interface:
		if dstValue.Elem().Kind() == reflect.Ptr || dstValue.Elem().Kind() == reflect.Interface {
			if !dstValue.Elem().IsNil() {
				return decode(src, dstValue.Elem().Addr().Interface())
			}
		}
		dstValue = dstValue.Elem()
		dstType = dstValue.Type()
		dstKind = dstValue.Kind()
	}

	switch srcKind {
	case reflect.Map:
		keyKind := srcType.Key().Kind()

		if keyKind == reflect.String {
			srcMap := src.(map[string]interface{})

			switch dstKind {
			case reflect.Map:
				if dstType.Key().Kind() == reflect.String {
					dstMap := dst.(map[string]interface{})
					for key, value := range srcMap {
						dstMap[key] = value
					}
				}
			case reflect.Struct:
				structMap := make(map[string]interface{})
				var inlineStruct interface{}
				var hasInlineMap bool
				var hasInlineStruct bool
				inlineMap := make(map[string]interface{})

				hasInline := false
				for i := 0; i < dstType.NumField(); i++ {
					field := dstType.Field(i)
					fVal := dstValue.Field(i)
					name, opts := scanTag(field)
					if name == "-" {
						continue
					}
					isInline := false
					_, isInline = opts["inline"]

					res := fVal

					if isInline {

						if fVal.Type() == reflect.TypeOf(inlineMap) {
							hasInlineMap = true

							if isInline {
								if !hasInline {
									hasInline = !hasInline
									res.Set(reflect.ValueOf(inlineMap))
								} else {
									fmt.Println("warning, structmap detected multiple inline tags on", dstType.String())
								}

							}
						} else if fVal.Kind() == reflect.Struct {
							inlineStruct = res.Addr().Interface()
						}

					} else {
						structMap[name] = res.Addr().Interface()
					}
				}

				for key, value := range srcMap {
					if _, ok := structMap[key]; ok {
						if err = decode(value, structMap[key]); err != nil {
							return
						}
					} else if hasInlineStruct {
						return decode(map[string]interface{}{key: value}, &inlineStruct)
					} else if hasInlineMap {
						var output interface{}
						if err = decode(value, &output); err != nil {
							return
						}
						inlineMap[key] = output
					}
				}
			}

		} else {
			dstValue.Set(srcValue)
		}
	case reflect.Slice, reflect.Array:
		outputSlice := make([]interface{}, dstValue.Len(), dstValue.Len())
		if inputSlice, ok := srcValue.Interface().([]interface{}); ok {
			for i, value := range inputSlice {
				outputSlice[i] = value
			}
		}
		dstValue.Set(reflect.ValueOf(outputSlice))
	default:
		dstValue.Set(srcValue)
	}

	return
}
