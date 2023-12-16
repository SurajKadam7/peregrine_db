package main

import (
	"fmt"
	"testing"
)

func Test_splitNode(t *testing.T) {
	type args struct {
		left  *node
		right *node
		ind   int
	}
	tests := []struct {
		name string
		args args
		want *node
	}{
		{
			name: "test",
			args: args{
				left: &node{
					cnt:        5,
					maxElement: 5,
					elements:   []element{{key: 20}, {key: 30}, {key: 40}, {key: 50}, {key: 60}},
				},
				ind:   3,
				right: initNode(true, 4),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitNode(tt.args.left, tt.args.right, tt.args.ind)
			left := tt.args.left
			right := tt.args.right

			fmt.Println(left.cnt)
			fmt.Println(left.elements)
			fmt.Println(right.cnt)
			fmt.Println(right.elements)
		})
	}
}

func Test_balance(t *testing.T) {
	type args struct {
		left  *node
		right *node
		ind   int
		key   int
	}
	tests := []struct {
		name        string
		args        args
		wantNewNode *node
	}{}

	for i := 0; i < 5; i++ {
		tests = append(tests, struct {
			name        string
			args        args
			wantNewNode *node
		}{
			name: fmt.Sprintf("test : %d", i),
			args: args{
				left: &node{
					isLeaf: true,
					cnt:        5,
					maxElement: 5,
					elements:   []element{{key: 20}, {key: 30}, {key: 40}, {key: 50},{key: 60}},
				},
				ind:   i,
				key:   5 + (10 * (i + 1)),
				right: initNode(true, 5),
			},
		})
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance(tt.args.left, tt.args.right, tt.args.ind, tt.args.key)
			left := tt.args.left
			right := tt.args.right

			fmt.Println(left.cnt)
			fmt.Println(left.elements)
			fmt.Println(right.cnt)
			fmt.Println(right.elements)
		})
	}
}
