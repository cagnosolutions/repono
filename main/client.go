package main

import (
	"fmt"

	"github.com/cagnosolutions/repono"
)

type User struct {
	Name   string
	Age    int
	Active bool
}

type Order struct {
	Item     string
	Quantity int
	Paid     bool
}

func main() {
	c := repono.Dial("localhost:9999")
	fmt.Println("adding stores...")
	c.AddStore("user")
	c.AddStore("order")
	fmt.Println("adding users...")
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%d", i)
		b := c.Add("user", s, User{
			Name:   s,
			Age:    i,
			Active: (i%2 == 0),
		})
		fmt.Printf("user %d add: %v\n", i, b)
	}
	fmt.Println("adding orders...")
	for i := 0; i < 10; i++ {
		s := fmt.Sprintf("%d", i)
		b := c.Add("order", s, Order{
			Item:     s,
			Quantity: i,
			Paid:     (i%2 == 0),
		})
		fmt.Printf("user %d add: %v\n", i, b)
	}
	fmt.Println("adding finished!")
}
