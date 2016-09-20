package dbase

import (
	"encoding/hex"
	"os"
	"path"
)

type lockMessage struct {
	k  string
	ch chan bool
}

type lockQItem struct {
	ch         chan bool
	next, last *lockQItem
}

type Locker chan lockMessage

type DBase string

func beginLocker() Locker {
	ch := make(Locker)

	go func() {
		q := make(map[string]*lockQItem, 0)
		for {
			req := <-ch
			qTop, ok := q[req.k]

			if req.ch != nil {
				//locking
				if !ok {
					req.ch <- true
					q[req.k] = &lockQItem{req.ch, nil, nil}
					continue
				}
				newItem := &lockQItem{req.ch, nil, nil}
				last := qTop.last
				if last == nil {
					last = qTop
				}
				last.next = newItem
				qTop.last = newItem

			} else {
				//unlocking
				if !ok {
					//ERROR tried to unlock something not locked
					continue
				}
				next := qTop.next
				if next != nil {
					next.ch <- true
					next.last = qTop.last
					q[req.k] = next
				} else {
					delete(q, req.k)
				}
			}

		}

	}()
	return ch
}

func (l Locker) Lock(s string) {
	ch := make(chan bool)
	l <- lockMessage{s, ch}
	_ = <-ch
}

func (l Locker) Unlock(s string) {
	l <- lockMessage{s, nil}
}

func newDB(root string) *DBase {

	return *DBase(&root)
}

/*TODO make Locking Version
func NewDB(root string, lock bool) *DBase {

	return &DBase{root, beginLocker()}
}*/

//Returns Bool OK true on success, false on fail
func (db DBase) WriteMap(key string, val []byte, hasLock bool) bool {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db, hexkey)
	f, err := os.Create(fname)
	if err != nil {
		return false
	}
	defer f.Close()

	f.Write(val)

	return true
}

func (db DBase) ReadMap(key string) []byte {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db.root, hexkey)

	f, err := os.Open(fname)
	defer f.Close()
	if err != nil {
		return nil
	}

	res := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		count, _ := f.Read(buf)
		res = append(res, buf[:count]...)
		if count < 1024 {
			return res
		}
		f.Seek(int64(count), os.SEEK_CUR)

	}

}
