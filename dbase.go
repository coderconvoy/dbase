package dbase

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
)

<<<<<<< HEAD
type DMapper interface {
	ReadMap(k string) ([]byte, error)
	WriteMap(k string, v []byte) bool
}

type DBase struct {
	root   string
	locker Locker
}

func NewDB(root string) *DBase {
=======
type DBase string

func NewDBase(s string) *DBase {
	res := DBase(s)
	return &res
}

/*TODO make Locking Version
func NewDB(root string, lock bool) *DBase {

>>>>>>> cf264fc0524646f19296a4c9e49a16fe028afa23
	return &DBase{root, beginLocker()}
}*/

<<<<<<< HEAD
//Returns Bool OK true on success, false on fail
func (db DBase) WriteMap(key string, val []byte) bool {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db.root, hexkey)
=======
//WriteMap Returns Bool OK true on success, false on fail
func (db DBase) WriteMap(key string, val []byte) {
	s := string(db)
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(s, hexkey)
>>>>>>> cf264fc0524646f19296a4c9e49a16fe028afa23
	f, err := os.Create(fname)
	if err != nil {
		return
	}
	defer f.Close()

	f.Write(val)

	return
}

<<<<<<< HEAD
func (db DBase) ReadMap(key string) ([]byte, error) {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db.root, hexkey)
	db.locker.Lock(fname)

	return ioutil.ReadFile(fname)

	/* Apparently I don't need to do all this work lol
=======
//ReadMap Takes a key and tries to find the result returns
//nil on empty
func (db DBase) ReadMap(key string) []byte {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(string(db), hexkey)

>>>>>>> cf264fc0524646f19296a4c9e49a16fe028afa23
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
=======
	}
>>>>>>> cf264fc0524646f19296a4c9e49a16fe028afa23
}
