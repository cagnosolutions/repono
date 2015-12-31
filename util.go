package repono

import (
	"bufio"
	"log"
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
