package main

import (
	"fmt"
	"strings"
)

func printNode(node *node) string {
	r := make([]string, 0)
	for _, val := range node.elements {
		r = append(r, fmt.Sprintf("%v ", val.key))
	}
	return strings.Join(r, "")
}

func main() {
	tree := createTree()
	display(tree)
}

func getSpaces(n int) string {
	spaces := make([]string, 0)
	for i := 0; i < n; i++ {
		spaces = append(spaces, " ")
	}
	return strings.Join(spaces, "")
}

func display(t *tree) {
	if t.root == nil{
		return 
	}
	n := t.root
	stack := make([][]*node, 0)

	stack = append(stack, []*node{n})
	bigStack := make([][]*node, 0)

	for len(stack) > 0 {
		curr := stack[0]
		bigStack = append(bigStack, curr)
		stack = stack[1:]
		currNode := make([]*node, 0)
		for _, s := range curr {
			if s.child == nil {
				continue
			}
			currNode = append(currNode, s.child...)
		}
		if len(currNode) > 0 {
			stack = append(stack, currNode)
		}
	}

	result := make([]string, 0)
	for len(bigStack) > 0 {
		curr := bigStack[0]
		bigStack = bigStack[1:]
		ss := []string{}
		for i, s := range curr {
			ss = append(ss, printNode(s))  
			if i < len(curr)-1 {
				ss = append(ss, " | ")
			}
			
		}
		sss := strings.Join(ss, "")
		result = append(result, sss)
	}

	for _, val := range result{
		fmt.Println(val)
	}
}

// 20 40
// 10 20 | 30 40 | 50

func createTree() *tree {
	var maxElement uint16 = 2
	t := New(maxElement)
	t.set(30, 0.3)
	t.set(20, 0.2)
	t.set(10, 0.1)
	t.set(5, 0.05)
	t.set(3, 0.03)
	t.set(7, 0.07)
	return t
}
