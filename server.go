package repono

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	ListenAndServe(":9999")
}

func ListenAndServe(port string) {
	log.Println("Server starting...")
	addr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on", port, "for connections...")
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Handling connection from %v\n", conn.RemoteAddr().String())
		go handleConn(conn)
	}
}

func handleConn(conn *net.TCPConn) {
	r, w := bufio.NewReader(conn), bufio.NewWriter(conn)
	for {
		b, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Got error, disconnecting client:", err)
			conn.Close()
			return
		} else {
			conn.SetDeadline(time.Now().Add(8 * time.Minute))
		}
		if len(b) > MAXBUF {
			log.Println("Client sent data exceeding size limit of 4096")
			w.Write([]byte("Exceeded data size (limit 4096)\n"))
			w.Flush()
			continue
		}
		bb := bytes.SplitN(dropCRLF(b), []byte{'|'}, 4)
		switch len(bb) {
		case 1: // ping, quit
			switch string(bb[0]) {
			case "ping":
				w.Write([]byte("pong\n"))
			case "quit":
				log.Println("Client quit")
				w.Write([]byte("Goodbye!\n"))
				w.Flush()
				conn.Close()
				return
			default:
				log.Printf("(hit case 1) bb[0]: %s\n", bb[0])
				w.Write([]byte("hit case 1\n"))
			}
		case 2:
			w.Write([]byte("hit case 2\n"))
		case 3:
			w.Write([]byte("hit case 3\n"))
		case 4:
			w.Write([]byte("hit case 4\n"))
		default:
			w.Write([]byte("hit default\n"))
		}
		w.Flush()
	}
	log.Println("Disconnecting", conn.RemoteAddr().String(), "Goodbye!")
	conn.Close()
	return
}
