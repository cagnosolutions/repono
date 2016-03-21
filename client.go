package repono

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
)

type Client struct {
	conn     *net.TCPConn
	w        *bufio.Writer
	r        *bufio.Reader
	ds       *DataStore
	embedded bool
}

func (c Client) read() []byte {
	b, err := c.r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		log.Printf("Encountered an error while reading: %s\n", err)
		c.conn.Close()
		return nil
	}
	return dropCRLF(b)
}

func (c Client) getBool(b []byte) bool {
	if b != nil && b[0] == 1 {
		return true
	}
	return false
}

func Dial(host string) *Client {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatal("Remote Addr: ", err)
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Fatal("Dial Addr: ", err)
	}
	return &Client{
		conn,
		bufio.NewWriter(conn),
		bufio.NewReader(conn),
		nil,
		false,
	}
}

func Embed(ds *DataStore) *Client {
	return &Client{
		nil,
		nil,
		nil,
		ds,
		true,
	}
}

func (c Client) Ping() bool {
	write(c.w, PING)
	return c.getBool(c.read())
}

func (c Client) Quit() {
	write(c.w, QUIT)
	c.conn.Close()
}

func (c Client) UUID() string {
	return UUID1()
}

func (c Client) AddStore(store string) bool {
	var b []byte
	if c.embedded {
		c.ds.AddStore(store)
	} else {
		write(c.w, encode([][]byte{ADDSTORE, []byte(store)}))
		b = c.read()
	}
	return c.getBool(b)
}

func (c Client) GetAll(store string, ptr interface{}) {
	var b []byte
	if c.embedded {
		b = c.ds.GetAll(store)
	} else {
		write(c.w, encode([][]byte{GETALL, []byte(store)}))
		b = c.read()
	}
	if err := json.Unmarshal(b, ptr); err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
	}
}

func (c Client) HasStore(store string) bool {
	var b []byte
	if c.embedded {
		b = c.ds.HasStore(store)
	} else {
		write(c.w, encode([][]byte{DELSTORE, []byte(store)}))
		b = c.read()
	}
	return c.getBool(b)
}

func (c Client) DelStore(store string) bool {
	var b []byte
	if c.embedded {
		b = c.ds.DelStore(store)
	} else {
		write(c.w, encode([][]byte{HASSTORE, []byte(store)}))
		b = c.read()
	}
	return c.getBool(b)
}

func (c Client) Put(store string, value interface{}) bool {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	var b2 []byte
	if c.embedded {
		b2 = c.ds.Put(store, b)
	} else {
		write(c.w, encode([][]byte{PUT, []byte(store), b}))
		b2 = c.read()
	}
	return c.getBool(b2)
}

func (c Client) Add(store, key string, value interface{}) bool {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	var b2 []byte
	if c.embedded {
		b2 = c.ds.Add(store, []byte(key), b)
	} else {
		write(c.w, encode([][]byte{ADD, []byte(store), []byte(key), b}))
		b2 = c.read()
	}
	return c.getBool(b2)
}

func (c Client) Set(store, key string, value interface{}) bool {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	var b2 []byte
	if c.embedded {
		b2 = c.ds.Set(store, []byte(key), b)
	} else {
		write(c.w, encode([][]byte{SET, []byte(store), []byte(key), b}))
		b2 = c.read()
	}
	return c.getBool(b2)
}

func (c Client) Get(store, key string, ptr interface{}) {
	var b []byte
	if c.embedded {
		b = c.ds.Get(store, []byte(key))
	} else {
		write(c.w, encode([][]byte{GET, []byte(store), []byte(key)}))
		b = c.read()
	}
	err := json.Unmarshal(b, ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
	}
}

func (c Client) Del(store, key string) bool {
	var b []byte
	if c.embedded {
		b = c.ds.Del(store, []byte(key))
	} else {
		write(c.w, encode([][]byte{DEL, []byte(store), []byte(key)}))
		b = c.read()
	}
	return c.getBool(b)
}

func (c Client) Has(store, key string) bool {
	var b []byte
	if c.embedded {
		b = c.ds.Has(store, []byte(key))
	} else {
		write(c.w, encode([][]byte{HAS, []byte(store), []byte(key)}))
		b = c.read()
	}
	return c.getBool(b)
}

func (c Client) Query(store string, ptr interface{}, query ...[]byte) bool {
	var b []byte
	if c.embedded {
		b = c.ds.Query(store, query)
	} else {
		write(c.w, encode(append([][]byte{QUERY, []byte(store)}, query...)))
		b = c.read()
	}
	err := json.Unmarshal(b, ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
		return false
	}
	return true
}

func encode(bb [][]byte) []byte {
	return bytes.Join(bb, []byte{byte(DELIM)})
}
