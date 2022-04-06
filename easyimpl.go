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

func EasyOf(collection TCollection) IRootEasySlice {
	checkTCollection(collection)
	return &easySlice{collection: reflect.ValueOf(collection), links: make([]linkInfo, 0), parallelProcessing: false}
}

func ParallelEasyOf(collection TCollection) IRootEasySlice {
	checkTCollection(collection)
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
