package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32
type Map struct {
	hash     Hash
	hashMap  map[int]string
	replicas int
	keys     []int
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		hashMap:  make(map[int]string),
		replicas: replicas,
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}
func (m *Map) Add(keys ...string) {
	for _, k := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + k)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = k
		}
		sort.Ints(m.keys)
	}
}
func (m *Map) Get(data string) string {
	if len(data) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(data)))
	search := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[search%len(m.keys)]]
}
