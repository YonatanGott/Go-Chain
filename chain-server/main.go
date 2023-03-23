package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Block: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP port for server")
	flag.Parse()

	app := NewServer(uint16(*port))
	app.Run()
}
