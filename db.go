package peregrinedb

import (
	"errors"
	"log"
	"os"
	"sort"
)

const (
	MaxValueSize int = 1 << 9 // in bytes
)

type Data struct {
}

type Index struct {
	Key    int64
	PageId int64
	Ind    int64
}

type DB struct {
	file string
	page Page
}

func Open(name string) *DB {
	mmap, err := os.OpenFile(name, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(mmap)
	}

	data, err := os.OpenFile("data_page_01", os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{
		file: name,
		page: Page{
			dataF: data,
			mmap:  mmap,
		},
	}
}

func (db *DB) Get(key int64) []byte {
	// read index page
	// find the data page from index page
	return nil
}

func (db *DB) Put(key int64, value []byte) error {
	if len(value) > int(MaxValueSize) {
		return errors.New("value size is large")
	}

	// buff, err := io.ReadAll(db.mmap)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	indM, indexs := db.page.getIndexPage()
	// updating the metadata
	indM.Size += 1
	indM.MaxKey = max(indM.MaxKey, key)
	indM.MinKey = min(indM.MaxKey, key)

	ind := sort.Search(len(indexs), func(i int) bool {
		// its like give me the smallest index where
		// index[i].key is not small than key
		return !(indexs[i].Key < key)
	})

	if indexs[ind].Key == key {
		// TODO update the value
		indexs[ind].Ind += 0
		return nil
	}

	// insert into index
	// insert into data

	return nil
}

func (db *DB) Delete(key int64) error {
	return nil
}

// sorted slice as an index
