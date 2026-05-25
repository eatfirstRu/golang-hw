package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	i, ok := c.items[key]
	if ok {
		i.Value = value
		c.queue.MoveToFront(i)
	} else {
		c.items[key] = c.queue.PushFront(value)
	}
	if len(c.items) > c.capacity {
		last := c.queue.Back()
		for k, v := range c.items {
			if v == last {
				delete(c.items, k)
			}
			break
		}
		c.queue.Remove(c.queue.Back())
	}
	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	i, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(i)
		return i.Value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}
