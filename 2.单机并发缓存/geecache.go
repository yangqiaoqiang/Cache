package Cache

import (
	"fmt"
	"log"
	"sync"
)

//缓存不存在调用回调函数，得到源数据
type Getter interface {
	Get(key string) ([]byte, error)
}

//接口型函数 GetterFunc 实现 Getter
type GetterFunc func(key string) ([]byte, error)

//Getter实现
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

//核心数据结构 Group
type Group struct {
	name      string //唯一name
	getter    Getter //缓存未命中的callback
	mainCache cache  //一开始实现并发缓存
}

var mu sync.RWMutex                  //读写锁
var groups = make(map[string]*Group) //保存Group组

// NewGroup 创建 Group 实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup 用来特定名称的 Group
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Cache核心方法Get
//缓存中存在直接获取，不存在通过回调函数添加
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}
	return g.load(key)
}
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key) // callback函数
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
