package main

import "fmt"

func main() {
	chann := make(chan int)
	chann <- 1
	fmt.Println("Hello, playground")
}
