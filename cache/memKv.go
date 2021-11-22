package cache

import (
	"errors"
	"fmt"
	"sync"
)

type MemKv struct {
	kv   map[string]interface{}
	lock *sync.RWMutex
}

func NewMemKv() *MemKv {
	return &MemKv{
		kv:   make(map[string]interface{}),
		lock: &sync.RWMutex{},
	}
}

func (mkv *MemKv) Get(key string) (interface{}, error) {
	mkv.lock.RLock()
	defer mkv.lock.RUnlock()
	v, ok := mkv.kv[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[MemKv] key %s not found", key))
	}
	return v, nil
}

func (mkv *MemKv) Put(key string, value interface{}) {
	mkv.lock.Lock()
	defer mkv.lock.Unlock()
	mkv.kv[key] = value
}

func (mkv *MemKv) Del(key string) {
	mkv.lock.Lock()
	defer mkv.lock.Unlock()
	delete(mkv.kv, key)
}
