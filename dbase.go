package dbase2

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
	Root string
}

//Returns Bool OK true on success, false on fail
func (db DBase) WriteMap(key string, val []byte) bool {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db.Root, hexkey)
	f, err := os.Create(fname)
	if err != nil {
		return false
	}
	defer f.Close()
	f.Write(val)
	return true
}

func (db DBase) ReadMap(key string) ([]byte, error) {
	hexkey := hex.EncodeToString([]byte(key))
	fname := path.Join(db.Root, hexkey)
	return ioutil.ReadFile(fname)
}
