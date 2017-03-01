package dbase

type DMapper interface {
	ReadMap(k string) []byte
	WriteMap(k string, v []byte)
}

type LockDMapper struct {
	db DMapper
	l  Locker
}

type lockMessage struct {
	k  string
	ch chan bool
}

type lockQItem struct {
	ch         chan bool
	next, last *lockQItem
}

type Locker chan lockMessage

func (l Locker) Lock(s string) {
	ch := make(chan bool)
	l <- lockMessage{s, ch}
	_ = <-ch
}

func (l Locker) Unlock(s string) {
	l <- lockMessage{s, nil}
}

func BeginLocker() Locker {
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

// NewLockDMapper returns a wrapper which will maintain locks tidily for any DMapper
func NewLockDMapper(m DMapper) *LockDMapper {
	return &LockDMapper{m, BeginLocker()}
}

func (self *LockDMapper) Read(k string, holdLock bool) []byte {
	self.l.Lock(k)
	if !holdLock {
		defer self.l.Unlock(k)
	}
	return self.db.ReadMap(k)

}

func (self *LockDMapper) Write(k string, v []byte, hasLock bool) {
	if !hasLock {
		self.l.Lock(k)
	}
	defer self.l.Unlock(k)
	self.db.WriteMap(k, v)
}

func (self *LockDMapper) Release(k string) {
	self.l.Unlock(k)
}
