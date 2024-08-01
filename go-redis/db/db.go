package db

import "fmt"

var KV_DB = DB{
	kv_pair: make(map[string]interface{}, 0),
}

type DB struct {
	kv_pair map[string]interface{}
}

func (db *DB) Set(key string, value interface{}) {
	db.kv_pair[key] = value
}

func (db *DB) Get(key string) (interface{}, error) {
	val, ok := db.kv_pair[key]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return val, nil
}
