package repono

import (
	"os"
	"path/filepath"
	"sort"
)

func Walk(start string) map[string][]string {
	w := Walker{
		StartDir: start[:len(start)-1],
		Stores:   make(map[string][]string),
	}
	filepath.Walk(w.StartDir, w.Texas)
	for k := range w.Stores {
		w.Stores[k] = w.Ranger(k)
	}
	return w.Stores
}

type Walker struct {
	StartDir string
	Stores   map[string][]string
}

// walks the db root and gathers all the stores/folders
func (w *Walker) Texas(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.Name() == w.StartDir {
		return nil
	}
	if info.IsDir() {
		w.Stores[info.Name()] = make([]string, 0)
		return filepath.SkipDir
	}
	return nil
}

// takes folder/store as key and walks files/docs...
// returns list of files/docs in this folder/store
func (w *Walker) Ranger(key string) []string {
	var files []string
	filepath.Walk(w.StartDir+"/"+key, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == key {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		files = append(files, info.Name())
		return nil
	})
	sort.Strings(files)
	return files
}
