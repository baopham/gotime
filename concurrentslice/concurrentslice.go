package concurrentslice

// Credit: https://github.com/dnaeon/gru/blob/master/utils/slice.go

import (
	"sync"
)

type Slice struct {
	mux   sync.RWMutex
	items []interface{}
}

func New() *Slice {
	return &Slice{
		items: make([]interface{}, 0),
	}
}

func (cs *Slice) Append(item interface{}) {
	cs.mux.Lock()
	cs.items = append(cs.items, item)
	cs.mux.Unlock()
}

func (cs *Slice) Size() int {
	return len(cs.items)
}

func (cs *Slice) Items() []interface{} {
	return cs.items
}

func (cs *Slice) Fill(f func(i int, v interface{})) {
	for i, v := range cs.Items() {
		f(i, v)
	}
}
