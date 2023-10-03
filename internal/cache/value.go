package cache

import (
	"sync"
	"time"
)

// Value creates a thread safe cache value accessible with Get
type Value[T any] struct {
	loaded bool
	value  T
	err    error
	mux    sync.Mutex
	ttl    *time.Time
}

// Get reads the cache value, calling the loader function to populate the cache putVal if it is not currently set.
// ttl sets the expiration date of the cache putVal, after which time it will be reloaded.  Setting ttl to nil will
// prevent the cache from expiring.
func (cv *Value[T]) Get(loader func() (T, error), ttl *time.Time) (T, error) {
	cv.mux.Lock()
	defer cv.mux.Unlock()
	if cv.loaded && cv.ttl != nil && cv.ttl.Before(time.Now()) {
		cv.ttl = nil
		cv.loaded = false
	}
	if !cv.loaded {
		cv.value, cv.err = loader()
		cv.ttl = ttl
		cv.loaded = true
	}
	return cv.value, cv.err
}
