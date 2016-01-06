package repono

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Store struct {
	Name   string
	shards *Shards
}

func NewStore(name string) *Store {
	return &Store{
		Name:   name,
		shards: NewShards(),
	}
}

func (st *Store) Add(key, val []byte) bool {
	ok := st.shards.Add(key, val)
	if ok {
		go WriteData(st.Name, key, val)
	}
	return ok
}

func (st *Store) Set(key, val []byte) {
	st.shards.Set(key, val)
	go WriteData(st.Name, key, val)
}

func (st *Store) Get(key []byte) []byte {
	return st.shards.Get(key)
}

func (st *Store) GetAll() [][]byte {
	var vals [][]byte
	for _, data := range st.shards.GetAll() {
		vals = append(vals, data.V)
	}
	return vals
}

func (st *Store) Del(key []byte) {
	if st.Has(key) {
		st.shards.Del(key)
		go DeleteData(st.Name, key)
	}
}

func (st *Store) Has(key []byte) bool {
	return st.shards.Has(key)
}

func (st *Store) Load(files []string) {
	for _, file := range files {
		data, err := ioutil.ReadFile(DB_PATH + st.Name + "/" + file)
		if err != nil {
			log.Fatal(err)
		}
		st.Add([]byte(strings.Split(file, ".")[0]), data)
	}
}

func (st *Store) Query(query [][]byte) [][]byte {
	var vals [][]byte
	for _, data := range st.shards.Query(query) {
		vals = append(vals, data.V)
	}
	return vals
}

func (st *Store) _Query(re *regexp.Regexp) [][]byte {
	var match [][]byte
	for _, data := range st.shards.GetAll() {
		if idx := re.FindIndex(data.V); idx != nil {
			match = append(match, data.V)
		}
	}
	return match
}
