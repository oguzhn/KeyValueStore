package cache

import "sync"

type Cache struct {
	mutex  *sync.RWMutex
	values map[string]string
}

func NewCache() *Cache {
	return &Cache{mutex: &sync.RWMutex{}, values: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, error) {
	c.mutex.RLock()
	value := c.values[key]
	c.mutex.RUnlock()
	return value, nil
}
func (c *Cache) Set(key, value string) error {
	c.mutex.Lock()
	c.values[key] = value
	c.mutex.Unlock()
	return nil
}

func (c *Cache) GetAll() (map[string]string, error) {
	c.mutex.RLock()
	allvalues := make(map[string]string)
	for i, v := range c.values {
		allvalues[i] = v
	}
	c.mutex.RUnlock()
	return allvalues, nil
}
