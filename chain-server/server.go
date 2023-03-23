package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"go-block/chain"
	"go-block/wallet"
)

var cache map[string]*chain.Blockchain = make(map[string]*chain.Blockchain)

type Server struct {
	port uint16
}

func NewServer(port uint16) *Server {
	return &Server{port}
}

func (server *Server) Port() uint16 {
	return server.port
}

func (server *Server) GetBlockchain() *chain.Blockchain {
	blockchain, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		blockchain = chain.NewBlockchain(minerWallet.BlockchainAddress(), server.port)
		cache["blockchain"] = blockchain
		log.Printf("private key %v", minerWallet.PrivateKeyStr())
		log.Printf("public key %v", minerWallet.PublicKeyStr())
		log.Printf("chain address key %v", minerWallet.BlockchainAddress())
	}
	return blockchain
}

func (server *Server) GetChain(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchain := server.GetBlockchain()
		marshal, _ := blockchain.MarshalJSON()
		io.WriteString(writer, string(marshal[:]))
	default:
		log.Printf("Error: Invalid Http Method")
	}
}

func (server *Server) Run() {
	http.HandleFunc("/", server.GetChain)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(server.Port())), nil))
}
