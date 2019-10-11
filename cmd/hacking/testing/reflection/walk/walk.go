package walk

import "reflect"

func walk(obj interface{}, fn func(input string)) {

	valueOf := getValue(obj)

	var numOfElements int
	var getValue func(int) reflect.Value

	switch valueOf.Kind() {
	case reflect.Struct:
		// Inspect the field
		numOfElements = valueOf.NumField()
		getValue = valueOf.Field
	case reflect.Slice, reflect.Array:
		numOfElements = valueOf.Len()
		getValue = valueOf.Index
	case reflect.String:
		fn(valueOf.String())
	case reflect.Map:
		for _, key := range valueOf.MapKeys() {
			walk(valueOf.MapIndex(key).Interface(), fn)
		}
	}

	for i := 0; i < numOfElements; i++ {
		walk(getValue(i).Interface(), fn)
	}
}

func getValue(val interface{}) reflect.Value {
	valueOf := reflect.ValueOf(val)

	if valueOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
	}
	return valueOf
}