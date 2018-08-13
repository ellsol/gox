package gox

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DiskCache struct {
	Values map[string]string
	Path   string
}

func NewDiskCache(filename string) (*DiskCache, error) {
	cache := new(DiskCache)
	cache.Path = filename

	err := cache.Load()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func (it *DiskCache) Get(key string) string {
	val, ok := it.Values[key]
	if ok {
		return val
	} else {
		return ""
	}
}

func (it *DiskCache) Save(key string, value string) error {
	m := it.Values
	m[key] = value

	return it.SaveToDisk()
}

func (it *DiskCache) SaveToDisk() error {
	json, err := json.Marshal(it.Values)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(it.Path, json, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (it *DiskCache) Load() error {
	values := make(map[string]string, 0)

	if _, err := os.Stat(it.Path); os.IsNotExist(err) {
		f, err := os.Create(it.Path)
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(it.Path)
	if err != nil {
		return err
	}

	if len(b) > 0 {
		err = json.Unmarshal(b, &values)

		if err != nil {
			return err
		}

		it.Values = values
	} else {
		it.Values = values
	}

	return nil
}
