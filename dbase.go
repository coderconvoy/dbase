package dbase

import (
	"encoding/hex"
	"os"
	"path"
)

type DMapper interface {
	ReadMap(k string) []byte
	WriteMap(k string, v []byte)
}

type DBase string

func NewDBase(s string) *DBase {
	res := DBase(s)
	return &res
}

/*TODO make Locking Version
func NewDB(root string, lock bool) *DBase {

	return &DBase{root, beginLocker()}
}*/

//Returns Bool OK true on success, false on fail
func (db DBase) WriteMap(key string, val []byte) bool {
	s := string(db)
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(s, hexkey)
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
	fname := path.Join(string(db), hexkey)

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
