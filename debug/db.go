package debug

import (
	"fmt"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type entry struct {
	table  string
	object interface{}
}

type DB struct {
	dbClient *r.Session
	dbName   string
	entries  []entry
	cap      int
}

func NewDB(dbClient *r.Session, dbName string, cap int) *DB {
	return &DB{
		dbClient: dbClient,
		dbName:   dbName,
		entries:  make([]entry, 0, cap),
		cap:      cap,
	}
}

func (db *DB) Flush() {
	fmt.Println("Flush")

	m := map[string][]interface{}{}
	for _, v := range db.entries {
		m[v.table] = append(m[v.table], v.object)
	}

	for t, entries := range m {
		fmt.Println("Insert")
		if err := r.DB(db.dbName).Table(t).Insert(entries).Exec(db.dbClient); err != nil {
			panic(err)
		}
	}

	db.entries = db.entries[:0]
}

func (db *DB) Insert(table string, object interface{}) {
	db.entries = append(db.entries, entry{table, object})
	if len(db.entries) >= db.cap {
		db.Flush()
	}
}
