package carrier

import (
	"fmt"
	"reflect"
)

//Transfer переводит значения из полей src в dst (связь по тэгу "tr")
func Transfer(src, dst interface{}) error {
	if src == nil || dst == nil {
		return fmt.Errorf("Value is nil")
	}
	in := reflect.ValueOf(src)
	out := reflect.ValueOf(dst)

	mapIn := make(map[string]reflect.Value)
	mapOut := make(map[string]reflect.Value)

	for i := 0; i < in.Elem().Type().NumField(); i++ {
		if key := in.Elem().Type().Field(i).Tag.Get("tr"); key != "" {
			mapIn[key] = in.Elem().Field(i).Addr()
		}
	}
	for i := 0; i < out.Elem().Type().NumField(); i++ {
		if key := out.Elem().Type().Field(i).Tag.Get("tr"); key != "" {
			mapOut[key] = out.Elem().Field(i).Addr()
		}
	}

	for i, v1 := range mapOut {
		if v2, ok := mapIn[i]; ok {
			if v1.Elem().Type().String() == v2.Elem().Type().String() {
				v1.Elem().Set(v2.Elem())
			} else {
				return fmt.Errorf("Types are not equal")
			}
		}
	}
	return nil
}
