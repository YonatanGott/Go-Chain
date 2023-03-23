package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (signature *Signature) String() string {
	return fmt.Sprintf("%064x%064x", signature.R, signature.S)
}

func StringToBigIntTuple(str string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(str[:64])
	by, _ := hex.DecodeString(str[64:])

	var bix big.Int
	var biy big.Int

	_ = bix.SetBytes(bx)
	_ = biy.SetBytes(by)

	return bix, biy
}

func PublicKeyFromString(str string) *ecdsa.PublicKey {
	x, y := StringToBigIntTuple(str)
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}
}

func PrivateKeyFromString(str string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	bytes, _ := hex.DecodeString(str[:])
	var bi big.Int
	_ = bi.SetBytes(bytes)
	return &ecdsa.PrivateKey{PublicKey: *publicKey, D: &bi}
}

func SignatureFromString(str string) *Signature {
	r, s := StringToBigIntTuple(str)
	return &Signature{R: &r, S: &s}
}
