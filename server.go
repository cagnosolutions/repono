package repono

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ListenAndServe(port string, path ...string) {
	if len(path) == 1 {
		DB_PATH = path[0]
	}
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
	catchSigInts()
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

// catch interrupt signals
func catchSigInts() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGSTOP)
	go func() {
		fmt.Printf("\nCaught signal: %v, exiting...\n", <-sig)
		os.Exit(0)
	}()
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
		if bytes.Equal(bb[0], QUERY) && len(bb) > 2 {
			write(w, ds.Query(string(bb[1]), bb[2:]))
			continue
		}
		cmd := bb[0]
		switch len(bb) {
		case 1:
			switch {
			case bytes.Equal(cmd, UUID):
				write(w, ds.UUID())
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
				write(w, ds.AddStore(store))
			case bytes.Equal(cmd, GETALL):
				write(w, ds.GetAll(store))
			case bytes.Equal(cmd, DELSTORE):
				write(w, ds.DelStore(store))
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
				for _, b := range bb {
					fmt.Printf("%s\n", b)
				}
				write(w, ERR4)
			}
		default:
			write(w, ERR5)
		}
	}
	log.Println("Disconnecting", conn.RemoteAddr().String(), "Goodbye!")
	conn.Close()
	return
}
