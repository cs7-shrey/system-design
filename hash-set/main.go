package main

import "fmt"

type Node struct {
	element string
	next *Node
}

type HashSet struct {
	size int
	container []*Node
}

func NewHashSet() *HashSet {
	return &HashSet{
		size: 0,
		container: make([]*Node, 16),
	}
}

func (hs *HashSet) Insert(x string) {
	// double calculation if Contains returns false
	if hs.Contains(x) {
		return
	}
	n := hs.size
	m := len(hs.container)

	if n+1 > m {
		hs.reOrg()
	}

	m = len(hs.container)
	loc := int(hash(x) % uint64(m))

	hs.size++

	head := hs.container[loc]
	node := &Node{element: x, next: head}
	hs.container[loc] = node
}

func (hs *HashSet) reOrg() {
	newM := 2 * len(hs.container)
	newCont := make([]*Node, newM)

	// loop over prev container
	for _, start := range hs.container {
		temp := start
		for temp != nil {
			newLoc := int(hash(temp.element) % uint64(newM))

			node := temp
			temp = temp.next

			head := newCont[newLoc]
			node.next = head
			newCont[newLoc] = node
		}
	}

	hs.container = newCont
}

func (hs *HashSet) Delete(x string) {
	mod := len(hs.container)
	loc := int(hash(x) % uint64(mod))

	head := hs.container[loc]
	if head == nil {
		return
	}

	// advance head
	if x == head.element {
		hs.container[loc] = head.next
		head.next = nil
		hs.size--
		return
	}

	prev := head
	temp := head.next
	
	for temp != nil {
		if x == temp.element {
			prev.next = temp.next
			temp.next = nil
			hs.size--
			return
		}
		prev = temp
		temp = temp.next
	}
}

func (hs *HashSet) Contains(x string) bool {
	mod := len(hs.container)
	loc := int(hash(x) % uint64(mod))

	head := hs.container[loc]
	if head == nil {
		return false
	}

	temp := head
	
	for temp != nil {
		if x == temp.element {
			return true
		}
		temp = temp.next
	}

	return false
}

func (hs *HashSet) Size() int {
	return hs.size
}

func hash(s string) uint64 {
	var h uint64 = 14695981039346656037

	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}

	return h
}

func main() {
	fmt.Println("Hello, World!")
}