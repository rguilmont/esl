package query

import (
	"reflect"
	"testing"
)

func TestGenerateEsQuery(t *testing.T) {
	type args struct {
		queryString string
		tsFrom      string
		tsAfter     string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateEsQuery(tt.args.queryString, tt.args.tsFrom, tt.args.tsAfter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateEsQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
