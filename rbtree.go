package rbtree

import (
	"sync"
	"sync/atomic"
)

type Tree struct {
	root  *node
	wlock *sync.Mutex
}

type node struct {
	key         string
	keyHash     int64
	value       []byte
	left, right *node
}

func New() *Tree {
}

func (t *Tree) Insert(key string, value []byte) {
}

func (t *Tree) Delete(key string) {
}

func (t *Tree) Get(key string) ([]byte, bool) {
}
