package walk_tags

import "reflect"

func walk(obj interface{}, fn func(key string, val interface{})) {
	val :=  reflect.ValueOf(obj)
	valType := val.Type()

	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			// Check the type of the struct field
			field := val.Field(i)
			switch field.Kind() {
			case reflect.Struct:
				walk(field.Interface(), fn)
			default:
				tag := valType.Field(i).Tag
				tagVal := tag.Get("whome")
				fieldVal := val.Field(i).Interface()
				fn(tagVal, fieldVal)
			}
		}
	case reflect.Ptr:
		val = val.Elem()
		walk(val.Interface(), fn)
	}
}
