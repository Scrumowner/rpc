package adapter

import (
	"fmt"
	"reflect"
)

func GetStructInfo(u interface{}, args ...func(*[]reflect.StructField)) StructInfo {
	val := reflect.ValueOf(u).Elem()
	var structFields []reflect.StructField

	for i := 0; i < val.NumField(); i++ {
		structFields = append(structFields, val.Type().Field(i))
	}

	for i := range args {
		if args[i] == nil {
			continue
		}
		args[i](&structFields)
	}

	var res StructInfo

	for _, field := range structFields {
		valueField := val.FieldByName(field.Name)
		res.Pointers = append(res.Pointers, valueField.Addr().Interface())
		res.Fields = append(res.Fields, field.Tag.Get("db"))
	}

	return res
}

type Tabler interface {
	TableName() string
	TableFieldAndType() string
}
type StructInfo struct {
	Fields   []string
	Pointers []interface{}
}

type HistorySearchAddress struct {
	Query  string `json:"query" db:"query" db_type:"VARCHAR(100)"`
	Result string `json:"result" db:"result" db_type:"VARCHAR(100)"`
	GeoLat string `json:"lat" db:"lat" db_type:"VARCHAR(100)"`
	GeoLon string `json:"lon" db:"lon" db_type:"VARCHAR(100)"`
}

func NewHistorySearchAddress() Tabler {
	return &HistorySearchAddress{}
}
func (u *HistorySearchAddress) TableName() string {
	return "history_search_address"
}

func (u *HistorySearchAddress) TableFieldAndType() string {
	val := reflect.ValueOf(u).Elem()
	tp := val.Type()
	var res string
	for i := 0; i < tp.NumField(); i++ {
		if i == tp.NumField()-1 {
			res += fmt.Sprintf("%s %s", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
			return res
		}
		res += fmt.Sprintf("%s %s,", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
	}
	return res
}

type SearchRequest struct {
	Query string `json:"query" db:"query" db_type:"VARCHAR(100)"`
}

func (u *SearchRequest) TableName() string {
	return "search_history"
}

func (u *SearchRequest) TableFieldAndType() string {
	val := reflect.ValueOf(u).Elem()
	tp := val.Type()
	var res string
	for i := 0; i < tp.NumField(); i++ {
		if i == tp.NumField()-1 {
			res += fmt.Sprintf("%s %s", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
			return res
		}
		res += fmt.Sprintf("%s %s,", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
	}
	return res
}

type ReworkedSearch struct {
	Result string `json:"result" db:"result" db_type:"VARCHAR(100)"`
	GeoLat string `json:"lat" db:"lat" db_type:"VARCHAR(100)"`
	GeoLon string `json:"lon" db:"lon" db_type:"VARCHAR(100)"`
}

func NewReworkerSearch() Tabler {
	return &ReworkedSearch{}
}
func (u *ReworkedSearch) TableName() string {
	return "address"
}
func (u *ReworkedSearch) TableFieldAndType() string {
	val := reflect.ValueOf(u).Elem()
	tp := val.Type()
	var res string
	for i := 0; i < tp.NumField(); i++ {
		if i == tp.NumField()-1 {
			res += fmt.Sprintf("%s %s", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
			return res
		}
		res += fmt.Sprintf("%s %s,", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
	}
	return res
}

func NewSearchRequest() Tabler {
	return &ReworkedSearch{}
}

func FilterByFields(fields ...int) func(fields *[]reflect.StructField) {
	return func(fs *[]reflect.StructField) {
		var res []reflect.StructField
		for _, fieldIndex := range fields {
			for idx, field := range *fs {
				if idx == fieldIndex {
					res = append(res, field)
				}
			}
		}
		*fs = res
	}
}

func FilterByTags(tags map[string]func(value string) bool) func(fields *[]reflect.StructField) {
	return func(fs *[]reflect.StructField) {
		var res []reflect.StructField
		for _, field := range *fs {
			for tag, filterFunc := range tags {
				if filterFunc(field.Tag.Get(tag)) {
					res = append(res, field)
					break
				}
			}
		}
		*fs = res
	}
}

type Geo struct {
	Lat    string `json:"lat" db:"lat" db_type:"VARCHAR(100)"`
	Lng    string `json:"lng" db:"lng" db_type:"VARCHAR(100)"`
	Result string `json:"result" db:"result" db_type:"VARCHAR(100)"`
	GeoLat string `json:"lat" db:"geolat" db_type:"VARCHAR(100)"`
	GeoLon string `json:"lon" db:"geolon" db_type:"VARCHAR(100)"`
}

func NewGeo() Tabler {
	return &Geo{}
}

func (u *Geo) TableName() string { return "geo" }
func (u *Geo) TableFieldAndType() string {
	val := reflect.ValueOf(u).Elem()
	tp := val.Type()
	var res string
	for i := 0; i < tp.NumField(); i++ {
		if i == tp.NumField()-1 {
			res += fmt.Sprintf("%s %s", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
			return res
		}
		res += fmt.Sprintf("%s %s,", tp.Field(i).Tag.Get("db"), tp.Field(i).Tag.Get("db_type"))
	}
	return res
}
