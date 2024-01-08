package peregrinedb

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestDB_Put(t *testing.T) {
	dataF, _ := os.OpenFile("dataF", os.O_RDWR|os.O_CREATE, 0666)
	indexF, _ := os.OpenFile("indexF", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)

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

	defer db.Close()

	type fields struct {
		db DB
	}
	type args struct {
		key   int64
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
				key:   i,
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
	dataF, err := os.OpenFile("dataF", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	indexF, err := os.OpenFile("indexF", os.O_RDONLY, 0666)
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

	defer db.Close()

	type args struct {
		key int64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	var start int64 = 10
	var limit int64 = 14
	for i := start; i < limit; i++ {
		tests = append(tests, struct {
			name string
			args args
			want []byte
		}{
			name: fmt.Sprint("test-", i),
			args: args{key: i},
			want: []byte(fmt.Sprint("raw_data_", i%10)),
		})

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := db.Get(tt.args.key)
			fmt.Printf("key : %v value : %v\n", tt.args.key, string(got))
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("DB.Get() = %v, want %v", got, tt.want)
			// }
		})
	}
}
