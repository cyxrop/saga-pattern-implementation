package main

import (
	"log"
	"os"
	"time"
)

func main() {
	log.Println("hello world")
	log.Println("test env:", os.Getenv("TESTENV"))

	c := 0
	for {
		time.Sleep(time.Second)
		c++
		log.Println("counter: ", c)
	}
}
