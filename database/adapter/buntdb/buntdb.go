package buntdb

import (
	"github.com/euforic/backend-base/database"
	"github.com/tidwall/buntdb"
)

const todoCollection = "todos"

var _ database.Adapter = &Database{} // or &myType{} or [&]myType if scalar

type Database struct {
	db      *buntdb.DB
	connStr string
}

func New(connStr string) *Database {
	d := Database{
		connStr: connStr,
	}
	return &d
}

// Connect opens a connection to to buntdb
func (d *Database) Connect() error {
	db, err := buntdb.Open(d.connStr)
	if err != nil {
		return err
	}

	if err := db.CreateIndex("todos", todoCollection+":*", buntdb.IndexString); err != nil {
		return err
	}

	d.db = db
	return nil
}

// Close closes out BuntDB connection
func (d Database) Close() error {
	return d.db.Close()
}
