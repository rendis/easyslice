package easyslice

import (
	"fmt"
	"reflect"
)

func checkTCollection(t TCollection) {
	r := reflect.TypeOf(t)
	if r.Kind() != reflect.Slice && r.Kind() != reflect.Array && r.Kind() != reflect.Map {
		panic(fmt.Sprintf("EasySlice: TCollection must be a [slice | array ], got '%v' (%v)", r.Kind(), t))
	}
}

func checkTSlice(s TPtrSlice) {
	r := reflect.TypeOf(s)
	if r.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("EasySlice: TPtrSlice must be a ptr to slice, got '%v'", r.Kind()))
	}

	if r.Elem().Kind() != reflect.Slice {
		panic(fmt.Sprintf("EasySlice: TPtrSlice must be a ptr to slice ('ptr' -> 'slice'), got '%v' -> '%v'", r.Kind(), r.Elem().Kind()))
	}
}
