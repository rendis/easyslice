package easyslice

import (
	"reflect"
	"sync"
)

type Easy struct {
	collection         reflect.Value
	links              []linkInfo
	parallelProcessing bool
}

type linkInfo struct {
	lType     linkType
	operation link
}

func EasyOf(collection TCollection) IEasy {
	return &Easy{collection: reflect.ValueOf(collection), links: make([]linkInfo, 0), parallelProcessing: false}
}

func (s *Easy) Filter(predicate TPredicate) IExtendedEasy {
	s.links = append(s.links, linkInfo{lType: linkPredicate, operation: predicate})
	return s
}

func (s *Easy) Map(mapper TMapper) IEasy {
	s.links = append(s.links, linkInfo{lType: linkMapper, operation: mapper})
	return s
}

func (s *Easy) CollectToList(o interface{}) {
	fn := chainLinks(s.links)
	r := reflect.ValueOf(o)
	r.Elem().Set(reflect.MakeSlice(r.Elem().Type(), 0, 0))
	for i := 0; i < s.collection.Len(); i++ {
		f, v := fn(s.collection.Index(i).Interface())
		if f {
			r.Elem().Set(reflect.Append(r.Elem(), reflect.ValueOf(v)))
		}
	}
}

func (s *Easy) ForEach(consumer TConsumer) {
	fn := chainLinks(s.links)
	for i := 0; i < s.collection.Len(); i++ {
		f, v := fn(s.collection.Index(i).Interface())
		if f {
			consumer(v)
		}
	}
}

func (s *Easy) FindFirst() IOptional {
	fn := chainLinks(s.links)
	for i := 0; i < s.collection.Len(); i++ {
		f, v := fn(s.collection.Index(i).Interface())
		if f {
			return Optional{v, true}
		}
	}
	return Optional{nil, false}
}

func (s *Easy) FindAny() IOptional {
	fn := chainLinks(s.links)
	workersCh := make(chan struct{}, getNumWorkers())
	successCh := make(chan interface{}, 1)
	defer close(workersCh)
	defer close(successCh)
	var wg sync.WaitGroup
	for i := 0; i < s.collection.Len(); i++ {
		select {
		case <-successCh:
			wg.Wait()
			return Optional{successCh, true}
		case workersCh <- struct{}{}:
			wg.Add(1)
			go func(i int) {
				defer func() {
					<-workersCh
					wg.Done()
				}()
				f, v := fn(s.collection.Index(i).Interface())
				if f {
					select {
					case successCh <- v:
					default:
					}
				}
			}(i)
		}
	}
	wg.Wait()
	select {
	case f := <-successCh:
		return &Optional{f, true}
	default:
	}
	return &Optional{nil, false}
}

func (s *Easy) AllMatch() bool {
	fn := chainLinks(s.links)
	workersCh := make(chan struct{}, getNumWorkers())
	failCh := make(chan struct{}, 1)
	defer close(workersCh)
	defer close(failCh)
	var wg sync.WaitGroup
	for i := 0; i < s.collection.Len(); i++ {
		select {
		case <-failCh:
			wg.Wait()
			return false
		case workersCh <- struct{}{}:
			wg.Add(1)
			go func(i int) {
				defer func() {
					<-workersCh
					wg.Done()
				}()
				f, _ := fn(s.collection.Index(i).Interface())
				if !f {
					select {
					case failCh <- struct{}{}:
					default:
					}
				}
			}(i)
		}
	}
	wg.Wait()
	select {
	case <-failCh:
		return false
	default:
	}
	return true
}

func (s *Easy) AnyMatch() bool {
	fn := chainLinks(s.links)
	workersCh := make(chan struct{}, getNumWorkers())
	successCh := make(chan struct{}, 1)
	defer close(workersCh)
	defer close(successCh)
	var wg sync.WaitGroup
	for i := 0; i < s.collection.Len(); i++ {
		select {
		case <-successCh:
			wg.Wait()
			return true
		case workersCh <- struct{}{}:
			wg.Add(1)
			go func(i int) {
				defer func() {
					<-workersCh
					wg.Done()
				}()
				f, _ := fn(s.collection.Index(i).Interface())
				if f {
					select {
					case successCh <- struct{}{}:
					default:
					}
				}
			}(i)
		}
	}
	wg.Wait()
	select {
	case <-successCh:
		return true
	default:
	}
	return false
}

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
