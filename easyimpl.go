package easyslice

import (
	"reflect"
)

type easySlice struct {
	collection         reflect.Value
	links              []linkInfo
	parallelProcessing bool
}

type linkInfo struct {
	lType     linkType
	operation link
}

// EasyOf returns a new sequential EasySlice instance.
// mapToSlice
//  - false:  convert map to a slice of its values (default is false)
//  - true: convert map to slice of Pair
func EasyOf(collection TCollection, mapToPair ...bool) IRootEasySlice {
	checkTCollection(collection)
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Map {
		v = mapToCollection(v, mapToPair)
	}
	return &easySlice{collection: v, links: make([]linkInfo, 0), parallelProcessing: false}
}

// ParallelEasyOf returns a new parallel EasySlice instance.
// mapToSlice
//  - false:  convert map to a slice of its values (default is false)
//  - true: convert map to slice of Pair
func ParallelEasyOf(collection TCollection, mapToPair ...bool) IRootEasySlice {
	checkTCollection(collection)
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Map {
		v = mapToCollection(v, mapToPair)
	}
	return &easySlice{collection: reflect.ValueOf(collection), links: make([]linkInfo, 0), parallelProcessing: true}
}

func (s *easySlice) evaluate(index int) (b bool, v interface{}) {
	v = s.collection.Index(index).Interface()
	b = true
	for _, l := range s.links {
		switch l.lType {
		case linkPredicate:
			if !l.operation.(TPredicate)(v) {
				b = false
				return
			}
		case linkMapper:
			v = l.operation.(TMapper)(v)
		}
	}
	return

}

func (s *easySlice) Filter(predicate TPredicate) IExtendedEasySlice {
	s.links = append(s.links, linkInfo{lType: linkPredicate, operation: predicate})
	return s
}

func (s *easySlice) Map(mapper TMapper) ISimpleEasySlice {
	s.links = append(s.links, linkInfo{lType: linkMapper, operation: mapper})
	return s
}

func (s *easySlice) CollectToSlice(slice TPtrSlice) {
	checkTSlice(slice)
	if s.parallelProcessing {
		pCollectToSlice(s, slice)
	} else {
		sCollectToSlice(s, slice)
	}
}

func (s *easySlice) ForEach(consumer TConsumer) {
	if s.parallelProcessing {
		pForEach(s, consumer)
	} else {
		sForEach(s, consumer)
	}
}

func (s *easySlice) FindFirst() IOptional {
	return sFindFirst(s)
}

func (s *easySlice) FindAny() IOptional {
	return pFindAny(s)
}

func (s *easySlice) AllMatch() bool {
	return pAllMatch(s)
}

func (s *easySlice) AnyMatch() bool {
	return pAnyMatch(s)
}

func mapToCollection(v reflect.Value, mapToPair []bool) reflect.Value {
	var r []interface{}
	iterator := v.MapRange()
	if len(mapToPair) == 0 || !mapToPair[0] {
		for iterator.Next() {
			r = append(r, iterator.Value().Interface())
		}
	} else {
		for iterator.Next() {
			r = append(r, Pair{Key: iterator.Key().Interface(), Value: iterator.Value().Interface()})
		}
	}
	return reflect.ValueOf(r)
}
