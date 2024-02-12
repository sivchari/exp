package maps_test

import (
	"reflect"
	"sort"
	"testing"

	"sivchari.github.io/exp/maps"
)

func TestKeys(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	tests := []struct {
		name string
		args args[int, string]
		want []int
	}{
		{
			name: "test1",
			args: args[int, string]{
				m: map[int]string{
					1: "a",
					2: "b",
					3: "c",
				},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Keys(tt.args.m)
			sort.Ints(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	type args[K comparable, V any] struct {
		m map[K]V
	}
	tests := []struct {
		name string
		args args[int, string]
		want []string
	}{
		{
			name: "test1",
			args: args[int, string]{
				m: map[int]string{
					1: "a",
					2: "b",
					3: "c",
				},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maps.Values(tt.args.m)
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
