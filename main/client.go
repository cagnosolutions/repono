package main

import (
	"fmt"

	"github.com/cagnosolutions/repono"
)

type User struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Active bool   `json:"active"`
}

type Order struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
	Paid     bool   `json:"paid"`
}

func main() {
	c := repono.Dial("localhost:9999")
	fmt.Println("adding stores...")
	c.AddStore("user")
	c.AddStore("order")
	fmt.Println("adding users...")

	uuid1 := c.UUID()
	c.Add("user", uuid1, User{
		Name:   uuid1,
		Age:    99,
		Active: false,
	})

	for i := 0; i < 10; i++ {
		s := c.UUID()
		b := c.Add("user", s, User{
			Name:   s,
			Age:    i,
			Active: (i%2 == 0),
		})
		fmt.Printf("user %d add: %v\n", i, b)
	}
	fmt.Println("adding orders...")
	for i := 0; i < 10; i++ {
		s := c.UUID()
		b := c.Add("order", s, Order{
			Item:     s,
			Quantity: i,
			Paid:     (i%2 == 0),
		})
		fmt.Printf("user %d add: %v\n", i, b)
	}
	fmt.Println("adding finished!")

	fmt.Println("getting a single user...")

	var user User
	c.Get("user", uuid1, &user)
	fmt.Printf("Get() -> %+v\n", user)

	fmt.Println("getting all users...")

	var users []User
	c.GetAll("user", &users)
	for _, u := range users {
		fmt.Printf("id: %s, user: %+v\n", u.Name, u)
	}
	fmt.Println()

	var q []User
	c.Query("user", &q, repono.C("age", repono.LT, "6"))
	for _, u := range q {
		fmt.Printf("id: %s, user: %+v\n", u.Name, u)
	}
}
