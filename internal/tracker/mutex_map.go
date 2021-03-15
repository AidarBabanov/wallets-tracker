package tracker

import (
	"fmt"
	"sync"
)

type MutexMap struct {
	m  map[string]*sync.Mutex //address-mutex
	mu *sync.RWMutex
}

func NewMutexMap() *MutexMap {
	mm := new(MutexMap)
	mm.m = make(map[string]*sync.Mutex)
	mm.mu = new(sync.RWMutex)
	return mm
}

func (mm *MutexMap) Lock(address string) {
	mm.mu.RLock()
	mutex, exist := mm.m[address]
	if !exist {
		mm.mu.RUnlock()
		mm.mu.Lock()
		var exist2 bool
		mutex, exist2 = mm.m[address]
		if !exist2 {
			mutex = new(sync.Mutex)
			mm.m[address] = mutex
		}
		mm.mu.Unlock()
	} else {
		mm.mu.RUnlock()
	}
	mutex.Lock()
}

func (mm *MutexMap) Unlock(address string) error {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	mutex, exist := mm.m[address]
	if !exist {
		return fmt.Errorf("unlocking empty mutex for address: %s", address)
	}
	mutex.Unlock()
	return nil
}
