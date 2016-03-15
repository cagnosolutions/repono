package repono

import (
	"bytes"
	"fmt"
)

func queryAuth(email, password string) QueryFunc {
	return func(doc string) bool {
		return bytes.Contains(doc, []byte(`"email":"`+email+`"`)) && bytes.Contains(doc, []byte(`"password":"`+password+`""`)) && bytes.Contains(doc, []byte(`"active":true"`))
	}
}

type DB struct {
	ss []string
}

func NewDB() *DB {
	return &DB{[]string{"greg", "scott", "rosalie", "kayla"}}
}

func testQ(doc []byte, q ...string) {
	if len(q)%2 != 0 {
		return
	}
	var i int
	for i = 0; i <= len(q); i += 2 {
		comp := []byte(fmt.Sprintf("%q:%s", q[i], q[i+1]))
		if bytes.Contains(doc, comp) {
			break
		}
	}
	if len(q)+2 == i {
		// append to res
	}

}

func (db *DB) Query(q QueryFunc) {
	for _, v := range db.ss {
		if q(v) {
			// append to results
		}
	}

}

type Q interface {
	Query(fn func() bool) []string
}

type QueryFunc func(doc string) bool
