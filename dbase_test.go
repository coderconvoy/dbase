package dbase

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	locker := beginLocker()

	locker.Lock("Hello")

	locker.Unlock("Hello")
	locker.Lock("Hello")

	locker.Unlock("Hello")
}

func TestLong(t *testing.T) {
	locker := beginLocker()
	ch := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			locker.Lock("Hello")
			fmt.Printf("Locked %d\n", n)
			time.Sleep(time.Second / 10)
			fmt.Printf("AllDone %d\n", n)
			locker.Unlock("Hello")
			ch <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		_ = <-ch
	}

}

func TestSaveLoad(t *testing.T) {
	k := "KEY"
	b := []byte("VALUE")
	db := NewDB("dbase_testdata/t1")
	db.WriteMap(k, b, false)
	b2 := db.ReadMap(k, false)
	if bytes.Compare(b, b2) != 0 {
		t.Fail()
	}
}
