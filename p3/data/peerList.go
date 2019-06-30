/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	Hold information for PeerList.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"sync"
)

// Hold the addresses and IDs of the peers.
type PeerList struct {
	selfId int32
	// K: IP address and port number, V: ID
	peerMap   map[string]int32
	maxLength int32
	// lock
	mux sync.Mutex
}

// Create a new PeerList object.
func NewPeerList(id int32, maxLength int32) PeerList {
	return PeerList{selfId: id, peerMap: make(map[string]int32), maxLength: maxLength}
}

// Add the address and ID of a peer.
func (peers *PeerList) Add(addr string, id int32) {
	peers.mux.Lock()
	peers.peerMap[addr] = id
	peers.mux.Unlock()
}

// Delete the address and ID of a peer.
func (peers *PeerList) Delete(addr string) {
	peers.mux.Lock()
	delete(peers.peerMap, addr)
	peers.mux.Unlock()
}

// From peerMap, get a list of addresses of maxLength closest peers.
func (peers *PeerList) Rebalance() []string {
	// peers.mux.Lock() // When func_test, comment this out.

	idToAddr := map[int32]string{}
	idList := []int{}
	for k, v := range peers.peerMap {
		idList = append(idList, int(v))
		idToAddr[v] = k
	}
	var closestIds, deleteIds []int
	sizeClosest := int(peers.maxLength)
	if len(idList) > sizeClosest {
		si := int(peers.selfId)
		sort.Ints(idList)
		closestIds, deleteIds = chooseClosest(idList, si, sizeClosest)
	} else {
		closestIds = idList
	}

	// Rebalance the map.
	for _, deleteId := range deleteIds {
		peers.Delete(idToAddr[int32(deleteId)])
	}

	// Get maxLength closest addresses.
	closestAddrs := []string{}
	for _, id := range closestIds {
		closestAddrs = append(closestAddrs, idToAddr[int32(id)])
	}

	// defer peers.mux.Unlock() // When func_test, comment this out.
	return closestAddrs
}

// Choose closest ids from selfId in the list. O(N)
func chooseClosest(list []int, selfId int, sizeClosest int) ([]int, []int) {
	closestIds := []int{}
	deleteIds := []int{}
	var left, right int
	sizeList := len(list)
	for i, v := range list {
		if v > selfId {
			left = i - sizeClosest/2
			right = i + sizeClosest/2
			if right > sizeList {
				right = right - sizeList
				closestIds = append(closestIds, list[:right]...)
				closestIds = append(closestIds, list[left:]...)
				deleteIds = append(deleteIds, list[right:left]...)
			} else if left < 0 {
				left = sizeList + left
				closestIds = append(closestIds, list[:right]...)
				closestIds = append(closestIds, list[left:]...)
				deleteIds = append(deleteIds, list[right:left]...)
			} else {
				closestIds = append(closestIds, list[left:right]...)
				deleteIds = append(deleteIds, list[:left]...)
				deleteIds = append(deleteIds, list[right:]...)
			}

			return closestIds, deleteIds
		}
	}
	left = sizeList - sizeClosest/2
	right = sizeClosest / 2
	closestIds = append(closestIds, list[:right]...)
	closestIds = append(closestIds, list[left:]...)
	deleteIds = append(deleteIds, list[right:left]...)

	return closestIds, deleteIds
}

// Show the content of the PeerList.
func (peers *PeerList) Show() string {
	peers.mux.Lock()

	fmt.Println("PeerList: ")
	var sb strings.Builder
	fmt.Printf("selfId: %d\n", peers.selfId)
	fmt.Println("peerMap: ")
	for k, v := range peers.peerMap {
		sb.WriteString(k)
		sb.WriteString(" ")
		sb.WriteString(strconv.Itoa(int(v)))
		sb.WriteString("\n")
	}

	defer peers.mux.Unlock()
	return sb.String()
}

// Assign the id.
func (peers *PeerList) Register(id int32) {
	peers.mux.Lock()
	peers.selfId = id
	fmt.Printf("SelfId=%v\n", id)
	peers.mux.Unlock()
}

// Copy the peerMap in the PeerList.
func (peers *PeerList) Copy() map[string]int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.peerMap
}

// Get my ID.
func (peers *PeerList) GetSelfId() int32 {
	peers.mux.Lock()
	defer peers.mux.Unlock()
	return peers.selfId
}

// Convert peerMap to Json string.
func (peers *PeerList) PeerMapToJson() (string, error) {
	peers.mux.Lock()
	jsonString, err := json.Marshal(peers.peerMap)
	if err != nil {
		fmt.Println("error:", err)
	}
	defer peers.mux.Unlock()
	return string(jsonString), nil
}

// Convert Json string to peerMap.
func (peers *PeerList) JsonToPeerMap(peerMapJson string) map[string]int32 {
	var peerMap map[string]int32
	err := json.Unmarshal([]byte(peerMapJson), &peerMap)
	if err != nil {
		fmt.Println("error:", err)
	}

	return peerMap
}

// Add another peerMap to my own peerMap.
func (peers *PeerList) InjectPeerMapJson(peerMapJsonStr string, selfAddr string) {
	incomingPeerMap := peers.JsonToPeerMap(peerMapJsonStr)

	for k, v := range incomingPeerMap {
		if k != selfAddr {
			peers.peerMap[k] = v
		}
	}
}

// Test Rebalance().
func TestPeerListRebalance() {
	fmt.Println("test 1: ")
	peers := NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected := NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	expected.Add("-1-1", -1)
	b := reflect.DeepEqual(peers, expected)
	fmt.Println(b)
	if !b {
		fmt.Printf("test 1, Expected %d, \nbut was %d\n", expected, peers)
	}

	fmt.Println("test 2: ")
	peers = NewPeerList(5, 2)
	peers.Add("1111", 1)
	peers.Add("4444", 4)
	peers.Add("-1-1", -1)
	peers.Add("0000", 0)
	peers.Add("2121", 21)
	peers.Rebalance()
	expected = NewPeerList(5, 2)
	expected.Add("4444", 4)
	expected.Add("2121", 21)
	b = reflect.DeepEqual(peers, expected)
	fmt.Println(b)
	if !b {
		fmt.Printf("test 2, Expected %d, \nbut was %d\n", expected, peers)
	}

	fmt.Println("test 3: ")
	peers = NewPeerList(5, 4)
	peers.Add("1111", 1)
	peers.Add("7777", 7)
	peers.Add("9999", 9)
	peers.Add("11111111", 11)
	peers.Add("2020", 20)
	peers.Rebalance()
	expected = NewPeerList(5, 4)
	expected.Add("1111", 1)
	expected.Add("7777", 7)
	expected.Add("9999", 9)
	expected.Add("2020", 20)
	b = reflect.DeepEqual(peers, expected)
	fmt.Println(b)
	if !b {
		fmt.Printf("test 3, Expected %d, \nbut was %d\n", expected, peers)
	}

}
