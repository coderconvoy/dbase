package dbase

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	locker := BeginLocker()

	locker.Lock("Hello")

	locker.Unlock("Hello")
	locker.Lock("Hello")

	locker.Unlock("Hello")
}

func TestLong(t *testing.T) {
	locker := BeginLocker()
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
	db := NewDBase("dbase_testdata/t1")
	db.WriteMap(k, b)
	b2 := db.ReadMap(k)
	if bytes.Compare(b, b2) != 0 {
		t.Fail()
	}
}

func TestLockedDBIncrements(t *testing.T) {
	db := NewLockDMapper(NewDBase("dbase_testdata/t2"))

	ch := make(chan bool)

	for i := 0; i < 100; i++ {
		go func(n int) {
			time.Sleep(100)
			a := db.Read("Poop", true)
			//time.Sleep(10)
			b := 0
			if a != nil {
				as := string(a)
				var err error
				b, err = strconv.Atoi(as)

				if err != nil {
					println(err)
				}
			}
			c := b + n
			fmt.Printf("%d + %d = %d\n", b, n, c)
			db.Write("Poop", []byte(strconv.Itoa(c)), true)
			ch <- true

		}(i)

	}

	for i := 0; i < 100; i++ {
		_ = <-ch
	}
}
