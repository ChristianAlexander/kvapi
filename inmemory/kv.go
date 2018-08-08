package inmemory

import (
	"sync"
)

// KV implements an in-memory key value store, scoped by uid.
type KV struct {
	mu sync.RWMutex

	// The outer key is uid, the inner key is the key of the key/value mapping.
	values map[string]map[string]string
}

// NewKV creates a new in-memory key value store.
func NewKV() *KV {
	return &KV{
		values: make(map[string]map[string]string),
	}
}

func (k *KV) ensureUIDMapExists(UID string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	if _, ok := k.values[UID]; !ok {
		k.values[UID] = make(map[string]string)
	}
}

// Get retrieves a value by key for the user.
func (k *KV) Get(UID string, key string) string {
	k.ensureUIDMapExists(UID)

	k.mu.RLock()
	defer k.mu.RUnlock()

	return k.values[UID][key]
}

// Set stores a value by key for the user.
func (k *KV) Set(UID string, key string, value string) {
	k.ensureUIDMapExists(UID)

	k.mu.Lock()
	defer k.mu.Unlock()

	k.values[UID][key] = value
}
