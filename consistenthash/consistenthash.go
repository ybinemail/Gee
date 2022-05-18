package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash           //采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，也方便测试时替换，默认为 crc32.ChecksumIEEE 算法
	replicas int            //虚拟节点倍数 replicas
	keys     []int          //哈希环 keys
	hashMap  map[int]string //虚拟节点与真实节点的映射表,键是虚拟节点的哈希值，值是真实节点的名称
}

func New(replicas int, fh Hash) *Map {
	m := &Map{
		hash:     fh,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 加真实节点/机器 到哈希环
func (m *Map) Add(keys ...string) {

	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) //对每一个真实节点 key，对应创建 m.replicas 个虚拟节点，虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点
			m.keys = append(m.keys, hash)                      //使用 m.hash() 计算虚拟节点的哈希值，使用 append(m.keys, hash) 添加到环上
			m.hashMap[hash] = key                              //在 hashMap 中增加虚拟节点和真实节点的映射关系
		}
	}

	sort.Ints(m.keys) //环上的哈希值排序
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key))) //计算 key 的哈希值

	// 顺时针找到第一个匹配的虚拟节点的下标 idx，
	idx := sort.Search(len(m.keys), func(i int) bool {
		// 从 m.keys 中获取到对应的哈希值
		return m.keys[i] >= hash
	})

	// 如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	return m.hashMap[m.keys[idx%len(m.keys)]] //通过 hashMap 映射得到真实的节点
}
