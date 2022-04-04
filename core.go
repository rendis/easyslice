package easyslice

import (
	"runtime"
	"sync"
)

type core struct {
	cpuNum        int
	cpuMultiplier int
}

var (
	coreInstance *core
	coreOnce     sync.Once
)

func getNumWorkers() int {
	coreOnce.Do(func() {
		coreInstance = &core{
			cpuNum:        runtime.NumCPU(),
			cpuMultiplier: 2,
		}
	})
	return coreInstance.cpuNum * coreInstance.cpuMultiplier
}

func SetCoreMultiplier(multiplier int) {
	coreOnce.Do(func() {
		coreInstance = &core{
			cpuNum:        runtime.NumCPU(),
			cpuMultiplier: multiplier,
		}
	})
}
