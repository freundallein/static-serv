package cache

import (
	"hash/fnv"
	"net/http"
	"net/http/httptest"
	"time"
)

// Cache - basic cache interface
type Cache interface {
	Get(key string) (*Item, bool)
	Set(key string, item *Item)
	GarbageCollect(timeToExpire time.Duration)
}

// New - cache factory
func New(expiration time.Duration) Cache {
	return &MapCache{
		store:      map[uint64]*Item{},
		expiration: expiration,
	}
}

// Item - basic cache item representation
type Item struct {
	header   http.Header
	data     []byte
	lastSeen time.Time
}

// NewItem - constructor
func NewItem(recorder *httptest.ResponseRecorder) *Item {
	return &Item{
		data:     recorder.Body.Bytes(),
		header:   recorder.Result().Header,
		lastSeen: time.Now().UTC(),
	}
}

// Header - header data getter
func (i *Item) Header() http.Header {
	return i.header
}

// Data - item's data getter
func (i *Item) Data() []byte {
	return i.data
}

//isExpired - check if current item is stale
func (i *Item) isExpired(expiration time.Duration) bool {
	return time.Since(i.lastSeen) > expiration
}

// constructKey - make hash from uri
func constructKey(uri string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(uri))
	return hash.Sum64()
}
