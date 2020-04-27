package lru

import "container/list"

// Cache is a LRU cache. not safe for concurrent access.
type Cache struct {
	// maximun that allow the memery to own
	maxBytes int64
	// current memory already used
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

// NewCache is the constructor of Cache
func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get get the value of the key
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		return kv.value, true
	}

	return nil, false
}

// RemoveOldest remove the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()

	if ele == nil {
		return
	}

	c.ll.Remove(ele)
	kv := ele.Value.(*entry)
	delete(c.cache, kv.key)

	c.nbytes -= (int64(len(kv.key)) + int64(kv.value.Len()))

	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Add add or update a value to cache
func (c *Cache) Add(key string, value Value) {
	// if exist
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += (int64(value.Len()) - int64(kv.value.Len()))
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// if over max memory then remove one until there is enpught memory again
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len get the item numbers of the whole cache
func (c *Cache) Len() int {
	return c.ll.Len()
}
