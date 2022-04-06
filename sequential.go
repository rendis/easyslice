package easyslice

import (
	"reflect"
)

func sCollectToSlice(s *easySlice, o interface{}) {
	reflectedSlice := reflect.ValueOf(o)
	reflectedSlice.Elem().Set(reflect.MakeSlice(reflectedSlice.Elem().Type(), 0, 0))
	actualIndex := 0
	for i := 0; i < s.collection.Len(); i++ {
		if f, v := s.evaluate(i); f {
			reflectedSlice.Elem().Set(reflect.Append(reflectedSlice.Elem(), reflect.ValueOf(v)))
			actualIndex++
		}
	}
}

func sForEach(s *easySlice, consumer TConsumer) {
	for i := 0; i < s.collection.Len(); i++ {
		if f, v := s.evaluate(i); f {
			consumer(v)
		}
	}
}

func sFindFirst(s *easySlice) IOptional {
	for i := 0; i < s.collection.Len(); i++ {
		if f, v := s.evaluate(i); f {
			return &Optional{v, true}
		}
	}
	return &Optional{nil, false}
}
