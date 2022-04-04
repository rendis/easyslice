package easyslice

type IEasy interface {
	IMap
	IFilter
	IConsumerEasy
	ICollectors
}

type IExtendedEasy interface {
	IEasy
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
	Map(TMapper) IEasy
}

type IFilter interface {
	Filter(TPredicate) IExtendedEasy
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
