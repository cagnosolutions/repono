package repono

import (
	"bufio"
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

// localize uuid call, no need to reach out to the server, right?
func (c Client) UUID() string {
	//write(c.w, UUID)
	//return string(c.read())
	return UUID1()
}

func (c Client) AddStore(s string) bool {
	write(c.w, encode([][]byte{ADDSTORE, []byte(s)}))
	return c.getBool()
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

func (c Client) DelStore(s string) bool {
	write(c.w, encode([][]byte{HASSTORE, []byte(s)}))
	return c.getBool()
}

func (c Client) Add(s, k string, v interface{}) bool {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	write(c.w, encode([][]byte{ADD, []byte(s), []byte(k), b}))
	return c.getBool()
}

func (c Client) Set(s, k string, v interface{}) bool {
	b, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling value: %s\n", err)
		return false
	}
	write(c.w, encode([][]byte{ADD, []byte(s), []byte(k), b}))
	return c.getBool()
}

func (c Client) Get(s, k string, ptr interface{}) {
	write(c.w, encode([][]byte{GET, []byte(s), []byte(k)}))
	err := json.Unmarshal(c.read(), ptr)
	if err != nil {
		log.Printf("Error unmarshaling value: %s\n", err)
	}
}

func (c Client) Del(s, k string) bool {
	write(c.w, encode([][]byte{GET, []byte(s), []byte(k)}))
	return c.getBool()
}

func (c Client) Has(s, k string) bool {
	write(c.w, encode([][]byte{HAS, []byte(s), []byte(k)}))
	return c.getBool()
}
