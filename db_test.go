package peregrinedb

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestDB_Put(t *testing.T) {
	dataF, _ := os.Create("dataF")
	indexF, _ := os.Create("indexF")

	// cleanning
	// defer os.Remove(dataF.Name())
	// defer os.Remove(indexF.Name())

	db := DB{
		file: "testFile",
		page: Page{
			dataF:  dataF,
			indexF: indexF,
		},
	}

	type fields struct {
		db DB
	}
	type args struct {
		key   []byte
		value []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{}

	var limit int64 = 14
	var start int64 = 10

	for i := start; i < limit; i++ {
		test := struct {
			name    string
			fields  fields
			args    args
			wantErr bool
		}{
			name: fmt.Sprint("test-", i%10),
			fields: fields{
				db: db,
			},

			args: args{
				key:   []byte(fmt.Sprint("", i)),
				value: []byte(fmt.Sprint("raw_data_", i%10)),
			},
		}
		tests = append(tests, test)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Put(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestDB_Get(t *testing.T) {
	dataF, err := os.Open("dataF")
	if err != nil {
		log.Fatal(err)
	}
	indexF, err := os.Open("indexF")
	if err != nil {
		log.Fatal(err)
	}

	// cleanning
	// defer os.Remove(dataF.Name())
	// defer os.Remove(indexF.Name())

	db := DB{
		file: "testFile",
		page: Page{
			dataF:  dataF,
			indexF: indexF,
		},
	}

	type args struct {
		key []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	var start int64 = 10
	var limit int64 = 12
	for i := start; i < limit; i++ {
		tests = append(tests, struct {
			name string
			args args
			want []byte
		}{
			name: fmt.Sprint("test-", i),
			args: args{key: []byte(fmt.Sprint("", i))},
			want: []byte(fmt.Sprint("raw_data_", i%10)),
		})

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := db.Get(tt.args.key)
			fmt.Println(string(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
