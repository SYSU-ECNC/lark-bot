package cache

import "sync"

type MemoryStore struct {
	sync.RWMutex
	Map map[string]interface{}
}

var memoryStore *MemoryStore

func init() {
	memoryStore = new(MemoryStore)
	memoryStore.Map = make(map[string]interface{})
}

// Get value by given key.
func Get(key string) interface{} {
	memoryStore.RLock()
	value := memoryStore.Map[key]
	memoryStore.RUnlock()
	return value
}

// Set value by given key and value.
func Set(key string, value interface{}) {
	memoryStore.Lock()
	memoryStore.Map[key] = value
	memoryStore.Unlock()
}
