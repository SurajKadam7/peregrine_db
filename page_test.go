package peregrinedb

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func Test_insertGetCombined(t *testing.T) {
	indM := IndexMeta{
		// Size:   0,
		MaxKey: []byte("index_10"),
		MinKey: []byte("index_10"),
	}

	indexs := []Index{}
	for i := 30; i < 35; i++ {
		indexs = append(indexs, Index{Key: []byte(fmt.Sprint("index_", i)), /*PageId: 123, Start: 1*/})
		indM.MaxKey = []byte(fmt.Sprint("index_", i))
		// indM.Size += 1s
	}
	type args struct {
		indM   IndexMeta
		indexs []Index
	}

	type want struct {
		indM   IndexMeta
		indexs []Index
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test1",
			args: args{indM: indM, indexs: indexs},
			want: want{indM: indM, indexs: indexs},
		},
	}
	// empty initialization
	mmap, err := os.CreateTemp("./", "indexFile")
	if err != nil {
		log.Fatal(err)
	}

	// os.Remove(mmap.Name())

	pg := Page{
		indexF: mmap,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg.insertIndex(tt.args.indM, tt.args.indexs)
			pg.indexF.Close()
			// indM, indexs := pg.getIndexPage()

			// if !reflect.DeepEqual(indM, tt.want.indM) {
			// 	t.Errorf("insertIndex() = %v, want %v", indM, tt.want.indM)
			// }
			// if !reflect.DeepEqual(indexs, tt.want.indexs) {
			// 	t.Errorf("insertIndex() = %v, want %v", indexs, tt.want.indexs)
			// }

			// fmt.Printf("%+v\n", indM)
			// fmt.Printf("%+v\n", indexs)
		})
	}
}

func TestPage_insertIntoDataMetaPage(t *testing.T) {
	f, _ := os.CreateTemp(".", "dataFile")
	// defer os.ReadFile(f.Name())
	type fields struct {
		dataF  *os.File
		indexF *os.File
	}
	type args struct {
		value []byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantStart int64
		wantEnd   int64
	}{}
	for i := 10; i < 15; i++ {
		tests = append(tests, struct {
			name      string
			fields    fields
			args      args
			wantStart int64
			wantEnd   int64
		}{
			name:   fmt.Sprint("test-", i%10),
			fields: fields{dataF: f},
			args:   args{value: []byte(fmt.Sprint("data_raw_", i))},
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &Page{
				dataF:  tt.fields.dataF,
				indexF: tt.fields.indexF,
			}
			pg.insertIntoDataMetaPage(tt.args.value)
			// gotStart, gotEnd := pg.insertIntoDataMetaPage(tt.args.value)
			// if gotStart != tt.wantStart {
			// 	t.Errorf("Page.insertIntoDataMetaPage() gotStart = %v, want %v", gotStart, tt.wantStart)
			// }
			// if gotEnd != tt.wantEnd {
			// 	t.Errorf("Page.insertIntoDataMetaPage() gotEnd = %v, want %v", gotEnd, tt.wantEnd)
			// }
		})
	}
}
