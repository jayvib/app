package walk

import "reflect"

func walk(obj interface{}, fn func(input string)) {

	valueOf := getValue(obj)

	walkValue := func(value reflect.Value){
		walk(value.Interface(), fn)
	}

	switch valueOf.Kind() {
	case reflect.String:
		fn(valueOf.String())
	case reflect.Struct:
		// Inspect the field
		for i := 0; i < valueOf.NumField(); i++ {
			walkValue(valueOf.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < valueOf.Len(); i++ {
			walkValue(valueOf.Index(i))
		}
	case reflect.Map:
		for _, key := range valueOf.MapKeys() {
			walkValue(valueOf.MapIndex(key))
		}
	}
}

func getValue(val interface{}) reflect.Value {
	valueOf := reflect.ValueOf(val)

	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	return valueOf
}