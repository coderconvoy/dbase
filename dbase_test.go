package dbase

import (
	"encoding/hex"
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

func TestEncode(t *testing.T) {
	s := []byte("Hello")
	h := hex.EncodeToString(s)
	n, _ := hex.DecodeString(h)
	fmt.Println(string(n))
}
