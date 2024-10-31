package utils

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	dbPath = "/var/lib/project-erp"
	db     *leveldb.DB
	dbLock sync.Mutex
)

func init() {
	var err error
	db, err = leveldb.OpenFile(dbPath, nil)
	if err != nil {
		panic(err)
	}
}
func GetValue(key string) (string, error) {
	dbLock.Lock()
	defer dbLock.Unlock()
	result, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
func SetValue(key string, value string) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	return db.Put([]byte(key), []byte(value), nil)
}
func DelValue(key string) error {
	dbLock.Lock()
	defer dbLock.Unlock()
	return db.Delete([]byte(key), nil)
}
func CloseDB() {
	dbLock.Lock()
	defer dbLock.Unlock()
	if db != nil {
		db.Close()
	}
}
