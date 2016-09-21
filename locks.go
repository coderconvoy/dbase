package dbase

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
