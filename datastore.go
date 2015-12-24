package repono

import (
	"regexp"
	"runtime"
	"sync"
)

var PATH string = "db/"

type DataStore struct {
	Stores map[string]*Store
	sync.RWMutex
}

func NewDataStore() *DataStore {
	ds := &DataStore{
		Stores: make(map[string]*Store, 0),
	}
	if PATH[len(PATH)-1] != '/' {
		PATH += "/"
	}
	// check if ds is empty; walk path and load; gc
	if data := Walk(PATH); len(data) > 0 {
		ds.Lock()
		for store, files := range data {
			st := NewStore(store)
			ds.Stores[store] = st
			st.Load(files)
		}
		ds.Unlock()
		runtime.GC()
	}
	return ds
}

func (ds *DataStore) AddStore(store string) {
	if _, ok := ds.GetStore(store); !ok {
		ds.Lock()
		ds.Stores[store] = NewStore(store)
		WriteStore(store)
		ds.Unlock()
	}
}

// for internal use only
func (ds *DataStore) GetStore(store string) (*Store, bool) {
	ds.RLock()
	st, ok := ds.Stores[store]
	ds.RUnlock()
	return st, ok
}

func (ds *DataStore) DelStore(store string) {
	if _, ok := ds.GetStore(store); !ok {
		ds.Lock()
		delete(ds.Stores, store)
		DeleteStore(store)
		ds.Unlock()
	}
}

func (ds *DataStore) HasStore(store string) bool {
	ds.RLock()
	_, ok := ds.Stores[store]
	ds.RUnlock()
	return ok
}

func (ds *DataStore) Add(store string, key, val []byte) bool {
	if st, ok := ds.GetStore(store); ok {
		return st.Add(key, val)
	}
	return false
}

func (ds *DataStore) Set(store string, key, val []byte) bool {
	if st, ok := ds.GetStore(store); ok {
		st.Set(key, val)
		return true
	}
	return false
}

func (ds *DataStore) Get(store string, key []byte) []byte {
	if st, ok := ds.GetStore(store); ok {
		return st.Get(key)
	}
	return nil
}

func (ds *DataStore) GetAll(store string) [][]byte {
	if st, ok := ds.GetStore(store); ok {
		return st.GetAll()
	}
	return nil
}

func (ds *DataStore) Del(store string, key []byte) {
	if st, ok := ds.GetStore(store); ok {
		st.Del(key)
	}
}

func (ds *DataStore) Has(store string, key []byte) bool {
	if st, ok := ds.GetStore(store); ok {
		return st.Has(key)
	}
	return false
}

func (ds *DataStore) Query(store string, re *regexp.Regexp) [][]byte {
	if st, ok := ds.GetStore(store); ok {
		return st.Query(re)
	}
	return nil
}

func (ds *DataStore) Import() {

}

func (ds *DataStore) Export() {

}
