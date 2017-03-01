package dbase

import (
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	locker := NewLocker()
	data := [][]string{
		[]string{"a"},
		[]string{"b", "c"},
		[]string{"a", "c"},
		[]string{"a", "b"},
	}
	mp := make(map[string]int)
	mp["a"] = 0
	mp["b"] = 0
	mp["c"] = 0

	ch := make(chan int)

	f := func(ss []string) {
		id := locker.Lock(ss...)
		n := 0
		for _, s := range ss {
			mp[s] = mp[s] + 1
			n += mp[s]
		}
		time.Sleep(time.Second / 500)

		n2 := 0
		for _, s := range ss {
			n2 += mp[s]
			mp[s] = mp[s] + 1

		}
		if n2 != n {
			t.Log("n2 != n")
			t.Fail()
		}
		locker.Unlock(id)
		ch <- len(ss) * 2
	}

	for i := 0; i < 100; i++ {
		go f(data[i%len(data)])
	}

	tot := 0
	for i := 0; i < 100; i++ {
		tot += <-ch
	}

	dtot := 0
	for _, v := range mp {
		dtot += v
	}
	if tot != dtot {
		t.Logf("dtot != tot , %d != %d", dtot, tot)
		t.Fail()
	}

}
