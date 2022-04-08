package easyslice

import (
	"reflect"
	"sync"
)

func pCollectToSlice(s *easySlice, o interface{}) {
	var wg sync.WaitGroup
	wg.Wait()

	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}

	var lock sync.Mutex
	reflectedSlice := reflect.ValueOf(o)
	reflectedSlice.Elem().Set(reflect.MakeSlice(reflectedSlice.Elem().Type(), 0, 0))
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		go func(start, end int) {
			defer wg.Done()
			var response []reflect.Value
			for j := start; j < end; j++ {
				f, v := s.evaluate(j)
				if f {
					response = append(response, reflect.ValueOf(v))
				}
			}
			lock.Lock()
			reflectedSlice.Elem().Set(reflect.Append(reflectedSlice.Elem(), response...))
			lock.Unlock()
		}(start, end)
	}
}

func pForEach(s *easySlice, consumer TConsumer) {
	var wg sync.WaitGroup
	defer wg.Wait()

	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}

	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				f, v := s.evaluate(j)
				if f {
					consumer(v)
				}
			}
		}(start, end)
	}
}

func pAllMatch(s *easySlice) bool {
	var wg sync.WaitGroup
	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}
	closer := NewCloser()
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				select {
				case <-closer.Done():
					return
				default:
					if f, _ := s.evaluate(j); !f {
						closer.Close()
						return
					}
				}
			}
		}(start, end)
	}
	wg.Wait()
	return !closer.Status()
}

func pAnyMatch(s *easySlice) bool {
	var wg sync.WaitGroup
	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}

	closer := NewCloser()
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				select {
				case <-closer.Done():
					return
				default:
					if f, _ := s.evaluate(j); f {
						closer.Close()
						return
					}
				}
			}
		}(start, end)
	}
	wg.Wait()
	return closer.Status()
}

func pFindAny(s *easySlice) IOptional {
	var wg sync.WaitGroup
	var foundChan = make(chan interface{}, 1)
	defer func() {
		close(foundChan)
	}()

	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}

	closer := NewCloser()
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				select {
				case <-closer.Done():
					return
				default:
					if f, v := s.evaluate(j); f {
						select {
						case foundChan <- v:
							closer.Close()
						default:
						}
						return
					}
				}
			}
		}(start, end)
	}
	wg.Wait()
	select {
	case found := <-foundChan:
		return &Optional{found: true, value: found}
	default:
		return &Optional{found: false}
	}
}
