package octopus

import (
	"sync"
)

type value struct {
	sync.RWMutex
	data map[interface{}]interface{}
}

// Create a new key-value pair
func (v *value) Set(key interface{}, val interface{}) {
	v.Lock()
	defer v.Unlock()
	if v.data == nil {
		v.data = make(map[interface{}]interface{})
	}
	v.data[key] = val
}

// Read a value by key
func (v *value) Get(key interface{}) (interface{}, bool) {
	v.RLock()
	defer v.RUnlock()
	val, ok := v.data[key]
	return val, ok
}

// Update a value by key
func (v *value) Update(key interface{}, val interface{}) bool {
	v.Lock()
	defer v.Unlock()
	_, ok := v.data[key]
	if ok {
		v.data[key] = val
	}
	return ok
}

// Delete a key-value pair
func (v *value) Delete(key interface{}) {
	v.Lock()
	defer v.Unlock()
	delete(v.data, key)
}
