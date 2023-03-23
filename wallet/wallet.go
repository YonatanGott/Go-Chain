package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"go-block/utils"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

type Transaction struct {
	senderPrivateKey   *ecdsa.PrivateKey
	recipientPublicKey *ecdsa.PublicKey
	senderAddress      string
	recipientAddress   string
	value              float32
}
type TransactionRequest struct {
	SenderPrivateKey           *string `json:"senderPrivatekey"`
	SenderBlockchainAddress    *string `json:"senderBlockchainAddress"`
	RecipientBlockchainAddress *string `json:"recipientBlockchainAddress"`
	SenderPublicKey            *string `json:"senderPublicKey"`
	Value                      *string `json:"value"`
}

func NewWallet() *Wallet {
	// ECDSA private key and public key (32 bytes)
	wallet := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	wallet.privateKey = privateKey
	wallet.publicKey = &wallet.privateKey.PublicKey
	// SHA-256 hashing on public key (32 bytes)
	hash2 := sha256.New()
	hash2.Write(wallet.publicKey.X.Bytes())
	hash2.Write(wallet.publicKey.Y.Bytes())
	digest2 := hash2.Sum(nil)
	// RIPEMD-160 hashing on the sha-256 hash (20 bytes)
	hash3 := ripemd160.New()
	hash3.Write(digest2)
	digest3 := hash3.Sum(nil)
	// Add version byte in front of ripemd-160 hash ( 0x00 for main )
	versionByte4 := make([]byte, 21)
	versionByte4[0] = 0x00
	copy(versionByte4[1:], digest3[:])
	// SHA-256 hashing on the extended ripemd-160 result
	hash5 := sha256.New()
	hash5.Write(versionByte4)
	digest5 := hash5.Sum(nil)
	// SHA-256 hashing on the previous sha-256 result
	hash6 := sha256.New()
	hash6.Write(digest5)
	digest6 := hash6.Sum(nil)
	// First 4 byes of the second sha-256 hash for checksum
	checksum := digest6[:4]
	// Add the 4 checksum bytes at the end of the ripemd-160 hash (25 bytes)
	checksumBytes := make([]byte, 25)
	copy(checksumBytes[:21], versionByte4[:])
	copy(checksumBytes[21:], checksum[:])
	// convert result from byte string to base58
	address := base58.Encode(checksumBytes)
	wallet.blockchainAddress = address
	return wallet
}

func (wallet *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", wallet.privateKey.D.Bytes())
}

func (wallet *Wallet) PublicKey() *ecdsa.PublicKey {
	return &wallet.privateKey.PublicKey
}

func (wallet *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", wallet.publicKey.X.Bytes(), wallet.publicKey.Y.Bytes())
}

func (wallet *Wallet) BlockchainAddress() string {
	return wallet.blockchainAddress
}

func (wallet *Wallet) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PrivateKey        string `json:"privateKey"`
		PublicKey         string `json:"publicKey"`
		BlockchainAddress string `json:"blockchainAddress"`
	}{
		PrivateKey:        wallet.PrivateKeyStr(),
		PublicKey:         wallet.PublicKeyStr(),
		BlockchainAddress: wallet.BlockchainAddress(),
	})
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, sender string, recipient string, value float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, recipient, value}
}

func (transaction *Transaction) GenerateSignature() *utils.Signature {
	marshal, _ := json.Marshal(transaction)
	hash := sha256.Sum256([]byte(marshal))
	r, s, _ := ecdsa.Sign(rand.Reader, transaction.senderPrivateKey, hash[:])
	return &utils.Signature{R: r, S: s}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"senderAddress"`
		Recipient string  `json:"recipientAddress"`
		Value     float32 `json:"value"`
	}{
		Sender:    transaction.senderAddress,
		Recipient: transaction.recipientAddress,
		Value:     transaction.value,
	})
}

func (transactionReq *TransactionRequest) Validate() bool {
	if transactionReq.SenderPrivateKey == nil ||
		transactionReq.RecipientBlockchainAddress == nil ||
		transactionReq.SenderBlockchainAddress == nil ||
		transactionReq.SenderPublicKey == nil ||
		transactionReq.Value == nil {
		return false
	}
	return true
}
