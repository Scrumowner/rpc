package db

import (
	"reflect"
	"strings"
)

type Scanner struct {
}

type FieldsTypes struct {
	Names []string
	Types map[string][]string
}

func NewScanner() *Scanner {
	return &Scanner{}
}

// Want pointer !!!
func (s *Scanner) GetFieldsTypes(u interface{}) FieldsTypes {
	fields := FieldsTypes{
		Types: make(map[string][]string),
	}
	val := reflect.ValueOf(u).Elem()
	for i := 0; i < val.NumField(); i++ {
		column := val.Type().Field(i).Tag.Get("db")
		fields.Names = append(fields.Names, column)
		typesRaw := val.Type().Field(i).Tag.Get("db_type")
		tags := strings.Split(typesRaw, ",")
		for j := 0; j < len(tags); j++ {
			fields.Types[column] = append(fields.Types[column], tags[j])
		}
	}
	return fields
}
