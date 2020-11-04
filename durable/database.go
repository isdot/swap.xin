package durable

import (
	"fmt"
	"log"

	"github.com/timshannon/badgerhold"
)

type Database struct {
	db *badgerhold.Store
}

func OpenDatabaseClient() *Database {
	conn := fmt.Sprintf("/root/swapxin/db")
	options := badgerhold.DefaultOptions
	options.Dir = conn
	options.ValueDir = conn

	db, err := badgerhold.Open(options)

	if err != nil {
		log.Fatal(err)
	}

	return &Database{db: db}
}

func (d *Database) UpdateMatching(data interface{}, query *badgerhold.Query, update func(interface{}) error) {
	d.db.UpdateMatching(data, query, update)
}

func (d *Database) Update(key interface{}, data interface{}) error {
	err := d.db.Update(key, data)
	return err
}

func (d *Database) Insert(data interface{}) error {
	key := badgerhold.NextSequence()
	err := d.db.Insert(key, data)
	return err
}

func (d *Database) Find(data interface{}, query *badgerhold.Query) error {
	err := d.db.Find(data, query)
	return err
}

func (d *Database) Delete(data interface{}, query *badgerhold.Query) error {
	err := d.db.DeleteMatching(data, query)
	return err
}
