package cache

import (
	"container/list"
)

type LRUIndex struct {
	size  int
	list  *list.List
	items map[string]*list.Element
}

func NewLRUIndex(size int) *LRUIndex {
	return &LRUIndex{
		size:  size,
		list:  list.New(),
		items: make(map[string]*list.Element),
	}
}

func (index *LRUIndex) IsAtCapacity() bool {
	return index.list.Len() >= index.size
}

func (index *LRUIndex) String() string {
	out := ""

	element := index.list.Front()

	if element == nil {
		return ""
	}

	for {
		if str, ok := element.Value.(string); ok {
			out += str
		}

		element = element.Next()

		if element == nil {
			break
		}

		out += ","
	}

	return out
}

func (index *LRUIndex) GetOldest() *string {
	element := index.list.Back()

	if element == nil {
		return nil
	}

	if str, ok := element.Value.(string); ok {
		return &str
	}

	return nil
}

func (index *LRUIndex) Update(key string) {
	element, exists := index.items[key]

	if exists {
		index.list.MoveToFront(element)
	} else {
		element := index.list.PushFront(key)
		index.items[key] = element
	}
}

func (index *LRUIndex) Remove(key string) {
	if element, exists := index.items[key]; exists {
		index.list.Remove(element)
		delete(index.items, key)
	}
}
