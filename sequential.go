package easyslice

import (
	"reflect"
	"sync"
)

func sCollectToList(s *easySlice, o interface{}) {
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

func sForEach(s *easySlice, consumer TConsumer) {
	fn := chainLinks(s.links)
	for i := 0; i < s.collection.Len(); i++ {
		f, v := fn(s.collection.Index(i).Interface())
		if f {
			consumer(v)
		}
	}
}

func sFindFirst(s *easySlice) IOptional {
	fn := chainLinks(s.links)
	for i := 0; i < s.collection.Len(); i++ {
		f, v := fn(s.collection.Index(i).Interface())
		if f {
			return &Optional{v, true}
		}
	}
	return &Optional{nil, false}
}

func sFindAny(s *easySlice) IOptional {
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
			return &Optional{successCh, true}
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

func sAllMatch(s *easySlice) bool {
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

func sAnyMatch(s *easySlice) bool {
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
