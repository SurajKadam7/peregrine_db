package peregrinedb

import (
	"os"
	"reflect"
	"testing"
)

func Test_insertGetCombined(t *testing.T) {
	indM := IndexMeta{
		Size:   2,
		MaxKey: 20,
		MinKey: 10,
	}
	indexs := []Index{{Key: 10, PageId: 123, Ind: 1}, {Key: 20, PageId: 123, Ind: 1}}
	for i := 30; i < 1000; i++ {
		indexs = append(indexs, Index{Key: int64(i), PageId: 123, Ind: 1})
		indM.MaxKey = int64(i)
		indM.Size += 1
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
	mmap, _ := os.Create("testMmap")

	pg := Page{
		mmap: mmap,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg.insertIndex(tt.args.indM, tt.args.indexs)

			indM, indexs := pg.getIndexPage()

			if !reflect.DeepEqual(indM, tt.want.indM) {
				t.Errorf("insertIndex() = %v, want %v", indM, tt.want.indM)
			}
			if !reflect.DeepEqual(indexs, tt.want.indexs) {
				t.Errorf("insertIndex() = %v, want %v", indexs, tt.want.indexs)
			}

			// fmt.Printf("%+v\n", indM)
			// fmt.Printf("%+v\n", indexs)
		})
	}
}
