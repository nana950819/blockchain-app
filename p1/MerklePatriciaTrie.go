/*
	Merkle Patricia Trie to be used in Project 3 and 4.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package p1

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Hold encoded prefix in ASCII and value,
// which is actual value in Leaf node or hash value of next node
// in Ext node.
type Flag_value struct {
	// ASCII value array.
	EncodedPrefix []uint8
	// If the node is Ext, 'value' is hash of the next node.
	// If the node is Leaf, 'value' is the string value inserted.
	Value string
}

// Hold node type, which is 0: Null, 1: Branch, 2: Ext or Leaf.
// If the node is not Branch, 'BranchValue' is default.
// If the node is Branch, 'FlagValue' is default.
type Node struct {
	// 0: Null, 1: Branch, 2: Ext or Leaf
	NodeType    int
	BranchValue [17]string
	FlagValue   Flag_value
}

// Hold root node of MerklePatriciaTrie and database for all nodes.
type MerklePatriciaTrie struct {
	// K: Node's hash value, V: Node
	Db map[string]Node
	// hash value of the root node
	Root    string
	ValueDb map[string]string
}

// // Create a new MerklePatriciaTrie. This is used in Project 1.
// func NewMPT() *MerklePatriciaTrie {
// 	// Initialize node.
// 	nullNode := createNewLeafOrExtNode(0, nil, "")

// 	mpt := &MerklePatriciaTrie{}
// 	mpt.Db = map[string]Node{}
// 	mpt.Root = putNodeInDb(nullNode, mpt.Db)

// 	return mpt
// }

// Convert key string to hex value array and append 16.
// If the key is "", then key_hex is 16.
func convert_string_to_hex(key string) []uint8 {
	length := 2*len(key) + 1
	key_hex := make([]uint8, length)
	for i, r := range key {
		key_hex[i*2] = uint8(r / 16)
		key_hex[i*2+1] = uint8(r % 16)
	}
	key_hex[length-1] = 16
	return key_hex
}

// Encode hex array, which ends with 16 if the node is Leaf,
// to ASCII values.
func compact_encode(hex_array []uint8) []uint8 {
	// TODO
	if len(hex_array) == 0 {
		return []uint8{}
	}
	var isLeaf uint8 = 0
	if hex_array[len(hex_array)-1] == 16 {
		isLeaf = 1
	}
	if isLeaf == 1 {
		hex_array = hex_array[:len(hex_array)-1]
	}
	var isOdd uint8 = uint8(len(hex_array) % 2)
	var flagInHexArray uint8 = 2*isLeaf + isOdd
	if isOdd == 1 {
		hex_array = append([]uint8{flagInHexArray}, hex_array...)
	} else {
		hex_array = append(append([]uint8{flagInHexArray}, 0), hex_array...)
	}
	// 'hex_array' now has an even length whose first nibble is the 'flagInHexArray'.
	length := len(hex_array) / 2
	encoded_prefix := make([]uint8, length)
	p := 0
	for i := 0; i < len(hex_array); i += 2 {
		encoded_prefix[p] = 16*hex_array[i] + hex_array[i+1]
		p++
	}

	return encoded_prefix
}

// Decode ASCII array to hex array without 16.
// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {
	// TODO
	if len(encoded_arr) == 0 {
		return []uint8{}
	}

	length := len(encoded_arr) * 2
	hex_array := make([]uint8, length)
	for i, ascii := range encoded_arr {
		hex_array[i*2] = ascii / 16
		hex_array[i*2+1] = ascii % 16
	}

	// Remove prefix and return hex array without 16.
	// If hex hex_array[0] is even, then cut first two. If hex hex_array[0] is odd, then cut first one.
	cut := 2 - hex_array[0]&1
	return hex_array[cut:]
}

// Return two values for testing.
// Get value with key in the MerklePatriciaTrie.
func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	// TODO
	if key == "" {
		return "", nil
	}

	keySearch := convert_string_to_hex(key)

	rootNode := mpt.Db[mpt.Root]
	encodedPrefix := rootNode.FlagValue.EncodedPrefix
	keyMPT := compact_decode(encodedPrefix)

	return get_helper(rootNode, keyMPT, keySearch, mpt.Db), nil
}

// // Get value with key in the MerklePatriciaTrie.
// func (mpt *MerklePatriciaTrie) Get(key string) string {
// 	// TODO
// 	if key == "" {
// 		return ""
// 	}

// 	keySearch := convert_string_to_hex(key)

// 	rootNode := mpt.Db[mpt.Root]
// 	encodedPrefix := rootNode.FlagValue.EncodedPrefix
// 	keyMPT := compact_decode(encodedPrefix)

// 	return get_helper(rootNode, keyMPT, keySearch, mpt.Db)
// }
func get_helper(node Node, keyMPT, keySearch []uint8, db map[string]Node) string {

	nodeType := node.NodeType
	switch nodeType {
	case 0:
		// Null node

		return ""
	case 1:
		// Branch

		if keySearch[0] == 16 {
			// Case B-1. Insert("a"), Insert("aa"), Get("a"), stack 2. keyMPT: [], keySearch: [16], matchLen: 2
			// Case B-3. Insert("aa"), Insert("a"), Get("aa"), stack 2. keyMPT: [], keySearch: [16], matchLen: 2

			return node.BranchValue[16]
		}

		if nextNode, ok := db[node.BranchValue[keySearch[0]]]; ok {
			// There is the next node in the Branch.
			// 'node' is now the next node.
			// Case B-1. Insert("a"), Insert("aa"), Get("aa"), stack 2. keyMPT: [], keySearch: [6 1 16]
			// Case B-2. Insert("a"), Insert("b"), Get("a"), stack 2. keyMPT: [], keySearch: [1 16]
			// Case B-3. Insert("aa"), Insert("a"), Get("aa"), stack 2. keyMPT: [], keySearch: [6 1 16]
			// Case C. Insert("a"), Insert("p"), Get("a"), stack 1. keyMPT: [], keySearch: [6 1 16]
			// Case D-1. Insert("a"), Insert("p"), Insert("abc"), Get("abc") stack 1. keyMPT: [], keySearch: [6 1 6 2 6 3 16]
			// Case D-1. Insert("a"), Insert("p"), Insert("abc"), Get("abc") stack 3. keyMPT: [], keySearch: [6 2 6 3 16]
			// Case D-3. Insert("a"), Insert("p"), Insert("A"), Get("A") stack 1. keyMPT: [], keySearch: [4 1 16]
			encodedPrefix := nextNode.FlagValue.EncodedPrefix
			keyMPT := compact_decode(encodedPrefix)
			return get_helper(nextNode, keyMPT, keySearch[1:], db)
		}

		// There is no link in the Branch.
		return ""

	case 2:
		// Ext or Leaf

		matchLen := prefixLen(keySearch, keyMPT)

		if matchLen != 0 {

			if len(keySearch) <= len(keyMPT) {
				// keySearch is shorter than keyMPT.
				return ""
			}

			if firstDigit := getFirstDigitOfAscii(node.FlagValue.EncodedPrefix); firstDigit == 0 || firstDigit == 1 || firstDigit == 2 {
				// 'node' is Ext node.
				// Case B-1. Insert("a"), Insert("aa"), Get("aa"), stack 1. keyMPT: [6 1], keySearch: [6 1 6 1 16], matchLen: 2
				// Case B-1. Insert("a"), Insert("aa"), Get("a"), stack 1. keyMPT: [6 1], keySearch: [6 1 16], matchLen: 2
				// Case B-2. Insert("a"), Insert("b"), Get("a"), stack 1. keyMPT: [6], keySearch: [6 1 16], matchLen: 1
				// Case B-3. Insert("aa"), Insert("a"), Get("aa"), stack 1. keyMPT: [6 1], keySearch: [6 1 6 1 16], matchLen: 2
				// Case D-1. Insert("a"), Insert("p"), Insert("abc"), Get("abc") stack 2. keyMPT: [1], keySearch: [1 6 2 6 3 16], matchLen: 1

				node = db[node.FlagValue.Value]
				// 'node' is now Branch node next to the Ext node.
				return get_helper(node, keyMPT[matchLen:], keySearch[matchLen:], db)
			}

			// 'node' is Leaf node.

			if keySearch[matchLen] == 16 && len(keyMPT) == matchLen {
				// Exact match.
				// Case A. Insert("a"), Get("a"). keyMPT: [6 1], keySearch: [6 1 16], matchLen: 2
				// Case B-1. Insert("a"), Insert("aa"), Get("aa"), stack 3. keyMPT: [1], keySearch: [1 16], matchLen: 1
				// Case B-3. Insert("aa"), Insert("a"), Get("aa"), stack 3. keyMPT: [1], keySearch: [1 16], matchLen: 1
				// Case C.   Insert("a"), Insert("p"), Get("a"), stack 2. keyMPT: [1], keySearch: [1 16], matchLen: 1
				// Case D-1. Insert("a"), Insert("p"), Insert("abc"), Get("abc") stack 4. keyMPT: [2 6 3], keySearch: [2 6 3 16], matchLen: 3
				// Case D-3. Insert("a"), Insert("p"), Insert("A"), Get("A") stack 2. keyMPT: [1], keySearch: [1 16]
				return node.FlagValue.Value
			}

			// 'node' is Leaf node and keySearch is shorter or longer than keyMPT.
			return ""

		} else if matchLen == 0 {

			if keySearch[matchLen] == 16 {
				// Case B-2. Insert("a"), Insert("b"), Get("a"), stack 3. keyMPT: [], keySearch: [16], matchLen: 0
				return node.FlagValue.Value
			}
			return ""
		}

	}

	return ""
}

// Insert value with key into the MerklePatriciaTrie.
func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// TODO
	if key == "" {
		return
	}

	mpt.ValueDb[key] = new_value

	keySearch := convert_string_to_hex(key)

	db := mpt.Db
	rootNode := db[mpt.Root]
	encodedPrefix := rootNode.FlagValue.EncodedPrefix
	keyMPT := compact_decode(encodedPrefix)

	newRootNode := insert_helper(rootNode, keyMPT, keySearch, new_value, db)
	delete(db, mpt.Root)
	mpt.Root = putNodeInDb(newRootNode, db)

	return
}
func insert_helper(node Node, keyMPT, keySearch []uint8, new_value string, db map[string]Node) Node {

	nodeType := node.NodeType
	switch nodeType {
	case 0:
		// Null node

		// Create a new Leaf node.
		node = createNewLeafOrExtNode(2, keySearch, new_value)
		return node
	case 1:
		// Branch

		if keySearch[0] == 16 {
			// Case E-2. stack 3. keyMPT: [], keySearch: [16]
			node.BranchValue[16] = new_value
			return node
		}

		if nextNode, ok := db[node.BranchValue[keySearch[0]]]; ok {
			// There is a link in the Branch.
			// 'nextNode' is the node next to the Branch. It could be Leaf, Ext, or Branch.
			// Case D-1. Insert("a"), Insert("p"), Insert("abc"), Get("abc"),
			// stack 1. keyMPT: [], keySearch: [6 1 6 2 6 3 16]
			// Case E-1. Insert("p"), Insert("aaaaa"), Insert("aaaap"), Insert("aa"), Get("aa"),
			// stack 1. keyMPT: [], keySearch: [6 1 6 1 16]
			// Case E-2. Insert("p"), Insert("aaaaa"), Insert("aaaap"), Insert("aaaa"), Get("aaaa"),
			// stack 1. keyMPT: [], keySearch: [6 1 6 1 6 1 6 1 16]
			// Case C-2. stack 2. C-3. stack 2. E-4b. stack 1.
			// Case D-2. stack 1.
			encodedPrefix := nextNode.FlagValue.EncodedPrefix
			keyMPT := compact_decode(encodedPrefix)

			nextNode = insert_helper(nextNode, keyMPT, keySearch[1:], new_value, db)
			delete(db, node.hash_node())
			node.BranchValue[keySearch[0]] = putNodeInDb(nextNode, db)
			return node
		}

		// There is no link in the Branch.
		// Case D-3. Insert("a"), Insert("p"), Insert("A"), Get("A"). keyMPT: [], keySearch: [4 1 16]
		leafNode := createNewLeafOrExtNode(2, keySearch[1:], new_value)
		delete(db, node.hash_node())
		node.BranchValue[keySearch[0]] = putNodeInDb(leafNode, db)
		return node

	case 2:
		// Ext or Leaf

		matchLen := prefixLen(keySearch, keyMPT)

		if matchLen != 0 {

			if firstDigit := getFirstDigitOfAscii(node.FlagValue.EncodedPrefix); firstDigit == 0 || firstDigit == 1 || firstDigit == 2 {
				// 'node' is Ext.

				if keySearch[matchLen] == 16 && len(keyMPT) > matchLen {
					// Partial match. keySearch is done and keyMPT is left.
					// Case E-1. stack 2. keyMPT: [1 6 1 6 1 6 1], keySearch: [1 6 1 16], matchLen: 3
					// E-1b. stack 2. keyMPT: [1 6 1 6 1], keySearch: [1 16], matchLen: 1
					extNode1 := createNewLeafOrExtNode(2, keyMPT[:matchLen], "")
					branchNode := Node{NodeType: 1, BranchValue: [17]string{}}
					extNode2 := createNewLeafOrExtNode(2, keyMPT[matchLen+1:], node.FlagValue.Value)
					delete(db, node.hash_node())
					branchNode.BranchValue[keyMPT[matchLen]] = putNodeInDb(extNode2, db)
					branchNode.BranchValue[16] = new_value
					extNode1.FlagValue.Value = putNodeInDb(branchNode, db)

					return extNode1
				}

				if len(keyMPT) == matchLen {
					// Partial match. keyMPT is done.
					// Case E-2. stack 2. keyMPT: [1 6 1 6 1 6 1], keySearch: [1 6 1 6 1 6 1 16], matchLen: 7
					// Case E-3. stack 2. keyMPT: [1 6 1 6 1 6 1], keySearch: [1 6 1 6 1 6 1 4 1 16], matchLen: 7
					// Case C-2. stack 1. C-3.
					branchNode := db[node.FlagValue.Value]
					branchNode = insert_helper(branchNode, nil, keySearch[matchLen:], new_value, db)
					delete(db, node.hash_node())
					node.FlagValue.Value = putNodeInDb(branchNode, db)

					return node
				}

			}

			// 'node' is Leaf.

			if keySearch[matchLen] == 16 && len(keyMPT) == matchLen {
				// Case A (Exact match). keyMPT: [6 1], keySearch: [6 1 16], matchLen: 2
				node.FlagValue.Value = new_value
				return node
			}

			// Case B-1 (Prefix match).
			// stack 1. keyMPT: [6 1], keySearch: [6 1 6 1 16], matchLen: 2
			// Case B-2 (Prefix match).
			// stack 1. keyMPT: [6 1], keySearch: [6 2 16], matchLen: 1
			// Case B-3 (Prefix match).
			// stack 1. keyMPT: [6 1 6 1] , keySearch: [6 1 16], matchLen: 2
			// Case D-1. stack 2. keyMPT: [1], keySearch: [1 6 2 6 3 16], matchLen: 1
			extNode := createNewLeafOrExtNode(2, keyMPT[:matchLen], node.FlagValue.Value)

			branchNode := insert_helper(extNode, keyMPT[matchLen:], keySearch[matchLen:], new_value, db)
			extNode.FlagValue.Value = putNodeInDb(branchNode, db)
			return extNode

		} else if matchLen == 0 {

			if keySearch[matchLen] == 16 && len(keyMPT) == 0 {
				// C-2. stack 3.
				node.FlagValue.Value = new_value
				delete(db, node.hash_node())
				return node
			}

			branchNode := Node{NodeType: 1, BranchValue: [17]string{}}
			if len(keyMPT) == 0 {
				// Case B-1 (Prefix match).
				// stack 2. keyMPT: [], keySearch: [6 1 16], matchLen: 0
				// Case D-1.
				// stack 3. keyMPT: [], keySearch: [6 2 6 3 16], matchLen: 0
				// D-3. stack 3.
				leafNode := createNewLeafOrExtNode(2, keySearch[matchLen+1:], new_value)
				branchNode.BranchValue[keySearch[matchLen]] = putNodeInDb(leafNode, db)
				branchNode.BranchValue[16] = node.FlagValue.Value

			} else if keySearch[matchLen] == 16 {
				// Case B-3 (Prefix match).
				// stack 2. keyMPT: [6 1], keySearch: [16], matchLen: 0
				leafNode := createNewLeafOrExtNode(2, append(keyMPT[matchLen+1:], 16), node.FlagValue.Value)
				branchNode.BranchValue[keyMPT[matchLen]] = putNodeInDb(leafNode, db)
				branchNode.BranchValue[16] = new_value

			} else {
				// Case B-2 (Prefix match).
				// stack 2. keyMPT: [1], keySearch: [2 16], matchLen: 0
				// Case C (Mismatch).
				// keyMPT: [6 1], keySearch: [7 0 16], matchLen: 0
				// D-2. stack 2.
				leafNode := createNewLeafOrExtNode(2, keySearch[matchLen+1:], new_value)
				delete(db, node.hash_node())
				branchNode.BranchValue[keySearch[matchLen]] = putNodeInDb(leafNode, db)

				if _, ok := db[node.FlagValue.Value]; ok {
					if len(keyMPT[1:]) != 0 {
						// E-4.
						extNode := createNewLeafOrExtNode(2, keyMPT[1:], node.FlagValue.Value)
						branchNode.BranchValue[keyMPT[0]] = putNodeInDb(extNode, db)
					} else {
						// E-4b
						branchNode.BranchValue[keyMPT[0]] = node.FlagValue.Value
					}

				} else {
					// B-2, C. D-2.
					leafNode = createNewLeafOrExtNode(2, append(keyMPT[matchLen+1:], 16), node.FlagValue.Value)
					branchNode.BranchValue[keyMPT[matchLen]] = putNodeInDb(leafNode, db)
				}
			}

			return branchNode
		}
	default:

	}

	return node
}

// Delete value in the MerklePatriciaTrie.
func (mpt *MerklePatriciaTrie) Delete(key string) string {
	// TODO
	if key == "" {
		return ""
	}

	keySearch := convert_string_to_hex(key)

	db := mpt.Db
	rootNode := db[mpt.Root]
	encodedPrefix := rootNode.FlagValue.EncodedPrefix
	keyMPT := compact_decode(encodedPrefix)

	newRootNode, ret := delete_helper(rootNode, keyMPT, keySearch, db)
	delete(db, mpt.Root)
	mpt.Root = putNodeInDb(newRootNode, db)

	return ret
}
func delete_helper(node Node, keyMPT, keySearch []uint8, db map[string]Node) (Node, string) {

	nodeType := node.NodeType
	switch nodeType {
	case 0:
		// 'node' is Null node.

		return node, ""
	case 1:
		// 'node' is Branch.

		if keySearch[0] == 16 {
			// Del-4. B-3. stack 2. Insert("aa"), Insert("a"), Delete("a"), stack 2. keyMPT: [], keySearch: [16], matchLen: 2
			// Del-5. stack 2.
			// Del-7.
			node.BranchValue[16] = ""

			if b, oneValue, index := getOnlyOneValueInBranch(node); b {
				// Only one value in the Branch. Rebalance.
				// The value is a link to the next node.
				// Del-4. stack 2.
				// Del-7.

				leftNode := db[oneValue]
				// 'leftNode' could be Leaf, Ext, or Branch.
				if leftNode.NodeType == 1 {
					// 'leftNode' is Branch.
					// Del-7.
					delete(db, node.hash_node())
					extNode := createNewLeafOrExtNode(2, []uint8{index}, oneValue)
					return extNode, ""

				} else if firstDigit := getFirstDigitOfAscii(leftNode.FlagValue.EncodedPrefix); firstDigit == 3 || firstDigit == 4 || firstDigit == 5 {
					// leftNode is Leaf.
					leftNode.FlagValue.EncodedPrefix = compact_encode(
						append([]uint8{index},
							append(compact_decode(leftNode.FlagValue.EncodedPrefix), 16)...))
				} else {
					// leftNode is Ext or Branch.
					// Del-5.
					leftNode.FlagValue.EncodedPrefix = compact_encode(
						append([]uint8{index}, compact_decode(leftNode.FlagValue.EncodedPrefix)...))

				}

				return leftNode, ""
			}

			return node, ""
		}

		if nextNode, ok := db[node.BranchValue[keySearch[0]]]; ok {
			// There is the next node in the Branch.
			// 'nextNode' is the node next to the Branch.
			// Del-3. Insert("a"), Insert("aa"), Delete("aa"), stack 2. keyMPT: [], keySearch: [6 1 16]
			// Del-8. stack 1. stack 2.
			encodedPrefix := nextNode.FlagValue.EncodedPrefix
			keyMPT := compact_decode(encodedPrefix)

			retNode, ret := delete_helper(nextNode, keyMPT, keySearch[1:], db)
			if retNode.NodeType == 0 {
				// 'retNode' is Null node.
				node.BranchValue[keySearch[0]] = ""

				if b, oneValue, index := getOnlyOneValueInBranch(node); b {
					// Only one value in the Branch. Rebalance.
					if node.BranchValue[16] != "" {
						// The value is in the last 16th elem.
						// Del-3. stack 2.
						leafNode := createNewLeafOrExtNode(2, []uint8{16}, oneValue)
						return leafNode, ret
					}
					// The value is a link to the next node.
					// Del-1. stack 2. Del-8.

					leftNode := db[oneValue]
					// 'leftNode' could be Leaf, Ext, or Branch.
					if leftNode.NodeType == 1 {
						// 'leftNode' is Branch.
						// Del-6. Del-8.
						delete(db, node.hash_node())
						extNode := createNewLeafOrExtNode(2, []uint8{index}, oneValue)
						return extNode, ret

					} else if firstDigit := getFirstDigitOfAscii(leftNode.FlagValue.EncodedPrefix); firstDigit == 3 || firstDigit == 4 || firstDigit == 5 {
						// leftNode is Leaf.
						leftNode.FlagValue.EncodedPrefix = compact_encode(
							append([]uint8{index},
								append(compact_decode(leftNode.FlagValue.EncodedPrefix), 16)...))
					} else {
						// leftNode is Ext.
						// Del-2.
						leftNode.FlagValue.EncodedPrefix = compact_encode(
							append([]uint8{index}, compact_decode(leftNode.FlagValue.EncodedPrefix)...))
					}

					return leftNode, ret
				}

				return node, ret
			}

			// Del-8. 'retNode' is Ext.
			node.BranchValue[keySearch[0]] = putNodeInDb(retNode, db)

			return node, ret
		}

		// There is no link in the Branch.
		return node, "path_not_found"

	case 2:
		// 'node' is Ext or Leaf.

		matchLen := prefixLen(keySearch, keyMPT)

		if matchLen != 0 {

			if firstDigit := getFirstDigitOfAscii(node.FlagValue.EncodedPrefix); firstDigit == 0 || firstDigit == 1 || firstDigit == 2 {
				// 'node' is Ext node.
				// Del-3. Insert("a"), Insert("aa"), Delete("aa"), stack 1. keyMPT: [6 1], keySearch: [6 1 6 1 16], matchLen: 2
				// Del-1. Insert("a"), Insert("b"), Delete("b"), stack 1.
				// Del-4. stack 1.
				// Del-2. stack 1.
				branchNode := db[node.FlagValue.Value]
				// 'node' is now Branch node next to the Ext node.

				retNode, ret := delete_helper(branchNode, keyMPT[matchLen:], keySearch[matchLen:], db)
				// 'retNode' could be Leaf, Ext, or Branch.
				if retNode.NodeType == 1 {
					// 'retNode' is Branch.
					delete(db, node.hash_node())
					node.FlagValue.Value = putNodeInDb(retNode, db)
					return node, ret

				} else if firstDigit := getFirstDigitOfAscii(retNode.FlagValue.EncodedPrefix); firstDigit == 3 || firstDigit == 4 || firstDigit == 5 {
					// 'retNode' is Leaf.
					retNode.FlagValue.EncodedPrefix = compact_encode(
						append(keyMPT,
							append(compact_decode(retNode.FlagValue.EncodedPrefix), 16)...))

				} else {
					// 'retNode' is Ext.
					// Del-2. Del-7.
					retNode.FlagValue.EncodedPrefix = compact_encode(
						append(keyMPT, compact_decode(retNode.FlagValue.EncodedPrefix)...))
				}
				delete(db, node.hash_node())
				node.FlagValue.Value = putNodeInDb(retNode, db)
				return retNode, ret
			}

			// 'node' is Leaf node.

			if keySearch[matchLen] == 16 && len(keyMPT) == matchLen {
				// Exact match.
				// Del-0. Just one node.
				// Del-3. Insert("a"), Insert("aa"), Delete("aa"), stack 3. keyMPT: [1], keySearch: [1 16], matchLen: 1
				delete(db, node.hash_node())
				nullNode := createNewLeafOrExtNode(0, nil, "")
				return nullNode, ""
			}

			// 'node' is Leaf node and keySearch is shorter or longer than keyMPT.
			return node, "path_not_found"

		} else if matchLen == 0 {

			if keySearch[matchLen] == 16 {
				// Del-1. Insert("a"), Insert("b"), Delete("b"), stack 3. keyMPT: [], keySearch: [16], matchLen: 0
				delete(db, node.hash_node())
				nullNode := createNewLeafOrExtNode(0, nil, "")
				return nullNode, ""
			}

			return node, "path_not_found"
		}

	}

	return node, "path_not_found"
}

// Check if there is only one value in the Branch node, and
// if yes, then return true, that value, and the index of that value
// in the node.
func getOnlyOneValueInBranch(node Node) (bool, string, uint8) {
	count := 0
	var index uint8 = 0
	oneValue := ""
	for i, str := range node.BranchValue {
		if str != "" {
			index = uint8(i)
			oneValue = str
			count++
		}
	}
	return count <= 1, oneValue, index
}

// Create a new node.
func createNewLeafOrExtNode(nodeType int, keyHex []uint8, newValue string) Node {
	encodedPrefix := compact_encode(keyHex)
	flagValue := Flag_value{EncodedPrefix: encodedPrefix, Value: newValue}
	node := Node{NodeType: nodeType, FlagValue: flagValue}
	return node
}

// Hash the node and put it into the database, and
// return the hash value.
func putNodeInDb(node Node, db map[string]Node) string {
	hash := node.hash_node()
	db[hash] = node
	return hash
}

// Return the length of the same prefix of array a and array b
// if both a and b have prefix in common.
func prefixLen(a []uint8, b []uint8) int {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}
	i := 0
	for i < length {
		if a[i] != b[i] {
			break
		}
		i++
	}
	return i
}

// Get the first digit of ASCII value.
// ex. if ASCII: 32, then return 3.
func getFirstDigitOfAscii(encodedPrefix []uint8) uint8 {
	firstDigit := encodedPrefix[0] / 10
	return firstDigit
}

// The functions below has been provided by the instructors.

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.NodeType {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.BranchValue {
			str += v
		}
	case 2:
		str = node.FlagValue.Value
	}

	// The instructor said "feel free to change that part of "hash_node()".
	// address := fmt.Sprintf("%p", node)
	// str = address + str

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

// Additional given code.
func (node *Node) String() string {
	str := "empty string"
	switch node.NodeType {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.BranchValue[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.BranchValue[16])
	case 2:
		encoded_prefix := node.FlagValue.EncodedPrefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.FlagValue.Value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.Db = make(map[string]Node)
	mpt.ValueDb = make(map[string]string)
	nullNode := createNewLeafOrExtNode(0, nil, "")
	mpt.Root = putNodeInDb(nullNode, mpt.Db)

	rentalMap := make(map[string]RentalInfo)
	jRentalMap, _ := EncodeRentalMapToJson(rentalMap)
	mpt.Insert("ads", jRentalMap)

	historyMap := make(map[string]History)
	jHistoryMap, _ := EncodeHistoryMapToJson(historyMap)
	mpt.Insert("history", jHistoryMap)

}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.Root)
	for hash := range mpt.Db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.Db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("Hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}
