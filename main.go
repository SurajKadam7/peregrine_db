package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"sort"
)

// B+ tree
// Node can have multiple sorted key value data
// Each leaf node is connected as a linked list
// Tree will grow in bottom up manner
// 1 4 8
// 0 3 7 10

func main() {
	var maxElement uint16 = 2
	t := initTree(maxElement)

	for i := 20; i < 1200; i += 10 {
		t.setInt(i, i*i)
	}

	for i := 10; i < 1200; i += 10 {
		t.getInt(i)
	}
	// t.setInt(20, 30)
	// t.setInt(30, 30)
	// t.setInt(40, 40)
	// t.setInt(50, 50)
	// t.setInt(60, 60)
	// t.setInt(70, 70)
	// t.setInt(80, 80)
	// t.setInt(90, 90)
	// t.setInt(100, 100)
	// t.setInt(110, 110)
	// t.setInt(120, 120)

	// getValues
	// t.getInt(20)
	// t.getInt(30)
	// t.getInt(40)
	// t.getInt(50)
	// t.getInt(60)
	// t.getInt(70)
	// t.getInt(80)
}

// 50 100
// 		30 50  		70 100   		110
// 20 30 | 40 50 | 60 70 | 80 90 | 100 110 | 120

// 90 170 |
// 50 90 | 130 170 | 210 250 
// 30 50 | 70 90 | 110 130| 150 170 | 190 210 | 230 250 | 270
// 20 30 | 40 50 | 60 70 | 80 90 | 100 110 | 120 130 | 140 150 | 160 170 | 180 190 | 200  210 | 220 230 | 240 250 \ 260 270 | 280 

type tree struct {
	root *node
}

type node struct {
	isLeaf     bool
	cnt        int       // count of element present
	maxElement uint16    // maximum element that each node can contain
	elements   []element // total elements
	parent     *node
	child      []*node
	next       *node // same level node for the leaf
}

type element struct {
	key   []byte
	value []byte
}

func initTree(maxElement uint16) *tree {
	n := &node{
		isLeaf:     true,
		cnt:        0,
		elements:   make([]element, 0),
		next:       nil,
		maxElement: maxElement,
	}
	t := &tree{
		root: n,
	}
	return t
}

// will give me the leaf node with the index
func findLeafNode(t *tree, key []byte) (node *node, ind int) {
	// recursion is like go to the right if small node find go down
	node = t.root
	for node != nil {
		ind = sort.Search(node.cnt, func(i int) bool {
			return bytes.Compare(node.elements[i].key, key) != -1
		})

		if node.isLeaf {
			return node, ind
		}
		// will remove this if condition
		if ind == node.cnt {
			// we can use the node.cnt as ind for nodeRef as it will always node elms + 1
			node = node.child[ind]
			continue
		}

		node = node.child[ind]
	}
	return node, ind
}

func initNode(isLeaf bool, maxElement uint16) *node {
	return &node{
		isLeaf:     isLeaf, // is this the best way to set leaf
		cnt:        0,
		maxElement: maxElement,
		elements:   make([]element, 0),
	}
}

func insertElement(n *node, ind int, key []byte) {
	e := element{
		key: key,
	}

	temp := make([]element, ind)
	copy(temp, n.elements[:ind])

	temp = append(temp, e) // temp = [39]
	temp = append(temp, n.elements[ind:]...)

	n.elements = temp
	n.cnt += 1
}

func splitNode(n *node, ind int) *node {
	mid := (n.maxElement / 2) + 1 //2
	newNode := initNode(n.isLeaf, n.maxElement)

	newNode.elements = n.elements[mid:] // 40
	n.elements = n.elements[:mid]       // 20 30

	n.cnt = int(mid)
	newNode.cnt = int(newNode.maxElement - mid)
	return newNode
}

func balance(n *node, ind int, key []byte) (newNode *node) {
	newNode = splitNode(n, ind)

	if newNode.isLeaf {
		// linking the leaf nodes together.
		n.next, newNode.next = newNode, n.next
	}

	// inserting the newNode
	if n.cnt > ind {
		insertElement(n, ind, key)
	} else {
		insertElement(newNode, ind-n.cnt, key) // 2-2, 0
	}
	return newNode
}

func (t *tree) insert(n *node, ind int, key []byte) {
	firstHalf := n
	for {
		if firstHalf.cnt < int(firstHalf.maxElement) && ind < int(firstHalf.maxElement) {
			insertElement(firstHalf, ind, key)
			return
		}

		secondHalf := balance(firstHalf, ind, key)

		if !firstHalf.isLeaf {
			// remove the last child
			firstHalf.child[firstHalf.cnt].parent = secondHalf
			secondHalf.child = append(secondHalf.child, firstHalf.child[firstHalf.cnt])

			firstHalf.child[firstHalf.cnt+1].parent = secondHalf
			secondHalf.child = append(secondHalf.child, firstHalf.child[firstHalf.cnt+1])

			firstHalf.child = firstHalf.child[:firstHalf.cnt]
		}

		if firstHalf.parent == nil {
			// changing the root node and creating a new parent
			temp := initNode(false, firstHalf.maxElement)
			temp.elements = append(temp.elements, firstHalf.elements[firstHalf.cnt-1])
			temp.cnt += 1

			temp.child = append(temp.child, firstHalf, secondHalf)
			firstHalf.parent = temp
			secondHalf.parent = temp
			t.root = temp
			return
		}

		// adding parent and child relation
		secondHalf.parent = firstHalf.parent
		firstHalf.parent.child = append(firstHalf.parent.child, secondHalf)

		// addition of the left key will be take care by the below key update
		key = firstHalf.elements[firstHalf.cnt-1].key
		ind = sort.Search(firstHalf.parent.cnt, func(i int) bool {
			return bytes.Compare(firstHalf.parent.elements[i].key, key) != -1
		})
		// calling recursivly for the parent of firstHalf
		firstHalf = firstHalf.parent
	}
}

func (t *tree) insertData(n *node, ind int, key []byte, value []byte) {
	t.insert(n, ind, key)
	// locking needs to be handled for the future
	if ind < n.cnt {
		n.elements[ind].value = value
		return
	}
	if n.next.elements == nil {
		log.Fatal("not able to insert the element")
	}
	// TODO need to handle carefully for index out of range
	n.next.elements[ind-n.cnt].value = value
}

func (t *tree) set(key []byte, value []byte) {
	// this function will alway leaf insert the data
	// some findLeaf function will find the leafNode to insert the data
	// if tree structure need to grow will use grow function which will grow the tree

	leafNode, ind := findLeafNode(t, key)
	t.insertData(leafNode, ind, key, value)
}

func (t *tree) get(key []byte) []byte {
	leaf, ind := findLeafNode(t, key)
	if ind >= leaf.cnt {
		fmt.Println("not found")
		return []byte{}
	}
	fmt.Println(leaf.elements[ind].value)
	return leaf.elements[ind].value
}

func (t *tree) setInt(a, b int) {
	a1 := make([]byte, 4)
	b1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(a1, uint32(a))
	binary.LittleEndian.PutUint32(b1, uint32(b))
	t.set(a1, b1)
}

func (t *tree) getInt(a int) {
	a1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(a1, uint32(a))
	t.get(a1)
}
