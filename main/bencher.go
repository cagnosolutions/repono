package main

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"time"
)

var data = [][][]byte{
	[][]byte{[]byte("r1"), []byte("r9"), []byte("r4"), []byte("r2"), []byte("r14")},
	[][]byte{[]byte("r7"), []byte("r5"), []byte("r8"), []byte("r11"), []byte("r15")},
	[][]byte{[]byte("r6"), []byte("r13"), []byte("r12"), []byte("r3"), []byte("r10")},
}

func showData() {
	for i, b1 := range data {
		for j, b2 := range b1 {
			fmt.Printf("[%d][%d][]byte{%s}\n", i, j, b2)
		}
	}
}

func find3Linear(b []byte) (int64, []byte) {
	t1 := time.Now().UnixNano()
	var res []byte
	for _, x := range data {
		for _, y := range x {
			if bytes.Equal(y, b) {
				res = y
				break
			}
		}
	}
	return time.Now().UnixNano() - t1, res
}

func find3Divide(b []byte) (int64, []byte) {
	t1 := time.Now().UnixNano()
	var wg sync.WaitGroup
	wg.Add(len(data))
	var res []byte
	for i := 0; i < len(data); i++ {
		go divide(data[i], b, &wg, &res)
	}
	wg.Wait()
	return time.Now().UnixNano() - t1, res
}

func divide(b [][]byte, b1 []byte, wg *sync.WaitGroup, res *[]byte) {
	for _, y := range b {
		if bytes.Equal(y, b1) {
			*res = y
			break
		}
	}
	wg.Done()
}

func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	fmt.Printf("Running with %d CPU's\n", cpus)
}

func main() {
	showData()
	n1, b1 := find3Linear([]byte("r5"))
	fmt.Printf("[linear]\ttime:\t%d,\tres:\t%s\n", n1, b1)
	n2, b2 := find3Divide([]byte("r5"))
	fmt.Printf("[divide]\ttime:\t%d,\tres:\t%s\n", n2, b2)
}
