package repono

import (
	"bufio"
	"bytes"
	"log"
)

var DELIM = '|'
var MAXBUF = 4096

var CRLF = []byte{'\r', '\n'}

var (
	UUID = []byte{'u', 'u', 'i', 'd'}
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
	PUT = []byte{'p', 'u', 't'}
	ADD = []byte{'a', 'd', 'd'}
	SET = []byte{'s', 'e', 't'}
	GET = []byte{'g', 'e', 't'}
	DEL = []byte{'d', 'e', 't'}
	HAS = []byte{'h', 'a', 's'}
	QUERY = []byte{'q', 'u', 'e', 'r','y'}
)

var (
	TRUE  = []byte{1}
	FALSE = []byte{0}
	ERR1  = []byte{'e', 'r', 'r', '1'}
	ERR2  = []byte{'e', 'r', 'r', '2'}
	ERR3  = []byte{'e', 'r', 'r', '3'}
	ERR4  = []byte{'e', 'r', 'r', '4'}
	ERR5  = []byte{'e', 'r', 'r', '5'}
	NIL   = []byte{}
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
	if n < 1 {
		log.Printf("Error, %d bytes to write\n", n)
		return
	} else if err != nil {
		log.Printf("Error writing: %s\n", err)
		return
	}
	n, err = w.Write(CRLF)
	if n < 1 {
		log.Printf("Error, %d bytes to write\n", n)
		return
	} else if err != nil {
		log.Printf("Error writing: %s\n", err)
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
