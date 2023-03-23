package chain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go-block/utils"
)

const (
	MINING_ZEROS  = 3
	MINING_SENDER = "GO CHAIN"
	MINING_REWARD = 0.1
)

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
	chainAddress    string
	port            uint16
}

type Block struct {
	timestamp    int64
	nonce        int
	prevHash     [32]byte
	transactions []*Transaction
}

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress *string  `json:"recipientBlockchainAddress"`
	SenderPublicKey            *string  `json:"senderPublicKey"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (block *Block) Print() {
	fmt.Printf("Time            %d\n", block.timestamp)
	fmt.Printf("Nonce           %d\n", block.nonce)
	fmt.Printf("Previous_Hash   %x\n", block.prevHash)
	for _, t := range block.transactions {
		t.Print()
	}
}

func (block *Block) Hash() [32]byte {
	bytes, _ := block.MarshalJSON()
	return sha256.Sum256([]byte(bytes))
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"prevHash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    block.timestamp,
		Nonce:        block.nonce,
		PreviousHash: fmt.Sprintf("%x", block.prevHash),
		Transactions: block.transactions,
	})
}

func NewBlock(nonce int, prevHash [32]byte, transactions []*Transaction) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		nonce:        nonce,
		prevHash:     prevHash,
		transactions: transactions,
	}
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (transaction *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	fmt.Printf("Sender Address            %s\n", transaction.senderAddress)
	fmt.Printf("Recipient Address         %s\n", transaction.recipientAddress)
	fmt.Printf("Value   %.2f\n", transaction.value)
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"senderAddress"`
		RecipientAddress string  `json:"recipientAddress"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    transaction.senderAddress,
		RecipientAddress: transaction.recipientAddress,
		Value:            transaction.value,
	})
}

func (blockchain *Blockchain) Print() {
	for index, value := range blockchain.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), index, strings.Repeat("=", 25))
		value.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("-", 25))
}

func (blockchain *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, signature *utils.Signature) bool {
	transaction := NewTransaction(sender, recipient, value)
	if sender == MINING_SENDER {
		blockchain.transactionPool = append(blockchain.transactionPool, transaction)
		return true
	}
	if blockchain.VerifyTransactionSignature(senderPublicKey, signature, transaction) {
		// if blockchain.CalculateTotalAmount(sender) < value {
		// 	log.Println("Error: insufficient funds")
		// 	return false
		// }
		blockchain.transactionPool = append(blockchain.transactionPool, transaction)
		return true
	} else {
		log.Println("Error: Couldn't Verify")
	}
	return false
}

func (blockchain *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, signature *utils.Signature, transaction *Transaction) bool {
	marshal, _ := json.Marshal(transaction)
	hash := sha256.Sum256([]byte(marshal))
	return ecdsa.Verify(senderPublicKey, hash[:], signature.R, signature.S)
}

func (blockchain *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range blockchain.transactionPool {
		transactions = append(transactions, NewTransaction(
			t.senderAddress, t.recipientAddress, t.value,
		))
	}
	return transactions
}

func (blockchain *Blockchain) CreateNewBlock(nonce int, prevHash [32]byte) *Block {
	block := NewBlock(nonce, prevHash, blockchain.transactionPool)
	blockchain.chain = append(blockchain.chain, block)
	blockchain.transactionPool = []*Transaction{}
	return block
}

func (blockchain *Blockchain) LastBlock() *Block {
	return blockchain.chain[len(blockchain.chain)-1]
}

func (blockchain *Blockchain) ValidProof(nonce int, prevHash [32]byte, transactions []*Transaction, zeros int) bool {
	zerosLead := strings.Repeat("0", zeros)
	blockGuess := Block{0, nonce, prevHash, transactions}
	hashGuess := fmt.Sprintf("%x", blockGuess.Hash())
	return hashGuess[:zeros] == zerosLead
}

func (blockchain *Blockchain) Mining() bool {
	blockchain.AddTransaction(MINING_SENDER, blockchain.chainAddress, MINING_REWARD, nil, nil)
	nonce := blockchain.ProofOfWork()
	previousHash := blockchain.LastBlock().Hash()
	blockchain.CreateNewBlock(nonce, previousHash)
	log.Println("action=mining", "status=success")
	return true
}

func (blockchain *Blockchain) CalculateTotalAmount(address string) float32 {
	var totalAmount float32 = 0.0
	for _, block := range blockchain.chain {
		for _, transaction := range block.transactions {
			value := transaction.value
			if address == transaction.recipientAddress {
				totalAmount += value
			}
			if address == transaction.senderAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (blockchain *Blockchain) ProofOfWork() int {
	transactions := blockchain.CopyTransactionPool()
	previousHash := blockchain.LastBlock().Hash()
	nonce := 0
	for !blockchain.ValidProof(nonce, previousHash, transactions, MINING_ZEROS) {
		nonce += 1
	}
	return nonce
}

func NewBlockchain(chainAddress string, port uint16) *Blockchain {
	block := &Block{}
	blockchain := new(Blockchain)
	blockchain.chainAddress = chainAddress
	blockchain.CreateNewBlock(0, block.Hash())
	blockchain.port = port
	return blockchain
}

func (blockchain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chains"`
	}{
		Blocks: blockchain.chain,
	})
}

func (transactionReq *TransactionRequest) Validate() bool {
	if transactionReq.RecipientBlockchainAddress == nil ||
		transactionReq.SenderBlockchainAddress == nil ||
		transactionReq.SenderPublicKey == nil ||
		transactionReq.Value == nil ||
		transactionReq.Signature == nil {
		return false
	}
	return true
}
