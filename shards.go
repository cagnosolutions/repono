package repono

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"sort"
	"sync"

	"github.com/cagnosolutions/bplus"
)

var (
	SHARDCOUNT = 64
	DEBUG      = false
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
	bucket := uint(h.Sum32()) % uint(SHARDCOUNT)
	return s[bucket]
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

func (s Shards) GetAllVals() [][]byte {
	var vals [][]byte
	for _, dataSet := range s.GetAll() {
		vals = append(vals, dataSet.V)
	}
	return vals
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

func (s Shards) Query(query [][]byte) DataSet {
	ch, active := make(chan []Data), 0
	for i := 0; i < SHARDCOUNT; i++ {
		if s[i].data.Len() > 0 {
			go search(s[i], ch, query)
			active++
		}
	}
	var dataSet DataSet
	if DEBUG {
		log.Printf(">>>>> Found %d active shards...\n", active)
	}
	for i := 0; i < active; i++ {
		select {
		case data := <-ch:
			if DEBUG {
				log.Printf(">>>>> Got data: %+v\n", data)
			}
			dataSet = append(dataSet, data...)
		}
	}
	close(ch)
	sort.Sort(dataSet)
	if DEBUG {
		fmt.Printf(">>>>> Sorted Found DataSet: %+v\n", dataSet)
	}
	return dataSet
}

func search(sh *Shard, ch chan []Data, query [][]byte) {
	sh.RLock()
	defer sh.RUnlock()
	enum, err := sh.data.SeekFirst()
	if err != nil {
		if err != io.EOF {
			log.Fatalf("Unknown tree.SeekFirst() error: %s\n", err)
		}
		ch <- []Data{}
		return // go an io.EOF
	}
	defer enum.Close()
	res := make([]Data, 0)
	var match bool
	for {
		match = true
		k, v, err := enum.Next()
		if err != nil {
			if err != io.EOF {
				log.Fatalf("Unknown enum.Next() error: %s\n", err)
			}
			break // got an io.EOF
		}
		for _, q := range query {
			if !bytes.Contains(v, q) {
				match = false
				break
			}
		}
		if match {
			res = append(res, Data{k, v})
		}
	}
	ch <- res
	return
}

func (s Shards) Iter() <-chan Data {
	ch := make(chan Data)
	go func() {
		for _, shard := range s {
			shard.RLock()
			enum, err := shard.data.SeekFirst()
			if err != nil {
				shard.RUnlock()
				continue
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
	K, V []byte
}

type DataSet []Data

func (ds DataSet) Len() int {
	return len(ds)
}

func (ds DataSet) Less(i, j int) bool {
	return bytes.Compare(ds[i].K, ds[j].V) == -1
}

func (ds DataSet) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}
