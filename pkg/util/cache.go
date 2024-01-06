package util

// implement a (very) simple cache with a map[string]string
// is a singleton

type Cache struct {
	cache map[string]interface{}
}

var theOnlyOneCache Cache

func NewCache() *Cache {
	if theOnlyOneCache.cache == nil {
		theOnlyOneCache.Clear()
	}

	return &theOnlyOneCache
}

func (c *Cache) Get(key string) (value interface{}, ok bool) {
	value, ok = c.cache[key]
	return
}

func (c *Cache) Set(key string, value interface{}) {
	c.cache[key] = value
}

func (c *Cache) Delete(key string) {
	delete(c.cache, key)
}

func (c *Cache) Clear() {
	//c.cache = make(map[string]interface{}, ninCacheSize)
	c.cache = map[string]interface{}{}
}
