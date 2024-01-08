package peregrinedb

import (
	"errors"
	"log"
	"math"
	"os"
	"sort"
)

const (
	MaxValueSize int = 1 << 9 // in bytes
	MaxKeySize   int = 1 << 4
)

type Data struct {
	PageId int64
	Start  int64
	End    int64
}

type Index struct {
	Key  int64
	Data Data
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
	db := &DB{
		file: name,
		page: Page{
			dataF:  data,
			indexF: mmap,
		},
	}
	// initialization
	db.init()
	return db

}

func (db *DB) Close() {
	if err := db.page.dataF.Close(); err != nil {
		log.Fatal(err)
	}
	if err := db.page.indexF.Close(); err != nil {
		log.Fatal(err)
	}
}

// func (db *DB) init() {
// 	meta := IndexMeta{
// 		// Size:   0,
// 		MinKey: make([]byte, 0),
// 		MaxKey: make([]byte, 0),
// 	}
// 	db.page.insertIndex(meta, []Index{})
// }

func (db *DB) init() {
	meta := IndexMeta{
		MinKey: math.MaxInt64,
		MaxKey: math.MinInt64,
	}
	db.page.insertIndex(meta, []Index{})
}

// func findIndex(key []byte, indexs []Index) int {
// 	ind := sort.Search(len(indexs), func(i int) bool {
// 		// its like give me the smallest index where
// 		// index[i].key is not small than key
// 		return !(bytes.Compare(indexs[i].Key, key) == -1)
// 	})
// 	return ind
// }

func findIndex(key int64, indexs []Index) int {
	ind := sort.Search(len(indexs), func(i int) bool {
		// its like give me the smallest index where index[i].key is not small than key
		return !(indexs[i].Key < key)
	})
	return ind
}

func (db *DB) Get(key int64) []byte {
	// read index page
	_, indexs := db.page.getIndexPage()
	ind := findIndex(key, indexs)

	// find the data page from index page
	s, e := indexs[ind].Data.Start, indexs[ind].Data.End
	value := db.page.getDataMetaPage(int(s), int(e))
	return value
}

func (db *DB) Put(key int64, value []byte) error {
	if len(value) > int(MaxValueSize) {
		return errors.New("value size is large")
	}

	indM, indexs := db.page.getIndexPage()
	// updating the metadata
	indM.Size += 1
	indM.MaxKey = max(indM.MaxKey, key)
	indM.MinKey = min(indM.MaxKey, key)

	ind := findIndex(key, indexs)

	// insert into data
	s, e := db.page.insertIntoDataMetaPage(value)

	// insert into index
	if ind < len(indexs) && indexs[ind].Key == key {
		indexs = append(indexs[:ind+1], indexs[ind:]...)
	}

	if ind == len(indexs) {
		indexs = append(indexs, Index{})
	}

	indexs[ind] = Index{
		Key: key,
		Data: Data{
			PageId: 1,
			Start:  s,
			End:    e,
		},
	}

	db.page.insertIndex(indM, indexs)

	return nil
}

func (db *DB) Delete(key int64) error {
	return nil
}

// sorted slice as an index

// func max(a, b []byte) []byte {
// 	if bytes.Compare(a, b) > 0 {
// 		return a
// 	}
// 	return b
// }

// func min(a, b []byte) []byte {
// 	if bytes.Compare(a, b) > 0 {
// 		return b
// 	}
// 	return a
// }
