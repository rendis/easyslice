package easyslice

type IEasySlice interface {
	IMap
	IFilter
	IConsumerEasy
	ICollectors
}

type IExtendedEasySlice interface {
	IEasySlice
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
	Map(TMapper) IEasySlice
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
	CollectToList(interface{})
}
