package main

import (
	"fmt"
	"github.com/heck-go/bytag"
)

type Address struct{
	Number int `mytag:"n"`
	
	Street string `mytag:"s"`
	
	ZipCode int `mytag:"z"`
}

type Person struct {
	Name string `mytag:"n"`
	
	Age int
	
	Address []Address `mytag:"ad"`
}

func main() {
	per := Person{}

	bytag.Bind("mytag", &per, map[string]interface{}{
		"n": "Teja",
		"Age": 29,
		"ad": map[string]interface{} {
			"n": 25,
			"s": "Street",
			"z": 16453,
		},
	})
	
	fmt.Println(per)
}
