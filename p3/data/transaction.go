/*
	CS686 Project 5: Build a blockchain application.
	Hold the transaction information.
	Author: Kei Fukutani
	Date  : May 12th, 2019
*/
package data

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

// Hold the content of the transaction information.
type Transaction struct {
	TxId           string `json:"txId"`
	RentalInfoJson string `json:"rentalInfoJson"`
	HistoryJson    string `json:"historyJson"`
	Sig            string `json:"sig"`
	TxFee          uint32 `json:"txFee"`
}

// Create a new Transaction.
func NewTx(jRentalInfo string, jHistory string, txFee uint32) Transaction {
	timeStamp := time.Now().UnixNano()

	// Generate rental ID.
	str := strconv.Itoa(int(timeStamp)) + jRentalInfo + jHistory + strconv.Itoa(int(txFee))
	sum := sha3.Sum256([]byte(str))
	txId := hex.EncodeToString(sum[:])

	transaction := Transaction{
		TxId:           txId,
		RentalInfoJson: jRentalInfo,
		HistoryJson:    jHistory,
		Sig:            "",
		TxFee:          txFee,
	}

	return transaction
}

// Convert Transaction to JSON string.
func (tx *Transaction) EncodeToJson() (string, error) {
	jTx, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(jTx), nil
}

// Convert Transaction's JSON string to Transaction.
func DecodeFromJson(jTx string) Transaction {
	bytes := []byte(jTx)
	var tx Transaction
	err := json.Unmarshal(bytes, &tx)
	if err != nil {
		fmt.Println("error:", err)
	}
	return tx
}
