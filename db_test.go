package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

var maxElement = 2

func genRand(a []int) []int {
	temp := make([]int, len(a))
	copy(temp, a)
	var randomArray []int
	for i := 0; i < 5; i++ { // Change 5 to the desired length of the new array
		randomIndex := rand.Intn(len(temp))
		randomArray = append(randomArray, temp[randomIndex])
		temp[randomIndex], temp[len(temp)-1] = temp[len(temp)-1], temp[randomIndex]
		temp = temp[:len(temp)-1]
	}
	return randomArray
}

func resArr(val string) []int {
	a := strings.Split(val, " ")
	res := []int{}
	for _, key := range a {
		n, _ := strconv.Atoi(key)
		res = append(res, n)
	}
	return res
}

func initTree(insertElm []int) *tree {
	t := New(uint16(maxElement))
	for _, key := range insertElm {
		t.set(key, float64(key)/10)
	}
	display(t)
	fmt.Printf("---------------tree---------------\n\n")
	return t
}
func deleteElements(t *tree, deleteElm []int) {
	for _, key := range deleteElm {
		t.deleteV2(key)
		display(t)
		if t.root != nil {
			fmt.Printf("------------------------------\n")
		}
	}
}
func Test_separate(t1 *testing.T) {
	// [30 10 25 20 15] [20 10 25 30 15]
	insert := "30 10 25 20 15"
	delete := "20 10 25 30 15"
	t := initTree(resArr(insert))
	eld := resArr(delete)
	deleteElements(t, eld)
}

func Test_tree_delete(t1 *testing.T) {
	// 1st
	tempValues := []int{10, 20, 15, 25, 30}

	for i := 0; i < 100; i++ {
		insert := genRand(tempValues)
		deleteElm := genRand(tempValues)
		fmt.Println(insert, deleteElm)

		t := initTree(insert)
		deleteElements(t, deleteElm)
		fmt.Printf("------------   done  ------------------\n\n")
	}
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

// 10
// 5 10  | 20 only copy 7 as a value
// 3 5  | 7 10  | 20  | 30
// ------------------------------
// 10
// 7 10  | 20
// 5 7  | 10  | 20  | 30
// ------------------------------
// ------------------------------
// 10 20
// 5 10  | 20  | 30
