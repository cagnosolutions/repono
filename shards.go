package repono

import (
	"bytes"
	"hash/fnv"
	"io"
	"log"
	"sort"
	"sync"

	"github.com/cagnosolutions/bplus"
)

var (
	SHARDCOUNT = 64
)

type Shards []*Shard

type Shard struct {
	data *bplus.Tree
	sync.RWMutex
}

func NewShards() *Shards {
	s := make(Shards, SHARDCOUNT)
	for i := 0; i < SHARDCOUNT; i++ {
		s[i] = &Shard{data: bplus.NewTree(bytes.Compare)}
		// NOTE: close tree
	}
	return &s
}

func (s Shards) shard(key []byte) *Shard {
	h := fnv.New32()
	h.Write(key)
	return s[uint(h.Sum32())%uint(SHARDCOUNT)]
}

func (s *Shards) Add(key, val []byte) bool {
	shard := s.shard(key)
	shard.Lock()
	defer shard.Unlock()
	_, ok := shard.data.Put(key, func(old []byte, exists bool) ([]byte, bool) {
		if exists {
			return old, false
		}
		return val, true
	})
	return ok
}

func (s *Shards) Set(key, val []byte) {
	shard := s.shard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.data.Set(key, val)
}

func (s Shards) Get(key []byte) []byte {
	shard := s.shard(key)
	shard.Lock()
	defer shard.Unlock()
	b, _ := shard.data.Get(key)
	return b
}

func (s Shards) GetAll() DataSet {
	var dataSet DataSet
	for data := range s.Iter() {
		dataSet = append(dataSet, data)
	}
	sort.Stable(dataSet)
	return dataSet
}

func (s *Shards) Del(key []byte) {
	if shard := s.shard(key); shard != nil {
		shard.Lock()
		shard.data.Delete(key)
		shard.Unlock()
	}
}

func (s Shards) Has(key []byte) bool {
	shard := s.shard(key)
	shard.Lock()
	defer shard.Unlock()
	_, ok := shard.data.Get(key)
	return ok
}

func (s Shards) Size() int {
	var n int
	for i := 0; i < SHARDCOUNT; i++ {
		shard := s[i]
		shard.RLock()
		n += shard.data.Len()
		shard.Unlock()
	}
	return n
}

func (s Shards) Iter() <-chan Data {
	ch := make(chan Data)
	go func() {
		for _, shard := range s {
			shard.RLock()
			enum, err := shard.data.SeekFirst()
			if err != nil {
				break
			}
			for {
				k, v, err := enum.Next()
				if err != nil {
					if err != io.EOF {
						log.Fatal(err)
					}
					break
				}
				ch <- Data{k, v}
			}
			enum.Close()
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

type Data struct {
	key, val []byte
}

type DataSet []Data

func (d DataSet) Len() int {
	return len(d)
}

func (d DataSet) Less(i, j int) bool {
	return bytes.Compare(d[i].key, d[j].key) == -1
}

func (d DataSet) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
