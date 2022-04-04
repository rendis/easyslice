package easyslice

import "reflect"

func chainLinks(links []linkInfo) fnLink {
	lastLink := true
	var currentFn fnLink
	for i := len(links) - 1; i >= 0; i-- {
		l := links[i]
		switch l.lType {
		case linkPredicate:
			if lastLink {
				lastLink = false
				currentFn = buildLastPredicate(l.operation.(TPredicate))
			} else {
				currentFn = buildPredicate(l.operation.(TPredicate), currentFn)
			}
		case linkMapper:
			if lastLink {
				lastLink = false
				currentFn = buildLastMapper(l.operation.(TMapper))
			} else {
				currentFn = buildMapper(l.operation.(TMapper), currentFn)
			}
		}
	}
	return currentFn
}

func buildLastPredicate(predicate TPredicate) fnLink {
	return func(v interface{}) (bool, interface{}) {
		return predicate(v), v
	}
}

func buildPredicate(predicate TPredicate, fn fnLink) fnLink {
	return func(v interface{}) (bool, interface{}) {
		if !predicate(v) {
			return false, reflect.Value{}
		}
		return fn(v)
	}
}

func buildLastMapper(mapper TMapper) fnLink {
	return func(v interface{}) (bool, interface{}) {
		return true, mapper(v)
	}
}

func buildMapper(mapper TMapper, fn fnLink) fnLink {
	return func(v interface{}) (bool, interface{}) {
		nv := mapper(v)
		return fn(nv)
	}
}
