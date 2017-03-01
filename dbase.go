package dbase

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
)

type DMapper interface {
	ReadMap(k string) ([]byte, error)
	WriteMap(k string, v []byte) bool
}

type DBase struct {
	root   string
	locker Locker
}

func NewDB(root string) *DBase {
	return &DBase{root, beginLocker()}
}

//Returns Bool OK true on success, false on fail
func (db DBase) WriteMap(key string, val []byte) bool {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db.root, hexkey)
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(val)

	return
}

func (db DBase) ReadMap(key string) ([]byte, error) {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db.root, hexkey)
	db.locker.Lock(fname)

	return ioutil.ReadFile(fname)

	/* Apparently I don't need to do all this work lol
		f, err := os.Open(fname)
		defer f.Close()
		if err != nil {
			return nil, err
		}

		res := make([]byte, 0)
		buf := make([]byte, 1024)
		for {
			count, _ := f.Read(buf)
			res = append(res, buf[:count]...)
			if count < 1024 {
				return res, nil
			}
			f.Seek(int64(count), os.SEEK_CUR)
	<<<<<<< HEAD
		}*/
}
