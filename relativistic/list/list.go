package list

import (
	rp "../../relativistic"
)

type List struct {
	rel        *rp.Relativistic
	head, tail *node
}

type node struct {
	value      []byte
	prev, next *node
}

func New(queueSize int) *List {
	return &List{
		rel: rp.New(queueSize),
	}
}

func (l *List) InsertAtTail(val []byte) {
	l.rel.Lock()
	defer l.rel.Unlock()

	newNode := l.alloc(val)
	newNode.prev = l.tail

	if l.tail == nil {
		l.tail = newNode
		l.head = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
}

func (l *List) InsertAtHead(val []byte) {
	l.rel.Lock()
	defer l.rel.Unlock()

	newNode := l.alloc(val)
	newNode.next = l.head

	if l.head == nil {
		l.tail = newNode
		l.head = newNode
	} else {
		l.head.prev = newNode
		l.head = newNode
	}
}

func (l *List) RemoveHead() ([]byte, bool) {
	l.rel.Lock()
	defer l.rel.Unlock()

	if l.head == nil {
		return nil, false
	}

	val := l.head.value

	if l.head == l.tail {
		l.head = nil
		l.tail = nil
	} else {
		oldHead := l.head
		l.head = l.head.next
		l.head.prev = nil
		l.rel.WaitForReaders()
		l.free(oldHead)
	}

	return val, true
}

func (l *List) RemoveTail() ([]byte, bool) {
	l.rel.Lock()
	defer l.rel.Unlock()

	if l.tail == nil {
		return nil, false
	}

	val := l.tail.value

	if l.head == l.tail {
		l.head = nil
		l.tail = nil
	} else {
		oldTail := l.tail
		l.tail = l.tail.prev
		l.tail.next = nil
		l.rel.WaitForReaders()
		l.free(oldTail)
	}

	return val, true
}

func (l *List) PeekHead() ([]byte, bool) {
	defer l.rel.EndRead(l.rel.StartRead())

	head := l.head

	if head == nil {
		return nil, false
	}

	return head.value, true
}

func (l *List) PeekTail() ([]byte, bool) {
	defer l.rel.EndRead(l.rel.StartRead())

	tail := l.tail

	if tail == nil {
		return nil, false
	}

	return tail.value, true
}

//Terrible
func (l *List) Snapshot() [][]byte {
	defer l.rel.EndRead(l.rel.StartRead())

	arr := make([][]byte, 100, 0)

	for pos := l.head; pos != nil; pos = pos.next {
		arr = append(arr, pos.value)
	}

	return arr
}

func (l *List) free(n *node) {
	//Just for now
	n.value = nil
}

func (l *List) alloc(val []byte) *node {
	return &node{value: val}
}
