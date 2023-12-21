package main

import (
	"fmt"
	"testing"
)

func Test_tree_delete(t1 *testing.T) {
	// 1st
	var maxElement uint16 = 2
	// t := New(maxElement)
	// t.set(30, 0.3)
	// t.set(20, 0.2)
	// t.set(10, 0.1)
	// t.set(8, 0.08)
	// t.delete(10) // pass

	// // 2nd
	// t = New(maxElement)
	// t.set(30, 0.3)
	// t.set(20, 0.2)
	// t.set(10, 0.1)
	// t.set(8, 0.08)
	// t.set(7, 0.07)
	// t.delete(10) // pass rightShif + parnetChange

	t := New(maxElement)
	t.set(30, 0.3)
	t.set(20, 0.2)
	t.set(10, 0.1)
	t.set(5, 0.05)
	t.set(3, 0.03)
	t.set(7, 0.07) // to avoid merging with the deleting 3 operation
	display(t)
	fmt.Println("------------------------------")
	t.deleteV2(3)
	display(t)
	fmt.Println("------------------------------")
	t.deleteV2(7)
	// t.deleteV2(10) // pass leftShif + parnetChange
	fmt.Println("------------------------------")
	// t.deleteV2(10)
	display(t)
}

// 10
// 5 10  | 20
// 3 5  | 7 10  | 20  | 30

// 10
// 7 10  | 20
// 5  7 |  10  | 20  | 30

// 7
// 3 7  | 20
//    3 | 5 7 | 20  | 30

// delte 3
// 10
// 7 10  | 20
// 5 10  | 20  | 30

// merge logic
// 10  | 20
// 5 10  | 20  | 30

// 20
// 5 20  | 30
