/*
	CS686 Project 3 and 4, Block in Blockchain.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package p2

import (
	"bytes"
	"cs686/cs686-blockchain-p3-kayfuku/p1"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"golang.org/x/crypto/sha3"
)

// Hold header of a block.
type Header struct {
	Height     int32  `json:"height"`
	Timestamp  int64  `json:"timeStamp"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Size       int32  `json:"size"`
	Nonce      string `json:"nonce"`
}

// A block. Hold header and value.
type Block struct {
	Header Header                `json:"header"`
	Value  p1.MerklePatriciaTrie `json:"value"`
}

// Create a new block.
func NewBlock(height int32, timestamp int64, parentHash string, mpt p1.MerklePatriciaTrie, nonce string) Block {
	block := Block{}
	block.Initial(height, timestamp, parentHash, mpt)
	block.Header.Nonce = nonce
	block.Header.Hash = block.hash_block()
	return block
}

// Initialize the block.
func (b *Block) Initial(height int32, timestamp int64, parentHash string, mpt p1.MerklePatriciaTrie) {
	header := Header{}
	header.Height = height
	header.Timestamp = timestamp
	header.ParentHash = parentHash

	b.Value = mpt
	header.Size = int32(len(convertToByteArray(b.Value)))

	b.Header = header
}

// Get the hash of the block.
func (b *Block) hash_block() string {
	hash_str := strconv.Itoa(int(b.Header.Height)) + strconv.Itoa(int(b.Header.Timestamp)) + b.Header.ParentHash +
		b.Value.Root + strconv.Itoa(int(b.Header.Size))

	sum := sha3.Sum256([]byte(hash_str))
	return hex.EncodeToString(sum[:])
}

// Get byte array of MerklePatriciaTrie.
func convertToByteArray(mpt p1.MerklePatriciaTrie) []uint8 {
	encBuf := new(bytes.Buffer)
	err := gob.NewEncoder(encBuf).Encode(&mpt)
	if err != nil {
		log.Fatal(err)
	}

	value := encBuf.Bytes()
	return value
}

// Convert Block's JSON string to Block.
func DecodeFromJson(jsonString string) Block {

	bytes := []byte(jsonString)
	var bj BlockJson
	err := json.Unmarshal(bytes, &bj)
	if err != nil {
		fmt.Println("error:", err)
	}

	block := BuildBlockFromBj(bj)

	return block
}

// This is used when converting Block to JSON string and vice versa.
type BlockJson struct {
	Hash       string            `json:"hash"`
	TimeStamp  int64             `json:"timeStamp"`
	Height     int32             `json:"height"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	Mpt        map[string]string `json:"mpt"`
	Nonce      string            `json:"nonce"`
}

// Convert BlockJson to Block.
func BuildBlockFromBj(bj BlockJson) Block {
	header := Header{}
	header.Hash = bj.Hash
	header.Timestamp = bj.TimeStamp
	header.Height = bj.Height
	header.ParentHash = bj.ParentHash
	header.Size = bj.Size
	header.Nonce = bj.Nonce

	mptMap := bj.Mpt
	mpt := buildMpt(mptMap)

	var block Block
	block.Header = header
	block.Value = mpt

	return block
}

// Take a map, take out every entry in the map, and insert them into MerklePatriciaTrie.
func buildMpt(mptMap map[string]string) p1.MerklePatriciaTrie {
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()

	mpt.Insert("ads", mptMap["ads"])
	mpt.Insert("history", mptMap["history"])

	// for key, value := range mptMap {
	// 	fmt.Printf("key: %s, value: %s\n", key, value)
	// 	mpt.Insert(key, value)
	// }

	fmt.Println("mptRoot:", mpt.Root)

	return mpt
}

// Convert Block to Block's JSON string.
func EncodeToJson(block Block) (string, error) {
	bj := BuildBjFromBlock(block)

	jsonString, err := json.Marshal(bj)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(jsonString), nil
}

// Convert Block to BlockJson.
func BuildBjFromBlock(block Block) BlockJson {
	bj := BlockJson{}
	header := block.Header

	bj.Hash = header.Hash
	bj.TimeStamp = header.Timestamp
	bj.Height = header.Height
	bj.ParentHash = header.ParentHash
	bj.Size = header.Size
	bj.Nonce = header.Nonce
	bj.Mpt = make(map[string]string)

	mptMap := block.Value.ValueDb
	for k, v := range mptMap {
		bj.Mpt[k] = v
	}

	return bj
}
