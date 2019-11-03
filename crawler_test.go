package main

import (
	"reflect"
	"testing"
)

func Test_extract(t *testing.T) {
	type args struct {
		domain string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extract(tt.args.domain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extract() = %v, want %v", got, tt.want)
			}
		})
	}
}
