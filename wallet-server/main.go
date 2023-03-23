package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Block: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP port for wallet server")
	gateway := flag.String("gateway", "http://localhost:5000", "TCP port for server")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
	app.Run()
}
