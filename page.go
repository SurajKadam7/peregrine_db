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
	Size   int64
	MaxKey int64
	MinKey int64
}

type DataMeta struct {
	Size     int64
	DataSize int64
}

type Page struct {
	dataF *os.File
	mmap  *os.File
}

func (pg *Page) getDataMetaPage(index int) (data DataMeta, value []byte) {
	// TODo Read form file
	buff, err := io.ReadAll(pg.dataF)
	if err != nil {
		log.Fatal(err)
	}
	dataM := (*DataMeta)(unsafe.Pointer(&buff[0]))

	offSet := (int(MaxValueSize) * (index - 1)) + int(DataMetaSize)
	value = buff[offSet : offSet+MaxValueSize]

	return *dataM, value
}

func (pg *Page) insertIntoDataMetaPage(value []byte) {

	n, err := pg.dataF.Write(value)
	if err != nil {
		log.Fatal(err)
	}

	if n != len(value) {
		log.Fatal("size of byte not match with byte written")
	}
}

func (pg *Page) getIndexPage() (indMeta IndexMeta, indexs []Index) {
	pg.mmap, _ = os.Open("testMmap")
	buff, err := io.ReadAll(pg.mmap)
	if err != nil {
		log.Fatal(err)
	}

	indM := (*IndexMeta)(unsafe.Pointer(&buff[0]))
	ptrToIndexMeta := unsafe.Pointer(&buff[IndexMetaSize])
	indexs = unsafe.Slice((*Index)(ptrToIndexMeta), indM.Size)
	return *indM, indexs
}

func (pg *Page) insertIndex(indM IndexMeta, indexs []Index) {
	sz := len(indexs)*int(IndexSize) + int(IndexMetaSize)
	buff := make([]byte, sz)
	{
		buffIndM := (*IndexMeta)(unsafe.Pointer(&buff[0]))

		buffIndM.MaxKey = indM.MaxKey
		buffIndM.MinKey = indM.MinKey
		buffIndM.Size = indM.Size
	}

	for i, index := range indexs {
		ind := int(IndexMetaSize) + (i * int(IndexSize))
		buffIndex := (*Index)(unsafe.Pointer(&buff[ind]))

		buffIndex.Ind = index.Ind
		buffIndex.Key = index.Key
		buffIndex.PageId = index.PageId
	}

	pg.mmap.Truncate(0)
	n, err := pg.mmap.Write(buff)
	if err != nil {
		log.Fatal(err)
	}

	if n != len(buff) {
		log.Fatal("size not match")
	}

	log.Println("data written on index page : ", n, len(buff))

	pg.mmap.Sync()
	pg.mmap.Close()
}
