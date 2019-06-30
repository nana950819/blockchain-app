/*
	CS686 Project 5: Build a blockchain application.
	Hold the history information.
	Author: Kei Fukutani
	Date  : May 13th, 2019
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

// Hold the content of the history information.
type History struct {
	HistoryId  string    `json:"historyId"`
	RentalId   string    `json:"rentalId"`
	BorrowerId int32     `json:"borrowerId"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Deposit    int32     `json:"deposit"`
	State      int8      `json:"state"`
}

// Create a new History.
func NewHistory(rentalId string, borrowerId int32, deposit int32) History {
	timeStamp := time.Now().UnixNano()

	// Generate history ID.
	str := strconv.Itoa(int(timeStamp)) + rentalId + strconv.Itoa(int(borrowerId))
	sum := sha3.Sum256([]byte(str))
	historyId := hex.EncodeToString(sum[:])

	history := History{HistoryId: historyId, RentalId: rentalId, BorrowerId: borrowerId,
		StartTime: time.Time{}, EndTime: time.Time{}, Deposit: deposit, State: 0}

	return history
}

// Convert History to JSON string.
func (history *History) EncodeToJson() (string, error) {
	jHistory, err := json.Marshal(history)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(jHistory), nil
}

// Convert History's JSON string to History.
func DecodeJsonToHistory(jHistory string) History {
	bytes := []byte(jHistory)
	var history History
	err := json.Unmarshal(bytes, &history)
	if err != nil {
		fmt.Println("error:", err)
	}
	return history
}

// Convert historyMap to JSON string.
func EncodeHistoryMapToJson(historyMap map[string]History) (string, error) {
	jHistoryMap, err := json.Marshal(historyMap)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(jHistoryMap), nil
}

// Convert JSON string to historyMap.
func DecodeJsonToHistoryMap(jHistoryMap string) (map[string]History, error) {
	bytes := []byte(jHistoryMap)
	var historyMap map[string]History
	err := json.Unmarshal(bytes, &historyMap)
	if err != nil {
		fmt.Println("error:", err)
	}

	return historyMap, nil
}
