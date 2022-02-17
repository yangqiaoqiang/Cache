package Cache

import (
	pb "Cache/geecachepb"
	"Cache/singleflight"
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
	name      string              //唯一name
	getter    Getter              //缓存未命中的callback
	mainCache cache               //一开始实现并发缓存
	peers     PeerPicker          //选择节点
	loader    *singleflight.Group //并发读取未缓存的数据时，缓存防止击穿，只提取一次
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
		loader:    &singleflight.Group{},
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

//使用 PickPeer 方法选择节点,若非本机节点,则调用 getFromPeer 从远程获取。
//若是本机节点或失败,则回退到 getLocally
func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFormPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
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

//新增 RegisterPeers 方法,将实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

//新增 getFromPeer 方法，使用实现了 PeerGetter 接口的 httpGetter
//从访问远程节点，获取缓存值。
func (g *Group) getFormPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
