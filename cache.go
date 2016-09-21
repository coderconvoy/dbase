package dbase

type CacheDB struct {
	db      DBase
	cache   map[string][]byte
	maxSize int
}

func NewCacheDB(root string, ms int) *CacheDB {
	return &CacheDB{
		NewDBase(root),
		make(map[string][]byte),
		ms,
	}
}

//TODO limit size of Cache in some way
func (self *CacheDB) ReadMap(k string) []byte {
	res, ok := self.cache[k]
	if ok {
		return res
	}
	v := self.db.ReadMap(k)
	self.cache[k] = v
	return v
}

func (self *CacheDB) WriteMap(k string, v []byte) {
	self.cache[k] = v
	self.db.WriteMap(k, v)
}
