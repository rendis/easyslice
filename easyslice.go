package easyslice

type IRootEasySlice interface {
	IMap
	IFilter
}

type ISimpleEasySlice interface {
	IRootEasySlice
	IConsumerEasy
	ICollectors
}

type IExtendedEasySlice interface {
	ISimpleEasySlice
	IDependentsClosers
}

type IDependentsClosers interface {
	IFindAny
	IFindFirst
	IAllMatch
	IAnyMatch
}

type IConsumerEasy interface {
	IForEach
}

type IMap interface {
	Map(TMapper) ISimpleEasySlice
}

type IFilter interface {
	Filter(TPredicate) IExtendedEasySlice
}

type IFindAny interface {
	FindAny() IOptional
}

type IFindFirst interface {
	FindFirst() IOptional
}

type IAllMatch interface {
	AllMatch() bool
}

type IAnyMatch interface {
	AnyMatch() bool
}

type IForEach interface {
	ForEach(TConsumer)
}

type ICollectors interface {
	CollectToSlice(TPtrSlice)
}
