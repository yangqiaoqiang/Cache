package consistenthash

//实现一致性哈希算法
import (
	"hash/crc32"
	"sort"
	"strconv"
)

//接口型函数 Hash 方便算法替换
type Hash func(data []byte) uint32

// Map 是一致性哈希算法数据结构
type Map struct {
	hash     Hash           //Hash函数默认为crc32.ChecksumIEEE
	replicas int            //虚拟节点个数
	keys     []int          //哈希环
	hashMap  map[int]string //键是虚拟节点的哈希值，值是真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//实现真实节点的 Add method
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//创建编号区分不同虚拟节点
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

//实现选择节点 Get method
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	//二分法选择虚拟节点
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
