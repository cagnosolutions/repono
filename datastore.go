package repono

import (
	"bytes"
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
	if PATH[len(PATH)-1] != '/' {
		PATH += "/"
	}
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

func (ds *DataStore) AddStore(store string) []byte {
	if _, ok := ds.getStore(store); !ok {
		ds.Lock()
		ds.Stores[store] = NewStore(store)
		WriteStore(store)
		ds.Unlock()
	}
	return TRUE
}

// for internal use only
func (ds *DataStore) getStore(store string) (*Store, bool) {
	ds.RLock()
	st, ok := ds.Stores[store]
	ds.RUnlock()
	return st, ok
}

func (ds *DataStore) DelStore(store string) []byte {
	if _, ok := ds.getStore(store); !ok {
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
	return []byte(UUID1())
}

// datastore generates and returns key on insert, behaves like add
func (ds *DataStore) Put(store string, val []byte) []byte {
	st, ok := ds.getStore(store)
	key := []byte(UUID1())
	if ok && st.Add(key, val) {
		return key
	}
	return NIL
}

func (ds *DataStore) Add(store string, key, val []byte) []byte {
	st, ok := ds.getStore(store)
	if ok && st.Add(key, val) {
		return TRUE
	}
	return FALSE
}

func (ds *DataStore) Set(store string, key, val []byte) []byte {
	if st, ok := ds.getStore(store); ok {
		st.Set(key, val)
		return TRUE
	}
	return FALSE
}

func (ds *DataStore) Get(store string, key []byte) []byte {
	if st, ok := ds.getStore(store); ok {
		return st.Get(key)
	}
	return NIL
}

func (ds *DataStore) GetAll(store string) []byte {
	if st, ok := ds.getStore(store); ok {
		return formatList(st.GetAll())
	}
	return NIL
}

func (ds *DataStore) Del(store string, key []byte) []byte {
	if st, ok := ds.getStore(store); ok {
		st.Del(key)
	}
	return TRUE
}

func (ds *DataStore) Has(store string, key []byte) []byte {
	st, ok := ds.getStore(store)
	if ok && st.Has(key) {
		return TRUE
	}
	return FALSE
}

func (ds *DataStore) Query(store string, query [][]byte) []byte {
	if st, ok := ds.getStore(store); ok {
		return formatList(st.Query(query))
	}
	return NIL
}

func (ds *DataStore) Import() {

}

func (ds *DataStore) Export() {

}

func formatList(bb [][]byte) []byte {
	if bb != nil {
		bb[0] = append([]byte{'['}, bb[0]...)
		bb[len(bb)-1] = append(bb[len(bb)-1], ']')
		return bytes.Join(bb, []byte{','})
	}
	return NIL
}
