package easyslice

import (
	"reflect"
)

type EasySlice struct {
	collection         reflect.Value
	links              []linkInfo
	parallelProcessing bool
}

type linkInfo struct {
	lType     linkType
	operation link
}

func EasyOf(collection TCollection) IEasySlice {
	return &EasySlice{collection: reflect.ValueOf(collection), links: make([]linkInfo, 0), parallelProcessing: false}
}

func EasyOfParallel(collection TCollection) IEasySlice {
	return &EasySlice{collection: reflect.ValueOf(collection), links: make([]linkInfo, 0), parallelProcessing: true}
}

func (s *EasySlice) Filter(predicate TPredicate) IExtendedEasySlice {
	s.links = append(s.links, linkInfo{lType: linkPredicate, operation: predicate})
	return s
}

func (s *EasySlice) Map(mapper TMapper) IEasySlice {
	s.links = append(s.links, linkInfo{lType: linkMapper, operation: mapper})
	return s
}

func (s *EasySlice) CollectToList(o interface{}) {
	// TODO: o validations
	sCollectToList(s, o)
}

func (s *EasySlice) ForEach(consumer TConsumer) {
	sForEach(s, consumer)
}

func (s *EasySlice) FindFirst() IOptional {
	return sFindFirst(s)
}

func (s *EasySlice) FindAny() IOptional {
	return sFindAny(s)
}

func (s *EasySlice) AllMatch() bool {
	return sAllMatch(s)
}

func (s *EasySlice) AnyMatch() bool {
	return sAnyMatch(s)
}
