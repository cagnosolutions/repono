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
	conn *net.TCPConn
	w    *bufio.Writer
	r    *bufio.Reader
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

func (c Client) getBool() bool {
	b := c.read()
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
	}
}

func (c Client) Ping() bool {
	write(c.w, PING)
	return c.getBool()
}

func (c Client) Quit() {
	write(c.w, QUIT)
	c.conn.Close()
}

func (c Client) UUID() string {
	return UUID1()
}

func (c Client) AddStore(store string) bool {
	write(c.w, encode([][]byte{ADDSTORE, []byte(store)}))
	return c.getBool()
}

func (c Client) GetAll(store string, ptr interface{}) {
	write(c.w, encode([][]byte{GETALL, []byte(store)}))
	err := json.Unmarshal(c.read(), ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
	}
}

func (c Client) HasStore(store string) bool {
	write(c.w, encode([][]byte{DELSTORE, []byte(store)}))
	return c.getBool()
}

func (c Client) DelStore(store string) bool {
	write(c.w, encode([][]byte{HASSTORE, []byte(store)}))
	return c.getBool()
}

func (c Client) Put(store string, value interface{}) bool {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	write(c.w, encode([][]byte{PUT, []byte(store), b}))
	return c.getBool()
}

func (c Client) Add(store, key string, value interface{}) bool {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	write(c.w, encode([][]byte{ADD, []byte(store), []byte(key), b}))
	return c.getBool()
}

func (c Client) Set(store, key string, value interface{}) bool {
	b, err := json.Marshal(value)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	write(c.w, encode([][]byte{SET, []byte(store), []byte(key), b}))
	return c.getBool()
}

func (c Client) Get(store, key string, ptr interface{}) {
	write(c.w, encode([][]byte{GET, []byte(store), []byte(key)}))
	err := json.Unmarshal(c.read(), ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
	}
}

func (c Client) Del(store, key string) bool {
	write(c.w, encode([][]byte{DEL, []byte(store), []byte(key)}))
	return c.getBool()
}

func (c Client) Has(store, key string) bool {
	write(c.w, encode([][]byte{HAS, []byte(store), []byte(key)}))
	return c.getBool()
}

func (c Client) Query(store string, ptr interface{}, query ...[]byte) bool {
	write(c.w, encode(append([][]byte{QUERY, []byte(store)}, query...)))
	err := json.Unmarshal(c.read(), ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
		return false
	}
	return true
}

func encode(bb [][]byte) []byte {
	return bytes.Join(bb, []byte{byte(DELIM)})
}
