package repono

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"time"
)

/*
type Server struct {
	ds *DataStore
}

func NewServer() *Server {
	return &Server{
		ds: NewDataStore(),
	}
}*/

func ListenAndServe(port string) {
	ds := NewDataStore()
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
		go handleConn(ds, conn)
	}
}

func handleConn(ds *DataStore, conn *net.TCPConn) {
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
		bb := bytes.SplitN(dropCRLF(b), []byte{byte(DELIM)}, 4)
		cmd := bb[0]
		switch len(bb) {
		case 1: // ping, quit
			switch {
			case bytes.Equal(cmd, PING):
				write(w, TRUE)
			case bytes.Equal(cmd, QUIT):
				log.Println("Client quit")
				conn.Close()
				return
			default:
				write(w, ERR1)
			}
		case 2:
			store := string(bb[1])
			switch {
			case bytes.Equal(cmd, ADDSTORE):
				ds.AddStore(store)
				write(w, TRUE)
			case bytes.Equal(cmd, GETALL):
				write(w, ds.GetAll(store))
			case bytes.Equal(cmd, DELSTORE):
				ds.DelStore(store)
				write(w, TRUE)
			case bytes.Equal(cmd, HASSTORE):
				write(w, ds.HasStore(store))
			default:
				write(w, ERR2)
			}
		case 3:
			store := string(bb[1])
			key := bb[2]
			switch {
			case bytes.Equal(cmd, GET):
				write(w, ds.Get(store, key))
			case bytes.Equal(cmd, DEL):
				write(w, ds.Del(store, key))
			case bytes.Equal(cmd, HAS):
				write(w, ds.Has(store, key))
			default:
				write(w, ERR3)
			}
		case 4:
			store := string(bb[1])
			key := bb[2]
			val := bb[3]
			switch {
			case bytes.Equal(cmd, ADD):
				write(w, ds.Add(store, key, val))
			case bytes.Equal(cmd, SET):
				write(w, ds.Set(store, key, val))
			default:
				write(w, ERR4)
			}
		default:
			write(w, ERR5)
		}
		w.Flush()
	}
	log.Println("Disconnecting", conn.RemoteAddr().String(), "Goodbye!")
	conn.Close()
	return
}
