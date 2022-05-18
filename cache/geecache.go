package cache

import (
	"fmt"
	"log"
	"sync"
)

// 定义接口 Getter 和 回调函数 Get(key string)([]byte, error)，参数是 key，返回值是 []byte
type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义函数类型 GetterFunc，并实现 Getter 接口的 Get 方法
type GetterFunc func(key string) ([]byte, error)

//函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group 可以认为是一个缓存的命名空间
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu sync.RWMutex

	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheByte int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheByte,
		},
	}

	groups[name] = g
	return g
}

//GetGroup 用来特定名称的 Group，这里使用了只读锁 RLock()，因为不涉及任何冲突变量的写操作
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]

	mu.RUnlock()

	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key si required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Cache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)

}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{
		b: colneBytes(bytes),
	}

	g.popluateCache(key, value)
	return value, nil
}

func (g *Group) popluateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
