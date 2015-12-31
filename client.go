package repono

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
)

var DELIM = '|'

var CRLF = []byte{'\r', '\n'}

var (
	PING = []byte{'p', 'i', 'n', 'g'}
	QUIT = []byte{'q', 'u', 'i', 't'}
)

var (
	ADDSTORE = []byte{'a', 'd', 'd', 's', 't', 'o', 'r', 'e'}
	GETALL   = []byte{'g', 'e', 't', 'a', 'l', 'l'}
	DELSTORE = []byte{'d', 'e', 'l', 's', 't', 'o', 'r', 'e'}
	HASSTORE = []byte{'h', 'a', 's', 's', 't', 'o', 'r', 'e'}
)

var (
	ADD = []byte{'a', 'd', 'd'}
	SET = []byte{'s', 'e', 't'}
	GET = []byte{'g', 'e', 't'}
	DEL = []byte{'d', 'e', 't'}
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

func encode(bb [][]byte) []byte {
	return bytes.Join(bb, []byte{byte(DELIM)})
}

func (c Client) getBool() bool {
	b := c.read()
	if b != nil && b[0] == 1 {
		return true
	}
	return false
}

func Dial(host string) *Client {
	laddr, err := net.ResolveTCPAddr("tcp", "localhost")
	if err != nil {
		log.Fatal(err)
	}
	raddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialTCP("tcp", laddr, raddr)
	if err != nil {
		log.Fatal(err)
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

func (c Client) AddStore(s string) {
	write(c.w, encode([][]byte{ADDSTORE, []byte(s)}))
}

func (c Client) GetAll(s string, ptr interface{}) {
	write(c.w, encode([][]byte{GETALL, []byte(s)}))
	err := json.Unmarshal(c.read(), ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
	}
}

func (c Client) HasStore(s string) bool {
	write(c.w, encode([][]byte{DELSTORE, []byte(s)}))
	return c.getBool()
}

func (c Client) DelStore(s string) {
	write(c.w, encode([][]byte{HASSTORE, s}))
}

func (c Client) Add(s, k string, v interface{}) bool {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return
	}
	c.write(encode([][]byte{ADD, s, k, b}))
	return c.getBool()
}

func (c Client) Set(s, k string, v interface{}) bool {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return
	}
	c.write(encode([][]byte{ADD, s, k, b}))
	return c.getBool()
}

func (c Client) Get(s, k string, ptr interface{}) {
	c.write(encode([][]byte{GET, s, k}))
	err := json.Unmarshal(c.read(), ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
	}
}

func (c Client) Del(s, k string) {
	c.write(encode([][]byte{GET, s, k}))
}

func (c Client) Has() {
}
