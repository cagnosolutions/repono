package main

import (
	"log"
	"sync"

	"github.com/cagnosolutions/repono"
)

var (
	n = 64
	g sync.WaitGroup
	s = "test"
)

func addN(i, n int, g *sync.WaitGroup) {
	c := repono.Dial("localhost:9999")
	log.Printf("Client[%d] connected...\nAdding %d records...\n", i, n)
	for j := 0; j < n; j++ {
		if !c.Add(s, c.UUID(), []byte{'{', '"', 'i', 'd', '"', ':', byte(j), '}'}) {
			log.Printf("Client[%d] failed to add record %d\n", i, j)
		}
	}
	c.Quit()
	g.Done()
}

func main() {

	c := repono.Dial("localhost:9999")
	c.AddStore(s)
	c.Quit()

	g.Add(n)

	log.Printf("Spinning up %d client goroutines...\n", n)
	for i := 0; i < n; i++ {
		go addN(i, n, &g)
	}
	g.Wait()
	log.Println("Finished")
}
