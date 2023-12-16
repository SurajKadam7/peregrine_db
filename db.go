package main

import (
	"log"
	"sort"
)

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
	key   int
	value float64
}

// 
func New(maxElement uint16) *tree {
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
func findLeafNode(t *tree, key int) (node *node, ind int) {
	// recursion is like go to the right if small node find go down
	node = t.root
	for node != nil {
		iterationRequireTofind[key] += 1
		ind = sort.Search(node.cnt, func(i int) bool {
			return node.elements[i].key >= key
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

func insertElement(n *node, ind int, key int) {
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

func copyChilds(left *node, right *node) {
	if len(left.child) > len(right.child) {
		for i := left.cnt; i < len(left.child); i++ {
			left.child[i].parent = right
			right.child = append(right.child, left.child[i])
		}
		left.child = left.child[:left.cnt]
		return
	}

	for i := 0; i < left.cnt; i++ {
		right.child[i].parent = left
		left.child = append(left.child, right.child[i])
	}
	right.child = right.child[left.cnt:]
}

func splitNode(n *node, right *node, ind int) {
	mid := (n.maxElement / 2) + 1 //2
	if int(mid) > ind {
		mid--
	}

	// decide who will be the left and who will be the right
	temp := n.elements
	n.elements, right.elements = temp[:mid], temp[mid:] // 5 20 30 ind = 2

	n.cnt = int(mid)
	right.cnt = int(right.maxElement - mid)
	// child split also need to be cone
	// it is alway like leftNode will copy only len(elemenet) childs

}

func balance(left *node, right *node, ind int, key int) {
	splitNode(left, right, ind)

	if right.isLeaf {
		// linking the leaf nodes together.
		left.next, right.next = right, left.next
	}

	// inserting the newNode following <= due to left bias
	if left.cnt <= right.cnt {
		insertElement(left, ind, key)
	} else {
		insertElement(right, ind-left.cnt, key) // 2-2, 0
	}

	if left.isLeaf {
		return
	}

	if right.elements[right.cnt-1].key < left.elements[left.cnt-1].key {
		right, left = left, right
	}

	copyChilds(left, right)
}

func (t *tree) insert(n *node, ind int, key int) {
	left := n
	for {
		if left.cnt < int(left.maxElement) && ind < int(left.maxElement) {
			insertElement(left, ind, key)
			return
		}
		right := initNode(left.isLeaf, left.maxElement)
		balance(left, right, ind, key)

		if left.parent == nil {
			// changing the root node and creating a new parent
			temp := initNode(false, left.maxElement)
			temp.elements = append(temp.elements, left.elements[left.cnt-1])
			temp.cnt += 1

			temp.child = append(temp.child, left, right)
			left.parent = temp
			right.parent = temp
			t.root = temp
			return
		}

		// adding parent and child relation
		right.parent = left.parent
		left.parent.child = append(left.parent.child, right)

		sort.Slice(left.parent.child, func(i, j int) bool {
			return left.parent.child[i].elements[0].key < left.parent.child[j].elements[0].key
		})

		// addition of the left key will be take care by the below key update
		key = left.elements[left.cnt-1].key
		ind = sort.Search(left.parent.cnt, func(i int) bool {
			return left.parent.elements[i].key >= key
		})
		// calling recursivly for the parent of firstHalf
		left = left.parent
	}
}

func (t *tree) insertData(n *node, ind int, key int, value float64) {
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

func (t *tree) set(key int, value float64) {
	// this function will alway leaf insert the data
	// some findLeaf function will find the leafNode to insert the data
	// if tree structure need to grow will use grow function which will grow the tree

	leafNode, ind := findLeafNode(t, key)
	t.insertData(leafNode, ind, key, value)
}

func (t *tree) get(key int) float64 {
	leaf, ind := findLeafNode(t, key)
	if ind >= leaf.cnt || leaf.elements[ind].key != key {
		return 0
	}
	return leaf.elements[ind].value
}
