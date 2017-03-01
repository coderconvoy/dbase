package dbase

type lockMessage struct {
	k  []string
	id uint64
	ch chan uint64
}

type Locker chan lockMessage

func (l Locker) Lock(s ...string) uint64 {
	ch := make(chan uint64)
	l <- lockMessage{s, 0, ch}
	return <-ch
}

func (l Locker) Unlock(id uint64) {
	l <- lockMessage{nil, id, nil}
}

func grabLock(id uint64, ks []string, mp map[string]uint64) (uint64, bool) {
	for _, k := range ks {
		_, cl := mp[k]
		if cl {
			return id, false
		}
	}
	id++
	for _, k := range ks {
		mp[k] = id
	}
	return id, true
}

func NewLocker() Locker {
	ch := make(Locker)

	go func() {
		q := []lockMessage{}
		locks := make(map[string]uint64)
		cid := uint64(0)
		for {
			req := <-ch

			if req.ch != nil {
				//locking
				var added bool
				cid, added = grabLock(cid, req.k, locks)
				if added {
					req.ch <- cid
					continue
				}
				q = append(q, req)
				continue
			}

			//unlocking
			for i, k := range locks {
				if k == req.id {
					delete(locks, i)
				}
			}
			//see what we can take out of the queue
			newq := []lockMessage{}
			for _, l := range q {
				var added bool
				cid, added = grabLock(cid, l.k, locks)
				if added {
					l.ch <- cid
				} else {
					newq = append(newq, l)
				}
			}
			q = newq

		}

	}()
	return ch
}
