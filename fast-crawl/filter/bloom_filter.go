package filter

import (
	"github.com/go-redis/redis"
	"log"
)

type BloomFilter struct {
	size     uint
	seeds    []uint
	set      BloomSet
	function func(uint, uint, string) uint
}

func NewBloomFilter(size uint, seeds []uint, set BloomSet, hashFunc func(uint, uint, string) uint) *BloomFilter {
	bf := new(BloomFilter)
	bf.size = size
	bf.seeds = seeds
	bf.function = hashFunc
	bf.set = set
	return bf
}

func (bf *BloomFilter) Add(value string) {
	values := make([]interface{}, 0)
	for _, s := range bf.seeds {
		values = append(values, bf.function(bf.size, s, value))
	}
	bf.set.SetAll(values)
}

func (bf *BloomFilter) Contains(value string) bool {
	if value == "" {
		return false
	}
	ret := make([]interface{}, 0)
	for _, s := range bf.seeds {
		ret = append(ret, bf.function(bf.size, s, value))
	}
	return bf.set.ContainAll(ret)
}

// hash 方法
func DefaultHash(cap uint, seed uint, value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*seed + uint(value[i])
	}
	return (cap - 1) & result
}

type BloomSet interface {
	Set(v interface{})
	SetAll(values []interface{})
	Contain(v interface{}) bool
	ContainAll(values []interface{}) bool
}

type RedisSet struct {
	length    uint
	client    *redis.Client
	key       string
	buffer    *Set
	bufferMax int
}

func NewRedisSet(length uint) *RedisSet {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	client.FlushDB()

	return &RedisSet{
		length,
		client,
		"key",
		NewSet(),
		10000,
	}
}

func (r *RedisSet) Set(i interface{}) {
	r.client.SAdd(r.key, i).Err()
}

func (r *RedisSet) SetAll(values []interface{}) {
	r.buffer.Add(values...)
	if r.buffer.Size() >= r.bufferMax {
		pipe := r.client.Pipeline()
		pipe.SAdd(r.key, r.buffer.ToSlice())
		_, err := pipe.Exec()
		if err != nil {
			log.Println(err)
		}
		r.buffer.Clear()
	}
}

func (r *RedisSet) Contain(i interface{}) bool {
	if r.buffer.Contains(i) {
		return true
	}
	return r.client.SIsMember(r.key, i).Val()
}

func (r *RedisSet) ContainAll(values []interface{}) bool {
	ret := true
	for _, v := range values {
		ret = ret && r.buffer.Contains(v)
	}
	if ret {
		return ret
	}

	pipe := r.client.Pipeline()
	for _, v := range values {
		pipe.SIsMember(r.key, v)
	}
	cmders, err := pipe.Exec()
	if err != nil {
		return false
	}
	ret = true
	for _, cmd := range cmders {
		val := cmd.(*redis.BoolCmd).Val()
		ret = ret && val
	}
	return ret
}

var Exists = struct{}{}

type Set struct {
	m map[interface{}]struct{}
}

func NewSet(items ...interface{}) *Set {
	// 获取Set的地址
	s := &Set{}
	// 声明map类型的数据结构
	s.m = make(map[interface{}]struct{})
	s.Add(items...)
	return s
}

func (s *Set) Add(items ...interface{}) error {
	for _, item := range items {
		s.m[item] = Exists
	}
	return nil
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

func (s *Set) ToSlice() []interface{} {
	result := make([]interface{}, 0)
	for k, _ := range s.m {
		result = append(result, k)
	}
	return result
}
