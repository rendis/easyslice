package easyslice

type TPredicate func(interface{}) bool

type TMapper func(interface{}) interface{}

type TConsumer func(interface{})

type TCollection interface{}

type link interface{}

type fnLink func(v interface{}) (bool, interface{})

type linkType int8

const (
	linkPredicate linkType = iota
	linkMapper
)
