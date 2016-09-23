package dbase

import (
	"encoding/hex"
	"os"
	"path"
)

type DBase string

func NewDBase(s string) *DBase {
	res := DBase(s)
	return &res
}

/*TODO make Locking Version
func NewDB(root string, lock bool) *DBase {

	return &DBase{root, beginLocker()}
}*/

//WriteMap Returns Bool OK true on success, false on fail
func (db DBase) WriteMap(key string, val []byte) {
	s := string(db)
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(s, hexkey)
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(val)

	return
}

//ReadMap Takes a key and tries to find the result returns
//nil on empty
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
