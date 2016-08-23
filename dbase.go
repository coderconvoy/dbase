package dbase

import (
	"fmt"
)

type lockMessage struct {
	doLock bool
	ls     string
	ch     chan bool
}

type lockQItem struct {
	m          lockMessage
	next, last *lockQItem
}

type Locker chan lockMessage

type DBase struct {
	root   string
	locker Locker
}

func beginLocker() Locker {
	ch := make(Locker)

	go func() {
		q := make(map[string]lockQItem, 0)
		locks := make(map[string]bool, 0)
		for {
			req := <-ch

			cl := l

		}

	}()
	return ch
}

func NewDB(root string) *DBase {
	return &DBase{root, beginLocker()}
}

func (db DBase) WriteMap(key string, val []byte) {
	_ = fmt.Sprintf("%x", []byte(key))
}

func (db DBase) ReadMap(key string) []byte {
	_ = fmt.Sprintf("%x", []byte(key))
	return nil
}
