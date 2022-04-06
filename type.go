package easyslice

type TPredicate func(interface{}) bool

type TMapper func(interface{}) interface{}

type TConsumer func(interface{})

type TCollection interface{}

type TPtrSlice interface{}

type link interface{}

type linkType int8

const (
	linkPredicate linkType = iota
	linkMapper
)
