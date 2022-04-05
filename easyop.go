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
	return &easySlice{collection: reflect.ValueOf(collection), links: make([]linkInfo, 0), parallelProcessing: false}
}

func ParallelEasyOf(collection TCollection) IRootEasySlice {
	return &easySlice{collection: reflect.ValueOf(collection), links: make([]linkInfo, 0), parallelProcessing: false}
}

func (s *easySlice) Filter(predicate TPredicate) IExtendedEasySlice {
	s.links = append(s.links, linkInfo{lType: linkPredicate, operation: predicate})
	return s
}

func (s *easySlice) Map(mapper TMapper) ISimpleEasySlice {
	s.links = append(s.links, linkInfo{lType: linkMapper, operation: mapper})
	return s
}

func (s *easySlice) CollectToList(o interface{}) {
	// TODO: o validations
	sCollectToList(s, o)
}

func (s *easySlice) ForEach(consumer TConsumer) {
	sForEach(s, consumer)
}

func (s *easySlice) FindFirst() IOptional {
	return sFindFirst(s)
}

func (s *easySlice) FindAny() IOptional {
	return sFindAny(s)
}

func (s *easySlice) AllMatch() bool {
	return sAllMatch(s)
}

func (s *easySlice) AnyMatch() bool {
	return sAnyMatch(s)
}
