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
	DELSTORE = []byte{'d', 'e', 'l', 's', 't', 'o', 'r', 'e'}
	GETALL   = []byte{'g', 'e', 't', 'a', 'l', 'l'}
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

func (c Client) write(b []byte) {
	n, err := c.w.Write(b)
	if n < 1 || err != nil {
		log.Printf("Wrote %d bytes, error: %s\n", n, err)
		return
	}
	n, err = c.w.Write(CRLF)
	if n < 1 || err != nil {
		log.Printf("Wrote %d bytes, error: %s\n", n, err)
		return
	}
	err = c.w.Flush()
	if err != nil {
		log.Printf("Error flushing write buffer: %s\n", err)
		return
	}
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

func dropCRLF(line []byte) []byte {
	if line[len(line)-1] == '\n' {
		drop := 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		line = line[:len(line)-drop]
	}
	return line
}

func encode(bb [][]byte) []byte {
	return bytes.Join(bb, []byte{DELIM})
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
	c.write(PING)
	b := c.read()
	if b != nil && b[0] == 1 {
		return true
	}
	return false
}

func (c Client) Quit() {
	c.write(QUIT)
	c.conn.Close()
}

func (c Client) AddStore(s string) {

}

func (c Client) DelStore(s string) {

}

func (c Client) GetAll(s string) {

}

func (c Client) Add(s, k string, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return
	}
	c.write(encode([][]byte{ADD, s, k, v, b}))
}

func (c Client) Set() {
}

func (c Client) Get() {
}

func (c Client) Del() {
}
