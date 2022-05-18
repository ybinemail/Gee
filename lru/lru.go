package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes  int64                         //允许使用的最大内存
	nbytes    int64                         //当前已使用的内存
	ll        *list.List                    //双向链表
	cache     map[string]*list.Element      //键是字符串，值是双向链表中对应节点的指针
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数，可以为 nil
}

//双向链表节点的
type entry struct {
	key   string
	value Value
}

// 定义接口类型
type Value interface {
	// 定义方法
	Len() int
}

// init
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {

		c.ll.MoveToFront(ele) // 双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾

		kv, _ := ele.Value.(*entry) //TODO 空接口类型的断言
		return kv.value, true
	}
	return
}

//缓存淘汰，移除最近最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele) // 移除链表中
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                //删除缓存map
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //更新内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 更新缓存
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele) //移动队尾
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value //更新缓存新值
	} else {
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		}) //新增缓存
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest() //判断是否超过容量
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
