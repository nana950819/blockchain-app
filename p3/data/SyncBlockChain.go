/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	Hold BlockChain with synchronization.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package data

import (
	"cs686/cs686-blockchain-p3-kayfuku/p1"
	"cs686/cs686-blockchain-p3-kayfuku/p2"
	"encoding/json"
	"fmt"
	"sync"
)

// Hold BlockChain with synchronization.
type SyncBlockChain struct {
	bc p2.BlockChain
	// lock
	mux sync.Mutex
}

// Create a new blockchain.
func NewBlockChain() SyncBlockChain {
	blockChain := p2.NewBlockChain()
	sbc := SyncBlockChain{bc: blockChain}
	return sbc
}

// Get the list of blocks at the height.
func (sbc *SyncBlockChain) Get(height int32) ([]p2.Block, bool) {
	sbc.mux.Lock()

	blockList, b := sbc.bc.Get(height)
	if !b {
		return nil, false
	}

	defer sbc.mux.Unlock()
	return blockList, true
}

// Get the block with the height and the hash.
func (sbc *SyncBlockChain) GetBlock(height int32, hash string) (p2.Block, bool) {
	// Do not put lock and unlock in this function!

	blockList, ok := sbc.Get(height)
	if !ok {
		return p2.Block{}, false
	}

	for _, block := range blockList {
		if block.Header.Hash == hash {
			return block, true
		}
	}

	return p2.Block{}, false
}

// This function returns the list of blocks of height "BlockChain.length".
func (sbc *SyncBlockChain) GetLatestBlocks() []p2.Block {
	sbc.mux.Lock()
	list := sbc.bc.GetLatestBlocks()
	defer sbc.mux.Unlock()
	return list
}

// This function takes a block as the parameter, and returns its parent block.
func (sbc *SyncBlockChain) GetParentBlock(block p2.Block) (p2.Block, error) {
	sbc.mux.Lock()
	block, err := sbc.bc.GetParentBlock(block)
	if err != nil {
		defer sbc.mux.Unlock()
		return p2.Block{}, err
	}
	defer sbc.mux.Unlock()
	return block, nil
}

// Insert the block to the BlockChain.
func (sbc *SyncBlockChain) Insert(block p2.Block) bool {
	sbc.mux.Lock()
	b := sbc.bc.Insert(block)
	defer sbc.mux.Unlock()
	return b
}

// Check if the block's parent exists in the current BlockChain.
func (sbc *SyncBlockChain) CheckParentHash(insertBlock p2.Block) bool {
	// fmt.Println("CheckParentHash() start.")
	heightToBlocks := sbc.bc.Chain
	height := insertBlock.Header.Height
	previousList := heightToBlocks[height-1]
	parentHash := insertBlock.Header.ParentHash
	// fmt.Println("previousList:", previousList)
	b := p2.ExistsInList(parentHash, previousList)

	return b
}

// Update MPT.
func (sbc *SyncBlockChain) UpdateMPT(insertBlock p2.Block) p2.Block {
	sbc.mux.Lock()
	fmt.Println("UpdateMPT() start.")
	heightToBlocks := sbc.bc.Chain
	height := insertBlock.Header.Height
	previousList := heightToBlocks[height-1]
	parentHash := insertBlock.Header.ParentHash

	var parentMpt p1.MerklePatriciaTrie
	for _, parentBlock := range previousList {
		if parentBlock.Header.Hash == parentHash {
			fmt.Println("parent block found. ")
			parentMpt = parentBlock.Value
		}
	}

	parentMptMap := parentMpt.ValueDb
	// Test
	fmt.Println("parentMptMap:")
	for k, v := range parentMptMap {
		fmt.Println("k:", k)
		fmt.Println("v:", v)
	}
	transactionMap := insertBlock.Value.ValueDb
	// Test
	fmt.Println("transactionMap:")
	for k, v := range transactionMap {
		fmt.Println("k:", k)
		fmt.Println("v:", v)
	}

	for txK, txV := range transactionMap {
		if txK == "ads" {
			parentRentalMap, _ := p1.DecodeJsonToRentalMap(parentMptMap[txK])
			txRentalMap, _ := p1.DecodeJsonToRentalMap(txV)
			fmt.Println("txRentalMap:")
			for k, v := range txRentalMap {
				fmt.Println("Rental ID:", k)
				fmt.Printf("Lender ID: %d, Asking Price: %d, Availability: %v\n", v.LenderId, v.AskingPrice, v.IsAvailable)
				parentRentalMap[k] = v
			}
			jParentRentalMap, _ := p1.EncodeRentalMapToJson(parentRentalMap)
			parentMpt.Insert(txK, jParentRentalMap)
		}
		if txK == "history" {
			parentHistoryMap, _ := p1.DecodeJsonToHistoryMap(parentMptMap[txK])
			txHistoryMap, _ := p1.DecodeJsonToHistoryMap(txV)
			fmt.Println("txHistoryMap:")
			for k, v := range txHistoryMap {
				fmt.Println("History ID:", k)
				fmt.Printf("Rental ID: %s, Borrower ID: %d, Start Time: %v, End Time: %v, Deposit: %d\n",
					v.RentalId, v.BorrowerId, v.StartTime, v.EndTime, v.Deposit)
				parentHistoryMap[k] = v
			}
			jParentHistoryMap, _ := p1.EncodeHistoryMapToJson(parentHistoryMap)
			parentMpt.Insert(txK, jParentHistoryMap)
		}
	}
	insertBlock.Value = parentMpt

	defer sbc.mux.Unlock()
	return insertBlock
}

// Take BlockChain's JSON strings and build BlockChain.
func (sbc *SyncBlockChain) UpdateEntireBlockChain(blockChainJson string) error {
	sbc.mux.Lock()

	err := sbc.bc.DecodeFromJson(blockChainJson)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	defer sbc.mux.Unlock()
	return nil
}

// Return BlockChain's Json string.
func (sbc *SyncBlockChain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	blockJsonList := []p2.BlockJson{}

	heightToBlocks := sbc.bc.Chain
	for _, list := range heightToBlocks {
		for _, b := range list {
			bj := p2.BuildBjFromBlock(b)
			blockJsonList = append(blockJsonList, bj)
		}
	}

	jsonStrings, err := json.Marshal(blockJsonList)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	defer sbc.mux.Unlock()
	return string(jsonStrings), nil
}

// Create a block at the highest in the current BlockChain.
func (sbc *SyncBlockChain) GenBlock(mpt p1.MerklePatriciaTrie) p2.Block {
	heighest := sbc.bc.Length
	listBlock, _ := sbc.Get(heighest)
	parentHash := listBlock[0].Header.Hash
	block := p2.NewBlock(heighest+1, 123456789, parentHash, mpt, "")
	return block
}

// Show the hash of the BlockChain.
func (sbc *SyncBlockChain) Show() string {
	return sbc.bc.Show()
}
