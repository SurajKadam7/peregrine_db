package peregrinedb

import (
	"io"
	"log"
	"os"
	"unsafe"
)

type DBMeta struct {
	DataPage int64
}

const (
	IndexMetaSize = unsafe.Sizeof(IndexMeta{})
	IndexSize     = unsafe.Sizeof(Index{})
	DataMetaSize  = unsafe.Sizeof(DataMeta{})
	DataSize      = unsafe.Sizeof(Data{})
)

type IndexMeta struct {
	// Size   int64
	MaxKey []byte
	MinKey []byte
}

type DataMeta struct {
	Size     int64
	DataSize int64
}

type Page struct {
	dataF  *os.File
	indexF *os.File
}

func (pg *Page) getDataMetaPage(start, end int) (value []byte) {
	// TODo Read form file
	buff, err := io.ReadAll(pg.dataF)
	if err != nil {
		log.Fatal(err)
	}

	value = buff[start:end]

	return value
}

func (pg *Page) insertIntoDataMetaPage(value []byte) (start, end int64) {
	if len(value) > MaxValueSize {
		log.Fatal("max size limit of value is exceeded")
	}
	_, err := pg.dataF.Write(value)
	if err != nil {
		log.Fatal(err)
	}

	err = pg.dataF.Sync()
	if err != nil {
		log.Fatal(err)
	}

	fi, err := pg.dataF.Stat()
	if err != nil {
		log.Fatal("stats : ", err)
	}
	end = fi.Size()
	start = end - int64(len(value))
	return
}

func (pg *Page) getIndexPage() (indMeta IndexMeta, indexs []Index) {
	buff, err := io.ReadAll(pg.indexF)
	if err != nil {
		log.Fatal(err)
	}

	if len(buff) <= int(IndexMetaSize) {
		return indMeta, nil
	}

	indM := (*IndexMeta)(unsafe.Pointer(&buff[0]))
	// ptrToIndexMeta := unsafe.Pointer(&buff[IndexMetaSize])
	// indexs = unsafe.Slice((*Index)(ptrToIndexMeta), indM.Size)
	return *indM, indexs
}

func (pg *Page) insertIndex(indM IndexMeta, indexs []Index) {
	sz := len(indexs)*int(IndexSize) + int(IndexMetaSize)

	buff := make([]byte, sz)
	{
		buffIndM := (*IndexMeta)(unsafe.Pointer(&buff[0]))

		buffIndM.MaxKey = indM.MaxKey
		buffIndM.MinKey = indM.MinKey
		// buffIndM.Size = indM.Size
	}

	for i, index := range indexs {
		ind := int(IndexMetaSize) + (i * int(IndexSize))
		buffIndex := (*Index)(unsafe.Pointer(&buff[ind]))

		buffIndex.Key = index.Key
		buffIndex.PageId = index.PageId
		buffIndex.Start = index.Start
		buffIndex.End = index.End
	}

	pg.indexF.Truncate(0)
	n, err := pg.indexF.Write(buff)
	if err != nil {
		log.Fatal(err)
	}

	if n != len(buff) {
		log.Fatal("size not match")
	}

	log.Println("data written on index page : ", n, len(buff))

	err = pg.indexF.Sync()
	if err != nil {
		log.Fatal(err)
	}
	_, err = pg.indexF.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}
}
