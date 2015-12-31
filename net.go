package repono

import (
	"bufio"
	"bytes"
	"log"
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
	HAS = []byte{'h', 'a', 's'}
)

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

func write(w *bufio.Writer, b []byte) {
	n, err := w.Write(b)
	if n < 1 || err != nil {
		log.Printf("Wrote %d bytes, error: %s\n", n, err)
		return
	}
	n, err = w.Write(CRLF)
	if n < 1 || err != nil {
		log.Printf("Wrote %d bytes, error: %s\n", n, err)
		return
	}
	err = w.Flush()
	if err != nil {
		log.Printf("Error flushing write buffer: %s\n", err)
		return
	}
}

func encode(bb [][]byte) []byte {
	return bytes.Join(bb, []byte{byte(DELIM)})
}
