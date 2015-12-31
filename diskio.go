package repono

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var mu sync.RWMutex

func WriteStore(store string) {
	mu.Lock()
	if err := os.MkdirAll(DB_PATH+store, 0755); err != nil {
		log.Fatal(err)
	}
	mu.Unlock()
}

func DeleteStore(store string) {
	mu.Lock()
	if err := os.RemoveAll(DB_PATH + store); err != nil {
		log.Fatal(err)
	}
	mu.Unlock()
}

func WriteData(store string, key, val []byte) {
	mu.Lock()
	if err := ioutil.WriteFile(DB_PATH+store+"/"+string(key)+".json", val, 0644); err != nil {
		log.Fatal(err)
	}
	mu.Unlock()
}

func DeleteData(store string, key []byte) {
	mu.Lock()
	if err := os.Remove(DB_PATH + store + "/" + string(key) + ".json"); err != nil {
		log.Fatal(err)
	}
	mu.Unlock()
}
