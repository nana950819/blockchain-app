/*
	For testing.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
	Note  : Before executing this, comment out the lock and unlock in
	        PeerList.Rebalance().
*/
package main

import (
	"cs686/cs686-blockchain-p3-kayfuku/p1"
	"cs686/cs686-blockchain-p3-kayfuku/p2"
	"cs686/cs686-blockchain-p3-kayfuku/p3"
	"cs686/cs686-blockchain-p3-kayfuku/p3/data"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/sha3"
)

func Test(t *testing.T) {

	// Test chooseClosest()
	fmt.Println("Test chooseClosest() start. ")
	fmt.Println("chooseClosest(list, 10, 4): ")
	list := []int{1, 3, 5, 7, 9, 11}
	ret := data.Test_chooseClosest(list, 10, 4)
	b := reflect.DeepEqual(ret, []int{1, 7, 9, 11})
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %d, but was %d", []int{1, 7, 9, 11}, ret)
	}

	fmt.Println("chooseClosest(list, 2, 4): ")
	ret = data.Test_chooseClosest(list, 2, 4)
	b = reflect.DeepEqual(ret, []int{1, 3, 5, 11})
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %d, but was %d", []int{1, 3, 5, 11}, ret)
	}

	fmt.Println("chooseClosest(list, 0, 4): ")
	ret = data.Test_chooseClosest(list, 0, 4)
	b = reflect.DeepEqual(ret, []int{1, 3, 9, 11})
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %d, but was %d", []int{1, 3, 9, 11}, ret)
	}

	fmt.Println("chooseClosest(list, 12, 4): ")
	ret = data.Test_chooseClosest(list, 12, 4)
	b = reflect.DeepEqual(ret, []int{1, 3, 9, 11})
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %d, but was %d", []int{1, 3, 9, 11}, ret)
	}

	fmt.Println("chooseClosest(list, 4, 6): ")
	ret = data.Test_chooseClosest(list, 4, 6)
	b = reflect.DeepEqual(ret, []int{1, 3, 5, 7, 9, 11})
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %d, but was %d", []int{1, 3, 5, 7, 9, 11}, ret)
	}
	fmt.Println()

	fmt.Println("Test PeerList.Rebalance() start.")
	fmt.Println("Please comment out the lock and unlock in PeerList.Rebalance().")
	data.TestPeerListRebalance()
	fmt.Println()

	// Test PeerList.PeerMapToJson().
	fmt.Println("Test PeerList.PeerMapToJson() start.")
	pl := data.NewPeerList(1, 32)
	pl.Add("addr1", 1)
	pl.Add("addr2", 2)
	jsonStr, _ := pl.PeerMapToJson()
	fmt.Println("jsonStr:", jsonStr) // jsonStr: {"addr1":1,"addr2":2}
	fmt.Println()

	// Test PeerList.Show()
	fmt.Println("Test PeerList.Show() start.")
	pl.Show()
	fmt.Println()

	// Test PeerList.JsonToPeerMap().
	fmt.Println("Test PeerList.JsonToPeerMap() start.")
	peerMap := pl.JsonToPeerMap(jsonStr)
	for k, v := range peerMap {
		fmt.Printf("%s %d\n", k, v)
	}
	// addr1 1
	// addr2 2

	// Test GetBlock().
	fmt.Println("Test SyncBlockChain.GetBlock() start.")
	sbc := data.NewBlockChain()
	height1blockList := sbc.GetLatestBlocks()
	firstBlock := height1blockList[0]
	mpt := p1.MerklePatriciaTrie{}
	b2 := sbc.GenBlock(mpt)
	sbc.Insert(b2)
	b3 := sbc.GenBlock(mpt)
	sbc.Insert(b3)
	jsonBlockStr2, _ := p2.EncodeToJson(b2)
	jsonBlockStr3, _ := p2.EncodeToJson(b3)
	fmt.Println("b2:", jsonBlockStr2)
	fmt.Println("b3:", jsonBlockStr3)
	block2, _ := sbc.GetBlock(2, b2.Header.Hash)
	jsonBlockStr2_retrieved, _ := p2.EncodeToJson(block2)
	fmt.Println("b2_retrieved:", jsonBlockStr2_retrieved)

	block3, _ := sbc.GetBlock(3, b3.Header.Hash)
	jsonBlockStr3_retrieved, _ := p2.EncodeToJson(block3)
	fmt.Println("b3_retrieved:", jsonBlockStr3_retrieved)
	fmt.Println()

	// Test GetLatestBlocks().
	fmt.Println("Test SyncBlockChain.GetLatestBlocks() start.")
	blockList := sbc.GetLatestBlocks()
	retBlock, _ := p2.EncodeToJson(blockList[0])
	b = reflect.DeepEqual(retBlock, jsonBlockStr3)
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %s, but was %s", jsonBlockStr3, retBlock)
	}

	// Test GetParentBlock().
	fmt.Println("Test SyncBlockChain.GetParentBlock() start.")
	parentBlock, _ := sbc.GetParentBlock(block3)
	retBlock, _ = p2.EncodeToJson(parentBlock)
	b = reflect.DeepEqual(retBlock, jsonBlockStr2)
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %s, but was %s", jsonBlockStr2, retBlock)
	}
	parentBlock, _ = sbc.GetParentBlock(block2)
	retBlock, _ = p2.EncodeToJson(parentBlock)
	mpt = p1.MerklePatriciaTrie{}
	mpt.Initial()
	jFirstBlock, _ := p2.EncodeToJson(firstBlock)
	b = reflect.DeepEqual(retBlock, jFirstBlock)
	fmt.Println(b)
	if !b {
		t.Errorf("Expected %s, but was %s", jFirstBlock, retBlock)
	}
	parentBlock, err := sbc.GetParentBlock(firstBlock)
	b = err != nil
	fmt.Println(b)
	fmt.Println(err)
	if !b {
		t.Errorf("Expected %s, but was %s", err, "nil")
	}

	// Test UpdateEntireBlockChain().
	fmt.Println("Test SyncBlockChain.UpdateEntireBlockChain() start.")
	jsonBlockChainStr, _ := sbc.BlockChainToJson()
	fmt.Println("jsonBlockChainStr:", jsonBlockChainStr)
	fmt.Println("UpdateEntireBlockChain() start.")
	sbc2 := data.NewBlockChain()
	sbc2.UpdateEntireBlockChain(jsonBlockChainStr)
	jsonBlockChainStr2, _ := sbc2.BlockChainToJson()
	fmt.Println("jsonBlockChainStr2:", jsonBlockChainStr2)
	fmt.Println()

	// // *** Test cases from Project 2. start. ***
	// // Comment out other parts of testing before doing this part.

	// fmt.Println("Test cases from Project 2.")
	// mpt := p1.MerklePatriciaTrie{}
	// mpt.Initial()
	// mpt.Insert("hello", "world")
	// mpt.Insert("charles", "ge")
	// b1 := p2.NewBlock(2, 1234567890, "genesis", mpt)
	// b2 := p2.NewBlock(3, 1234567890, b1.Header.Hash, mpt)

	// fmt.Println("b1")
	// fmt.Printf("Hash: %s\n", b1.Header.Hash)
	// fmt.Printf("Timestamp: %d\n", b1.Header.Timestamp)
	// fmt.Printf("Height: %d\n", b1.Header.Height)
	// fmt.Printf("ParentHash: %s\n", b1.Header.ParentHash)
	// fmt.Printf("Size: %d\n", b1.Header.Size)

	// fmt.Println("b2")
	// fmt.Printf("Hash: %s\n", b2.Header.Hash)
	// fmt.Printf("Timestamp: %d\n", b2.Header.Timestamp)
	// fmt.Printf("Height: %d\n", b2.Header.Height)
	// fmt.Printf("ParentHash: %s\n", b2.Header.ParentHash)
	// fmt.Printf("Size: %d\n", b2.Header.Size)

	// // Test DecodeFromJson(): a JSON string to a Block.
	// jsonString := "{\"hash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"timeStamp\": 1234567890, \"height\": 1, \"parentHash\": \"genesis\", \"size\": 1174, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}"
	// block := p2.DecodeFromJson(jsonString)
	// // fmt.Printf("block: %+v\n", block)
	// fmt.Println("block")
	// fmt.Printf("Hash: %s\n", block.Header.Hash)
	// fmt.Printf("Timestamp: %d\n", block.Header.Timestamp)
	// fmt.Printf("Height: %d\n", block.Header.Height)
	// fmt.Printf("ParentHash: %s\n", block.Header.ParentHash)
	// fmt.Printf("Size: %d\n", block.Header.Size)
	// str, _ := block.Value.Get("hello")
	// fmt.Println("str: ", str)
	// str, _ = block.Value.Get("charles")
	// fmt.Println("str: ", str)

	// // Test DecodeFromJson(): a Block to a JSON string.
	// jsonString, _ = p2.EncodeToJson(block)
	// fmt.Println("jsonString: ", jsonString)

	// // Test bc.EncodeToJson: Blocks to JSON strings.
	// bc := p2.NewBlockChain()
	// bc.Insert(b1)
	// bc.Insert(b2)
	// jsonStrings, err := bc.EncodeToJson()
	// if err != nil {
	// 	fmt.Println(err)
	// 	t.Fail()
	// }
	// fmt.Println("jsonStrings: ", jsonStrings)
	// ret := bc.Length
	// fmt.Println(ret) // 3
	// if ret != 3 {
	// 	t.Errorf("Expected %d, but was %d", 3, ret)
	// }
	// fmt.Println()
	// // Test cases from Project 2. end.

	// Test random nonce.
	fmt.Println("Test random nonce start.")
	nonce := getRandHexStr(16)
	fmt.Println("nonce:", nonce)
	mpt = p1.MerklePatriciaTrie{}
	nonceHash := CalcNonceHash(b3.Header.Hash, mpt.Root)
	fmt.Println("nonceHash:", nonceHash)

	// Test time.
	timeStamp := time.Now().UnixNano()
	fmt.Println("time:", timeStamp)
	fmt.Println()

	// // Test Proof of Work.
	// fmt.Println("Test Proof of Work start...")
	// for {
	// 	nonceHash := CalcNonceHash(b3.Header.Hash, mpt.Root)
	// 	if strings.HasPrefix(nonceHash, "000000") {
	// 		fmt.Println("Nonce found! nonceHash:", nonceHash)
	// 		break
	// 	}
	// }
	// fmt.Println()

	// Test Canonical().
	fmt.Println("Test Canonical() start.")
	p3.SBC = sbc
	highestList := p3.SBC.GetLatestBlocks()
	chainNum := 0
	for _, b := range highestList {
		chainNum++
		fmt.Printf("Chain #%d:\n", chainNum)
		p3.DisplayAncestors(b)
	}
	fmt.Println()
	// Fork version.
	fmt.Println("Fork version: ")
	timeStamp = time.Now().UnixNano()
	b3_2 := p2.NewBlock(3, timeStamp, b2.Header.Hash, mpt, "")
	p3.SBC.Insert(b3_2)
	highestList = p3.SBC.GetLatestBlocks()
	chainNum = 0
	for _, b := range highestList {
		chainNum++
		fmt.Printf("Chain #%d:\n", chainNum)
		p3.DisplayAncestors(b)
	}
	fmt.Println()

	// Test rentalInfo.EncodeToJson().
	fmt.Println("Test rentalInfo.EncodeToJson() start.")
	rentalInfo1 := p1.NewRentalInfo(123, 100)
	jRentalInfo, _ := rentalInfo1.EncodeToJson()
	fmt.Println("jRentalInfo:", jRentalInfo)
	fmt.Println()

	// Test rentalInfos.EncodeToJson().
	fmt.Println("Test rentalInfos.EncodeToJson() start.")
	var rentalInfos p1.RentalInfos
	rentalInfos = append(rentalInfos, rentalInfo1)
	rentalInfo2 := p1.NewRentalInfo(333, 200)
	rentalInfos = append(rentalInfos, rentalInfo2)
	jRentalInfos, _ := rentalInfos.EncodeToJson()
	fmt.Println("jRentalInfos:", jRentalInfos)
	fmt.Println()

	// Test EncodeMapToJson().
	fmt.Println("Test EncodeMapToJson() start.")
	rentalInfo := p1.NewRentalInfo(200, 50)
	rentalMap := make(map[string]p1.RentalInfo)
	rentalMap[rentalInfo.RentalId] = rentalInfo
	jRentalMap, _ := p1.EncodeMapToJson(rentalMap)
	fmt.Println("jRentalMap:", jRentalMap)
	fmt.Println()

	// Test DecodeJsonToMap().
	fmt.Println("Test DecodeJsonToMap() start.")
	ret_rentalMap, _ := data.DecodeJsonToMap(jRentalMap)
	for k, v := range ret_rentalMap {
		fmt.Println("k:", k)
		fmt.Printf("v.LenderId: %d, v.AskingPrice: %d, v.IsAvailable: %v", v.LenderId, v.AskingPrice, v.IsAvailable)
	}
	fmt.Println()

	fmt.Println()
	fmt.Println("done.")
}

// BlockJson For testing.
type BlockJson struct {
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Hash       string            `json:"hash"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	MPT        map[string]string `json:"mpt"`
}

// Calculate a random nonceHash.
func CalcNonceHash(parentHash string, mptRoot string) string {
	var str string
	nonce := getRandHexStr(16)
	str = parentHash + nonce + mptRoot
	sum := sha3.Sum256([]byte(str))
	nonceHash := hex.EncodeToString(sum[:])
	return nonceHash
}

// Generate a random hex string of fixed length.
// https://stackoverflow.com/questions/46904588/efficient-way-to-to-generate-a-random-hex-string-of-a-fixed-length-in-golang
const letterBytes = "abcdef0123456789"
const (
	letterIdxBits = 4                    // 4 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src1 = rand.NewSource(time.Now().UnixNano())
var src2 = rand.New(rand.NewSource(time.Now().UnixNano()))

// getRandHexStr returns a random hexadecimal string of length n.
// https://stackoverflow.com/questions/46904588/efficient-way-to-to-generate-a-random-hex-string-of-a-fixed-length-in-golang
func getRandHexStr(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src1.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src1.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
