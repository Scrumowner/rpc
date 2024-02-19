package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TableName(t *testing.T) {
	type args struct {
		u Tabler
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Struct method db",
			args: args{u: &GeoIntoDb{}},
			want: "geo",
		},
		{
			name: "Struct method address",
			args: args{u: &SearchIntoDb{}},
			want: "address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.args.u.TableName()
			assert.Equalf(t, result, tt.want, fmt.Sprintf("Resutl : %s but want : %s", result, tt.want))
		})
	}
}
