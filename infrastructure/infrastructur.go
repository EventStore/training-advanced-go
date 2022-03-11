package infrastructure

import "reflect"

func GetValueType(t interface{}) reflect.Type {
	v := reflect.ValueOf(t)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Type()
}
