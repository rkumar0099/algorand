package main

import (
	"bytes"
	"fmt"
)

func main() {
	d := bytes.Join([][]byte{
		[]byte("Rabindar"),
		[]byte("Kumar"),
	}, nil)

	fmt.Println("Hello")
}
