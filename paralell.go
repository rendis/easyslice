package easyslice

import (
	"reflect"
	"sync"
	"sync/atomic"
)

func pCollectToList(s *easySlice, o interface{}) {
	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	var wg sync.WaitGroup
	var lock sync.Mutex

	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}

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
	wg.Wait()
}

func pForEach(s *easySlice, consumer TConsumer) {
	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	var wg sync.WaitGroup

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
	wg.Wait()
}

func pAllMatch(s *easySlice) bool {
	workerNum := 1
	subSliceSize := s.collection.Len() / workerNum
	var wg sync.WaitGroup
	var failedChan = make(chan struct{}, 1)
	defer func() {
		close(failedChan)
	}()

	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}

	var exitFlag int32 = 0
	for i := 0; i < workerNum; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if exitFlag > 0 {
					return
				}
				f, _ := s.evaluate(j)
				if !f {
					select {
					case failedChan <- struct{}{}:
						atomic.AddInt32(&exitFlag, 1)
					default:
					}
					return
				}
			}
		}(start, end)
	}
	wg.Wait()
	select {
	case <-failedChan:
		return false
	default:
		return true
	}
}

func pAnyMatch(s *easySlice) bool {
	workerNum := getNumWorkers() * 10
	subSliceSize := s.collection.Len() / workerNum
	var wg sync.WaitGroup

	if subSliceSize == 0 {
		subSliceSize = 1
		workerNum = s.collection.Len()
	}

	var exitFlag uint32 = 0
	for i := 0; i < workerNum; i++ {
		start := i * subSliceSize
		end := start + subSliceSize
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if exitFlag > 0 {
					return
				}
				f, _ := s.evaluate(j)
				if f {
					atomic.AddUint32(&exitFlag, 1)
					return
				}
			}
		}(start, end)
	}
	wg.Wait()
	if exitFlag > 0 {
		return true
	}
	return false
}

func pFindAny(s *easySlice) IOptional {
	workerNum := getNumWorkers()
	subSliceSize := s.collection.Len() / workerNum
	var wg sync.WaitGroup
	var foundChan = make(chan interface{}, 1)
	var closeChan = make(chan struct{}, 1)
	defer func() {
		close(closeChan)
		close(foundChan)
	}()

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
				select {
				case <-closeChan:
					closeChan <- struct{}{}
					return
				default:
					f, v := s.evaluate(j)
					if f {
						select {
						case foundChan <- v:
							closeChan <- struct{}{}
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
