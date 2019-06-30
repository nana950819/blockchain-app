/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	Hold information for RegisterData.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package data

import (
	"encoding/json"
	"fmt"
)

// Hold the assigned ID and JSON string of peerMap.
type RegisterData struct {
	AssignedId  int32  `json:"assignedId"`
	PeerMapJson string `json:"peerMapJson"`
	PortNum     string `json:"portNum"`
}

// Create a new RegisterData object.
func NewRegisterData(id int32, peerMapJson string) RegisterData {
	rd := RegisterData{AssignedId: id, PeerMapJson: peerMapJson}
	return rd
}

// Convert the RegisterData object to JSON string.
func (data *RegisterData) EncodeToJson() (string, error) {
	rd := RegisterData{}
	jsonString, err := json.Marshal(rd)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	return string(jsonString), nil
}
