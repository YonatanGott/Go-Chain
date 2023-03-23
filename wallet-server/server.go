package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"go-block/chain"
	"go-block/utils"
	"go-block/wallet"
)

const templateDir = "wallet-server/templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (server *WalletServer) Port() uint16 {
	return server.port
}

func (server *WalletServer) Gateway() string {
	return server.gateway
}

func (server *WalletServer) Index(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		template, _ := template.ParseFiles(path.Join(templateDir, "index.html"))
		template.Execute(writer, "")

	default:
		log.Printf("Error: Invalid Http Method")
	}
}

func (server *WalletServer) Wallet(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		writer.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		marshal, _ := myWallet.MarshalJSON()
		io.WriteString(writer, string(marshal[:]))

	default:
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: Invalid Http Method")
	}
}

func (server *WalletServer) CreateTransaction(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var transactionReq wallet.TransactionRequest
		err := decoder.Decode(&transactionReq)
		if err != nil {
			log.Printf("%v", err)
			io.WriteString(writer, string(utils.JsonStatus("Fail")))
			return
		}
		if !transactionReq.Validate() {
			log.Println("Error: missing fields")
			io.WriteString(writer, string(utils.JsonStatus("Fail")))
			return
		}

		publicKey := utils.PublicKeyFromString(*transactionReq.SenderPublicKey)
		privateKey := utils.PrivateKeyFromString(*transactionReq.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*transactionReq.Value, 32)
		if err != nil {
			log.Println("Error: parse error")
			io.WriteString(writer, string(utils.JsonStatus("Fail")))
			return
		}
		value32 := float32(value)

		writer.Header().Add("Content-Type", "application/json")

		transaction := wallet.NewTransaction(privateKey, publicKey,
			*transactionReq.SenderBlockchainAddress, *transactionReq.RecipientBlockchainAddress, value32)
		signature := transaction.GenerateSignature().String()

		blockTransaction := &chain.TransactionRequest{
			SenderBlockchainAddress:    transactionReq.SenderBlockchainAddress,
			RecipientBlockchainAddress: transactionReq.RecipientBlockchainAddress,
			SenderPublicKey:            transactionReq.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signature,
		}
		marshal, _ := json.Marshal(blockTransaction)
		buf := bytes.NewBuffer(marshal)

		res, _ := http.Post(server.Gateway()+"/transactions", "application/json", buf)
		if res.StatusCode == 201 {
			io.WriteString(writer, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(writer, string(utils.JsonStatus("fail")))

	default:
		writer.WriteHeader(http.StatusBadRequest)
		log.Printf("Error: Invalid Http Method")
	}
}

func (server *WalletServer) Run() {
	http.HandleFunc("/", server.Index)
	http.HandleFunc("/wallet", server.Wallet)
	http.HandleFunc("/transaction", server.CreateTransaction)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(server.Port())), nil))
}
