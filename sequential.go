package easyslice

import (
	"reflect"
	"sync"
)

func sCollectToList(s *easySlice, o interface{}) {
	reflectedSlice := reflect.ValueOf(o)
	reflectedSlice.Elem().Set(reflect.MakeSlice(reflectedSlice.Elem().Type(), 0, 0))
	actualIndex := 0
	for i := 0; i < s.collection.Len(); i++ {
		if f, v := s.evaluate(i); f {
			reflectedSlice.Elem().Set(reflect.Append(reflectedSlice.Elem(), reflect.ValueOf(v)))
			actualIndex++
		}
	}
}

func sForEach(s *easySlice, consumer TConsumer) {
	for i := 0; i < s.collection.Len(); i++ {
		if f, v := s.evaluate(i); f {
			consumer(v)
		}
	}
}

func sFindFirst(s *easySlice) IOptional {
	for i := 0; i < s.collection.Len(); i++ {
		if f, v := s.evaluate(i); f {
			return &Optional{v, true}
		}
	}
	return &Optional{nil, false}
}

func sFindAny(s *easySlice) IOptional {
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
				if f, v := s.evaluate(i); f {
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
				if f, _ := s.evaluate(i); !f {
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
	return s.FindFirst().IsPresent()
}
