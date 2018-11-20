package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type Transaction struct {
	From string
	To  string
	BTC  float64
	Reason string
}

var Transactions []Transaction

func Tsend(transaction Transaction)  {

	Transactions = append(Transactions, transaction)
}

func Reset(){
	i :=len(Transactions)
	Transactions=append(Transactions[:i])
}

func PackTransaction() string  {
	var t string
	for _,s := range Transactions {
		temp := [][]byte{
			FloatToByte(s.BTC),
			[]byte(s.From),
			[]byte(s.To),
			[]byte(s.Reason),
		}


	data := bytes.Join(temp,[]byte{})
	hash :=sha256.Sum256(data)

	t1 :=hex.EncodeToString(hash[:])
	t=strings.Join([]string{t,t1},"")
	}
	value := sha256.Sum256([]byte(t))
	return hex.EncodeToString(value[:])
}
