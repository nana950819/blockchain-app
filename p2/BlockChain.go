/*
	CS686 Project 3 and 4, Blockchain.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package p2

import (
	"cs686/cs686-blockchain-p3-kayfuku/p1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"golang.org/x/crypto/sha3"
)

// Hold blockchain.
type BlockChain struct {
	// K: Block's height, V: List of Block.
	Chain map[int32][]Block
	// The highest block's height.
	Length int32
}

// Create a new blockchain with the first block.
func NewBlockChain() BlockChain {
	bc := BlockChain{}
	bc.Chain = make(map[int32][]Block)
	bc.Length = 0
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()
	block := NewBlock(1, 123456789, "genesis", mpt, "")
	bc.Insert(block)
	return bc
}

// This function takes a height as the argument,
// returns the list of blocks stored in that height or None
// if the height doesn't exist.
func (bc *BlockChain) Get(height int32) ([]Block, bool) {
	list, ok := bc.Chain[height]
	if !ok {
		return nil, false
	}
	return list, true
}

// This function returns the list of blocks of height "BlockChain.length".
func (bc *BlockChain) GetLatestBlocks() []Block {
	list, _ := bc.Get(bc.Length)
	return list
}

// This function takes a block as the argument, use its height
// to find the corresponding list in blockchain's Chain map.
// If the list has already contained that block's hash, ignore it
// because we don't store duplicate blocks; if not, insert the block into the list.
func (bc *BlockChain) Insert(block Block) bool {
	heightToBlocks := bc.Chain
	height := block.Header.Height
	// Commented out for the following project.
	// previousList := heightToBlocks[height-1]
	// parentHash := block.Header.ParentHash

	if height < 0 {
		return false
	}
	if height > bc.Length+1 {
		newList := []Block{}
		newList = append(newList, block)
		heightToBlocks[height] = newList
		bc.Length = height
		return true
	}
	if bc.Length == 0 {
		newList := []Block{}
		newList = append(newList, block)
		heightToBlocks[height] = newList
		bc.Length++
		return true
	}
	if height == bc.Length+1 {
		// Check if the parent exists. Commented out for the following project.
		// if ExistsInList(parentHash, previousList) {
		newList := []Block{}
		newList = append(newList, block)
		heightToBlocks[height] = newList
		bc.Length++
		return true
		// }
	}

	list := heightToBlocks[height]
	hash := block.Header.Hash
	// Check if it is in the list.
	if ExistsInList(hash, list) {
		// Ignore the block.
		fmt.Println("The block was ignored because it is already in the current BlockChain.")
		return false
	}
	// Check if the parent exists. Commented out for the following project.
	// if ExistsInList(parentHash, previousList) {
	list = append(list, block)
	heightToBlocks[height] = list
	return true
	// }

	// return false
}

// Take hash of a block and a list of blocks, and check to see if
// the block exists in the list.
func ExistsInList(hash string, list []Block) bool {
	for _, blockInList := range list {
		if blockInList.Header.Hash == hash {
			return true
		}
	}
	return false
}

// This function takes a block as the parameter, and returns its parent block.
func (bc *BlockChain) GetParentBlock(block Block) (Block, error) {
	heightToBlocks := bc.Chain
	height := block.Header.Height
	previousList := heightToBlocks[height-1]
	parentHash := block.Header.ParentHash

	for _, b := range previousList {
		if b.Header.Hash == parentHash {
			return b, nil
		}
	}

	return Block{}, errors.New("Parent block not found.")
}

// Return JSON strings from the BlockChain.
func (bc *BlockChain) EncodeToJson() (string, error) {
	blockJsonList := []BlockJson{}

	heightToBlocks := bc.Chain
	for _, list := range heightToBlocks {
		for _, b := range list {
			bj := BuildBjFromBlock(b)
			blockJsonList = append(blockJsonList, bj)
		}
	}

	jsonStrings, err := json.Marshal(blockJsonList)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}

	return string(jsonStrings), nil
}

// Take JSON strings and build BlockChain.
func (bc *BlockChain) DecodeFromJson(jsonStrings string) error {
	bytes := []byte(jsonStrings)
	var bjs []BlockJson
	err := json.Unmarshal(bytes, &bjs)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	for _, bj := range bjs {
		block := BuildBlockFromBj(bj)
		bc.Insert(block)
	}

	return nil
}

// Show the hash of the BlockChain.
func (bc *BlockChain) Show() string {
	rs := ""
	var idList []int
	for id := range bc.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range bc.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is a BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}
