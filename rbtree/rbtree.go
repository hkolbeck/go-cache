package rbtree

import (
	"sync"
	"sync/atomic"
)

type Tree struct {
	root  *node
	wlock *sync.Mutex
	epoch uint64
}

type reader struct {
	uint64 epoch
	reader *next
}

type node struct {
	key                 string
	keyHash             int64
	value               []byte
	left, right, parent *node
	color               int
}

const (
	RED = iota
	BLACK
)

func New() *Tree {
}

func (t *Tree) Put(key string, value []byte) {
	defer c.endWrite(c.startWrite())
}

func (t *Tree) Get(key string) ([]byte, bool) {
}

func (t *Tree) Delete(key string) {
	defer t.endWrite(t.startWrite())
}

/* Tree helper functions */
func (t *Tree) next(n *node) {
}

func (t *Tree) prev(n *node) {
}
