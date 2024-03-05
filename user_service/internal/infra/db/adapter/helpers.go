package adapter

import (
	"reflect"
	"user/internal/models"
)

func GetFieldsAndPointers(user models.Userer) *FieldAndPointers {
	val := reflect.ValueOf(user).Elem()
	filds := &FieldAndPointers{
		Fields: make(map[string]interface{}),
	}
	for i := 0; i < val.NumField(); i++ {
		t := val.Type().Field(i)
		fildsName := t.Tag.Get("db")
		filds.FieldsName = append(filds.FieldsName, fildsName)
		if val.Field(i).Addr().Interface() == nil {
			continue
		}
		filds.Fields[fildsName] = val.Field(i).Addr().Interface()
	}
	return filds
}

type FieldAndPointers struct {
	FieldsName []string
	Fields     map[string]interface{}
}
