package cache

import (
	"Gcache/singleflight"
	"errors"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var groups = map[string]*Group{}
var rmu sync.RWMutex

type Group struct {
	name      string
	getter    Getter
	mainCache Cache
	peers     PeerPicker
	loader    *singleflight.Group
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}
func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				return value, nil
			}
			log.Println("cache failed to get from peer", err)
		}
		return g.GetLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}
	rmu.Lock()
	defer rmu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: Cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}
func GetGroup(name string) *Group {
	rmu.RLock()
	defer rmu.RUnlock()
	return groups[name]
}
func (g *Group) Get(key string) (value ByteView, err error) {
	if len(key) == 0 {
		return ByteView{}, errors.New("key is required")
	}
	if get, ok := g.mainCache.get(key); ok {
		log.Println("cache hit")
		return get, nil
	}
	return g.load(key)

}

// GetLocally
//func (g *Group) Load(key string) (value ByteView, err error) {
//	return g.GetLocally(key)
//}
func (g *Group) GetLocally(key string) (value ByteView, err error) {
	get, err := g.getter.Get(key)
	if err != nil {
		return
	}
	value = ByteView{b: cloneBytes(get)}
	g.population(key, value)
	return
}
func (g *Group) population(key string, value ByteView) {
	g.mainCache.add(key, value)
}
