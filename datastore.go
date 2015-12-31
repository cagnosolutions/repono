package repono

import (
	"log"
	"regexp"
	"runtime"
	"sync"
)

type DataStore struct {
	Stores map[string]*Store
	sync.RWMutex
}

func NewDataStore() *DataStore {
	ds := &DataStore{
		Stores: make(map[string]*Store, 0),
	}
	if DB_PATH[len(DB_PATH)-1] != '/' {
		DB_PATH += "/"
	}
	if data := Walk(DB_PATH); len(data) > 0 {
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

func (ds *DataStore) AddStore(store string) []byte {
	if _, ok := ds.GetStore(store); !ok {
		ds.Lock()
		ds.Stores[store] = NewStore(store)
		WriteStore(store)
		ds.Unlock()
	}
	return TRUE
}

// for internal use only
func (ds *DataStore) GetStore(store string) (*Store, bool) {
	ds.RLock()
	st, ok := ds.Stores[store]
	ds.RUnlock()
	return st, ok
}

func (ds *DataStore) DelStore(store string) []byte {
	if _, ok := ds.GetStore(store); !ok {
		ds.Lock()
		delete(ds.Stores, store)
		DeleteStore(store)
		ds.Unlock()
	}
	return TRUE
}

func (ds *DataStore) HasStore(store string) []byte {
	ds.RLock()
	_, ok := ds.Stores[store]
	ds.RUnlock()
	if ok {
		return TRUE
	}
	return FALSE
}

func (ds *DataStore) UUID() []byte {
	return UUID1()
}

func (ds *DataStore) Add(store string, key, val []byte) []byte {
	st, ok := ds.GetStore(store)
	if ok && st.Add(key, val) {
		return TRUE
	}
	return FALSE
}

func (ds *DataStore) Set(store string, key, val []byte) []byte {
	if st, ok := ds.GetStore(store); ok {
		st.Set(key, val)
		return TRUE
	}
	return FALSE
}

func (ds *DataStore) Get(store string, key []byte) []byte {
	if st, ok := ds.GetStore(store); ok {
		return st.Get(key)
	}
	return NIL
}

func (ds *DataStore) GetAll(store string) []byte {
	if st, ok := ds.GetStore(store); ok {
		return formatList(st.GetAll())
	}
	return NIL
}

func (ds *DataStore) Del(store string, key []byte) []byte {
	if st, ok := ds.GetStore(store); ok {
		st.Del(key)
	}
	return TRUE
}

func (ds *DataStore) Has(store string, key []byte) []byte {
	st, ok := ds.GetStore(store)
	if ok && st.Has(key) {
		return TRUE
	}
	return FALSE
}

func (ds *DataStore) Query(store, query string) []byte {
	re, err := regexp.Compile(query)
	if err != nil {
		log.Fatal(err)
	}
	if st, ok := ds.GetStore(store); ok {
		return formatList(st.Query(re))
	}
	return nil
}

func (ds *DataStore) Import() {

}

func (ds *DataStore) Export() {

}
