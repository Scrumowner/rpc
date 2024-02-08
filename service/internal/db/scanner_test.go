package db

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestScanner_GetFieldsTypes(t *testing.T) {
	type args struct {
		u interface{}
	}
	tests := []struct {
		name string
		args args
		want FieldsTypes
	}{
		{
			name: "Valid #1",
			args: args{
				u: &struct {
					name string `db:"name" db_type:"text"`
					lat  string `db:"lat" db_type:"text"`
					lon  string `db:"lon" db_type:"text"`
					id   int    `db:"id" db_type:"integer"`
				}{},
			},
			want: FieldsTypes{
				Names: []string{"name", "lat", "lon", "id"},
				Types: map[string][]string{
					"name": []string{"text"},
					"lat":  []string{"text"},
					"lon":  []string{"text"},
					"id":   []string{"integer"},
				},
			},
		},
		{
			name: "Valid #2",
			args: args{
				u: &struct {
					name   string `db:"name" db_type:"text"`
					lat    string `db:"lat" db_type:"text"`
					lon    string `db:"lon" db_type:"text"`
					id     int    `db:"id" db_type:"integer"`
					result string `db:"result" db_type:"text,varchar(55)"`
				}{},
			},
			want: FieldsTypes{
				Names: []string{"name", "lat", "lon", "id", "result"},
				Types: map[string][]string{
					"name":   []string{"text"},
					"lat":    []string{"text"},
					"lon":    []string{"text"},
					"id":     []string{"integer"},
					"result": []string{"text", "varchar(55)"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{}
			got := s.GetFieldsTypes(tt.args.u)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldsTypes() = %v, want %v", got, tt.want)
			}
			assert.NotNil(t, got, "GetFieldsTypes() result shoud not be nil")
		})
	}
}
