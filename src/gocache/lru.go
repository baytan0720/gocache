package gocache

import "fmt"

type lrulist struct {
	head *node
	tail *node
}

type node struct {
	pre  *node
	next *node
	data interface{}
}

func makeLrulist() *lrulist {
	head := new(node)
	tail := new(node)
	head.next = tail
	tail.pre = head
	return &lrulist{
		head: head,
		tail: tail,
	}
}

func (l *lrulist) insertHead(node *node) {
	l.head.next.pre = node
	node.next = l.head.next
	l.head.next = node
	node.pre = l.head
}

func (l *lrulist) insertTail(node *node) {
	l.tail.pre.next = node
	node.pre = l.head.next
	l.tail.pre = node
	node.next = l.head
}

func (l *lrulist) delete(node *node) {
	node.pre.next = node.next
	node.next.pre = node.pre
}

func (l *lrulist) deletelast() (key interface{}) {
	key = l.tail.pre.data
	l.tail.pre.pre.next = l.tail
	l.tail.pre = l.tail.pre.pre
	return
}

func (l *lrulist) Show() {
	node := l.head.next
	for node != l.tail {
		fmt.Print(node.data)
		fmt.Print(" ")
		node = node.next
	}
	fmt.Println("")
}
