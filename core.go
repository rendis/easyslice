package easyslice

import (
	"runtime"
	"sync"
)

type core struct {
	numWorkers int
}

var (
	coreInstance *core
	coreOnce     sync.Once
)

func getNumWorkers() int {
	coreOnce.Do(func() {
		coreInstance = &core{
			numWorkers: runtime.NumCPU(),
		}
	})
	return coreInstance.numWorkers
}

func SetCoreMultiplier(multiplier float32) {
	coreOnce.Do(func() {
		nw := int(float32(runtime.NumCPU()) * multiplier)
		if nw < 1 {
			nw = 1
		}
		coreInstance = &core{
			numWorkers: nw,
		}
	})
}
