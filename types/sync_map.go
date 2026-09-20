package types

import "sync"

type SyncMap[TKey comparable, TValue any] struct {
	sync.RWMutex

	data map[TKey]TValue
}

func NewSyncMap[TKey comparable, TValue any]() *SyncMap[TKey, TValue] {
	return &SyncMap[TKey, TValue]{
		data: make(map[TKey]TValue),
	}
}

func (receiver *SyncMap[TKey, TValue]) Get(key TKey) (TValue, bool) {
	receiver.RLock()
	defer receiver.RUnlock()

	value, ok := receiver.data[key]
	return value, ok
}

func (receiver *SyncMap[TKey, TValue]) Set(key TKey, value TValue) {
	receiver.Lock()
	defer receiver.Unlock()

	receiver.data[key] = value
}

func (receiver *SyncMap[TKey, TValue]) Delete(key TKey) {
	receiver.Lock()
	defer receiver.Unlock()
	delete(receiver.data, key)
}
