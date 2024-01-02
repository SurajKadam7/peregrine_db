package main

import (
	"fmt"
	"log"
	"sort"
)

type tree struct {
	root *node
}

type node struct {
	isLeaf     bool
	cnt        int      // count of element present
	maxElement uint16   // maximum element that each node can contain
	elements   elements // total elements
	parent     *node
	child      []*node
	next       *node // same level node for the leaf
	degree     int
}

type element struct {
	key   int
	value float64
}

func New(maxElement uint16) *tree {
	n := &node{
		isLeaf:     true,
		cnt:        0,
		elements:   make([]element, 0),
		next:       nil,
		degree:     1,
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
	leafNode, ind := findLeafNode(t, key)
	t.insertData(leafNode, ind, key, value)
}

func removeElement(n *node, ind int) {
	n.elements = append(n.elements[:ind], n.elements[ind+1:]...)
	n.cnt -= 1
}

func hasSpace(child []*node, ind int, maxElement uint16) bool {
	return child[ind].cnt <= int(maxElement)/2
}

func leftMerge(n *node, ind int) {
	leftChild := n.child[ind-1]
	rightChild := n.child[ind]

	leftChild.elements = append(leftChild.elements, rightChild.elements...)
	leftChild.cnt += rightChild.cnt

	leftChild.child = append(leftChild.child, rightChild.child...)

	index := leftChild.cnt - 1
	n.elements[ind-1] = leftChild.elements[index]

	// delete the node
	if leftChild.isLeaf {
		leftChild.next = rightChild.next
	}

	n.child = append(n.child[:ind], n.child[ind+1:]...)
}

func rightMerge(n *node, ind int) {
	rightChild := n.child[ind+1]
	leftChild := n.child[ind]

	rightChild.elements = append(leftChild.elements, rightChild.elements...)
	rightChild.cnt += leftChild.cnt

	// coping the child elements
	rightChild.child = append(leftChild.child, rightChild.child...)

	if rightChild.isLeaf && ind-1 > 0 {
		n.child[ind-1].next = leftChild.next
	}

	n.child = append(n.child[:ind], n.child[ind+1:]...)
}

func mergeIfPossible(n *node, ind int) {
	if n == nil || n.isLeaf {
		return
	}

	// check for the left merge
	if ind-1 >= 0 && hasSpace(n.child, ind-1, n.maxElement) {
		leftMerge(n, ind)
		return
	}

	// checking for the right merge
	if ind+1 < len(n.child) && hasSpace(n.child, ind+1, n.maxElement) {
		rightMerge(n, ind)
		return
	}

	// complete child node is deleted
	if n.child[ind].cnt == 0 {
		if ind-1 > 0 {
			n.child[ind-1].next = n.child[ind].next
		}
		n.child = append(n.child[:ind], n.child[ind+1:]...)
	}
}

func handleNilParent(t *tree, n *node) {
	if n.cnt == 0 && len(n.child) == 1 {
		t.root = n.child[0]
		return
	}

	if n.cnt == 0 {
		i := n.child[0].cnt - 1
		n.elements = append(n.elements, n.child[0].elements[i])
		n.cnt += 1
	}
}

func (t *tree) deleteData(n *node, ind int, key int) {
	curr := n
	for {
		removeElement(curr, ind)
		mergeIfPossible(curr, ind)

		if curr.parent == nil {
			handleNilParent(t, curr)
			return
		}

		curr = curr.parent

		newInd := sort.Search(n.parent.cnt, func(i int) bool {
			return curr.elements[i].key >= key
		})

		if newInd == curr.cnt || curr.elements[newInd].key != key {
			mergeIfPossible(curr, ind)
			return
		}
		ind = newInd
	}

}

func (t *tree) delete(key int) {
	leafNode, ind := findLeafNode(t, key)
	t.deleteData(leafNode, ind, key)
}

func (t *tree) get(key int) float64 {
	leaf, ind := findLeafNode(t, key)
	if ind >= leaf.cnt || leaf.elements[ind].key != key {
		return 0
	}
	return leaf.elements[ind].value
}

func (t *tree) getRange(key int, limit int) []element {
	leaf, ind := findLeafNode(t, key)
	if ind >= leaf.cnt || leaf.elements[ind].key != key {
		return []element{}
	}
	elements := make([]element, limit)
	for i := 0; i < limit; i++ {
		elements = append(elements, leaf.elements[ind])
		ind += 1
		if ind == leaf.cnt {
			leaf = leaf.next
			ind = 0
			if leaf == nil {
				return elements
			}
		}
	}
	return elements
}

type elements []element

func (e *elements) removeAt(ind int) {
	copy((*e)[ind:], (*e)[ind+1:])
	(*e) = (*e)[:len(*e)-1]
}

func leftStill(n *node, ind int, m int) {
	// this many I can still
	l := m
	m = n.child[ind-1].cnt - m
	e := make(elements, l+n.child[ind].cnt)
	c := make([]*node, l+n.child[ind].cnt+1)

	copy(e, n.child[ind-1].elements[m:])
	copy(e[l:], n.child[ind].elements)
	n.child[ind].elements = e

	n.child[ind-1].elements = n.child[ind-1].elements[:m]

	if n.child[ind].child != nil {
		copy(c, n.child[ind-1].child[m:])
		copy(c[l:], n.child[ind].child)
		if c[len(c)-1] == nil {
			c = c[:len(c)-1]
		}
		n.child[ind].child = c
		n.child[ind-1].child = n.child[ind-1].child[:m]
	}
}

func rightStill(n *node, ind int, m int) {
	n.child[ind].elements = append(n.child[ind].elements, n.child[ind+1].elements[:m]...)
	n.child[ind+1].elements = n.child[ind+1].elements[m:]

	if n.child[ind].child != nil {
		n.child[ind].child = append(n.child[ind].child, n.child[ind+1].child[:m]...)
		n.child[ind+1].child = n.child[ind+1].child[m:]
	}
}

func (t *tree) deleteV2(key int) {
	if t == nil || t.root == nil {
		fmt.Println("tree is empty")
		return
	}
	leaf, i := findLeafNode(t, key)
	if i >= leaf.cnt || leaf.elements[i].key != key {
		return
	}

	removeV2(leaf, i, key)
	if t.root.cnt <= 1 && len(t.root.child) == 1 {
		t.root = t.root.child[0]
		t.root.parent = nil
	} else if t.root.cnt == 0 && len(t.root.child) > 1 {
		panic("something is wrong with the tree")
	} else if t.root.cnt == 0 {
		// tree is empty now
		t.root = nil
	}
}

func findChildNodePosition(parent *node, key int) int {
	if parent == nil {
		return 0
	}
	// need to find first node where my child was
	// index is second last node of node whose index is higher that key
	i := sort.Search(len(parent.child), func(i int) bool {
		ind := parent.child[i].cnt
		return parent.child[i].elements[ind-1].key > key
	})
	if i == 0 {
		i++
	}
	return i - 1
}

// growChildAndRemove will take currentNode, key, and index of the key in the currentNode
// it will balance tree and returns
func growChildAndRemoveV2(n *node, keyInd int, key int) {
	min := getMinNodes(n)

	parent := n.parent
	i := findChildNodePosition(parent, key)
	if n.cnt > 0 && n.elements[keyInd].key == key {
		n.elements.removeAt(keyInd)
		n.cnt -= 1
	}

	if parent == nil {
		return
	}

	if n.cnt == 0 {
		// removing the node from the parent childs
		copy(parent.child[i:], parent.child[i+1:])
		parent.child = parent.child[:len(parent.child)-1]
		if len(n.child) > 0{
			n.child[0].parent = parent.child[i]
			parent.child[i].child = append(parent.child[i].child , n.child[0])
		}
		return
	}

	// swapping the parent element if I have deleted a rightmost element in
	if parent.child[i].cnt == keyInd {
		pI := i
		if i == parent.cnt { /* when deleted elem from rightmost child */
			pI -= 1
		}
		n.parent.elements[pI] = n.parent.child[i].elements[keyInd-1]
	}

	// TODO : need to change this condtions
	if n.cnt >= min {
		return
	}

	switch {
	case i > 0 && parent.child[i-1].cnt > min:
		still := parent.child[i-1].cnt - min
		leftStill(parent, i, min)
		parent.child[i].cnt += still
		parent.child[i-1].cnt -= still
		// changing the parent element with the current last element in the child
		parent.elements[i-1] = parent.child[i-1].elements[parent.child[i-1].cnt-1]

	case i+1 < len(parent.child) && parent.child[i+1].cnt > min:
		still := parent.child[i+1].cnt - min
		rightStill(parent, i, still)
		parent.child[i].cnt += still
		parent.child[i+1].cnt -= still

		parent.elements[i] = parent.child[i].elements[parent.child[i].cnt-1]

	default:
		// include the merge logic
		if i+1 >= len(parent.child) {
			i -= 1
		}
		// copying all the values from current node to the right node and distroy self
		leftStill(parent, i+1, parent.child[i].cnt)
		if parent.child[i].child != nil {
			for _, cc := range parent.child[i+1].child {
				if cc.parent == parent.child[i+1] {
					break
				}
				cc.parent = parent.child[i+1]
			}
		}
		parent.child[i+1].cnt += parent.child[i].cnt
		// distroying the node
		// removing
		copy(parent.child[i:], parent.child[i+1:])
		parent.child = parent.child[:len(parent.child)-1]

		// deleting left element of current index

		parent.elements.removeAt(i)
		parent.cnt -= 1

	}
}

func getMinNodes(n *node) int {
	return n.degree - 1
}

func removeV2(n *node, ind int, key int) {
	for {
		// delete the element

		// grow treee
		growChildAndRemoveV2(n, ind, key)

		// below flow only work when we removed right most element from node
		// need to go up
		if n.parent == nil || n.parent.cnt == 0 {
			return
		}
		n = n.parent
		// finding the index where my key is prenset in the elements
		ind = sort.Search(n.cnt, func(i int) bool {
			return n.elements[i].key > key
		})
		if ind == 0 {
			ind += 1
		}
		ind -= 1
	}
}
