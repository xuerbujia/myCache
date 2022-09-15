package lru

import "container/list"

type Cache struct {
	cache     map[string]*list.Element
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	OnEvicted func(key string, value Value)
}
type Value interface {
	Len() int
}
type entry struct {
	key   string
	value Value
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
func (c *Cache) Get(key string) (value Value, ok bool) {
	if node, ok := c.cache[key]; ok {
		c.ll.MoveToFront(node)
		kv := node.Value.(*entry)
		return kv.value, true
	}
	return
}
func (c *Cache) Remove() {
	node := c.ll.Back()
	if node != nil {
		c.ll.Remove(node)
		kv := node.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
func (c *Cache) Add(key string, value Value) {
	if node, ok := c.cache[key]; ok {
		c.ll.MoveToFront(node)
		kv := node.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		n := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = n
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.Remove()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

//type ListNode struct {
//	Key  string
//	Val  string
//	Next *ListNode
//	prev *ListNode
//}

//func (this *Cache) Get(key string) string {
//	node := this.hash[key]
//	if node == nil {
//		return ""
//	}
//	if node != this.head {
//		prev := node.prev
//		next := node.Next
//		if next == nil {
//			this.tail = prev
//		} else {
//			next.prev = prev
//		}
//		prev.Next = next
//
//		this.head.prev = node
//		node.Next = this.head
//		this.head = node
//	}
//
//	return node.Val
//}
//func (this *Cache) Put(key string, val string) {
//	node := &ListNode{Key: key, Val: val}
//	if n, ok := this.hash[key]; !ok {
//		this.hash[key] = node
//	} else {
//		n.Val = val
//		this.Get(key)
//	}
//	if this.head == nil {
//		this.head = node
//		this.tail = node
//	} else {
//		this.head.prev = node
//		node.Next = this.head
//		this.head = node
//	}
//	if len(this.hash) == this.l {
//		delete(this.hash, this.tail.Key)
//		tail := this.tail.prev
//		this.tail.prev.Next = nil
//		this.tail = tail
//	}
//}
