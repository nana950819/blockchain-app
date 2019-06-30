/*
	CS686 Project 5: Build a blockchain application.
	Hold the rental information.
	Author: Kei Fukutani
	Date  : May 10th, 2019
*/
package p1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

// Hold the content of the rental information.
type RentalInfo struct {
	RentalId    string `json:"rentalId"`
	LenderId    int32  `json:"lenderId"`
	AskingPrice int32  `json:"askingPrice"`
	IsAvailable bool   `json:"isAvailable"`
}

// Create a new RentalInfo.
func NewRentalInfo(lenderId int32, askingPrice int32) RentalInfo {
	timeStamp := time.Now().UnixNano()

	// Generate rental ID.
	str := strconv.Itoa(int(timeStamp)) + strconv.Itoa(int(lenderId)) + strconv.Itoa(int(askingPrice))
	sum := sha3.Sum256([]byte(str))
	rentalId := hex.EncodeToString(sum[:])

	rentalInfo := RentalInfo{RentalId: rentalId, LenderId: lenderId, AskingPrice: askingPrice, IsAvailable: true}

	return rentalInfo
}

// Convert RentalInfo to JSON string.
func (rentalInfo *RentalInfo) EncodeToJson() (string, error) {
	jRentalInfo, err := json.Marshal(rentalInfo)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(jRentalInfo), nil
}

// Convert RentalInfo's JSON string to RentalInfo.
func DecodeJsonToRentalInfo(jRentalInfo string) RentalInfo {
	bytes := []byte(jRentalInfo)
	var rentalInfo RentalInfo
	err := json.Unmarshal(bytes, &rentalInfo)
	if err != nil {
		fmt.Println("error:", err)
	}
	return rentalInfo
}

// Convert rentalMap to JSON string.
func EncodeRentalMapToJson(rentalMap map[string]RentalInfo) (string, error) {
	jRentalMap, err := json.Marshal(rentalMap)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(jRentalMap), nil
}

// Convert JSON string to rentalMap.
func DecodeJsonToRentalMap(jRentalMap string) (map[string]RentalInfo, error) {
	bytes := []byte(jRentalMap)
	var rentalMap map[string]RentalInfo
	err := json.Unmarshal(bytes, &rentalMap)
	if err != nil {
		fmt.Println("error:", err)
	}

	return rentalMap, nil
}
