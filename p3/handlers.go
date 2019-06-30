/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	This is a client and a server in the node.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package p3

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cs686/cs686-blockchain-p3-kayfuku/p1"
	"cs686/cs686-blockchain-p3-kayfuku/p2"
	"cs686/cs686-blockchain-p3-kayfuku/p3/data"

	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/sha3"
)

// Create a tunnel to the TA's server "mc07.cs.usfca.edu:6688" from "localhost:6688".
// var TA_SERVER = "http://localhost:6688"
var TA_SERVER = "localhost:6688"
var FIRST_NODE_ADDR = "localhost:6686"
var SECOND_NODE_ADDR = "localhost:6687"

// (Changed?)
// var REGISTER_SERVER = TA_SERVER + "/peer"
// var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var SELF_ADDR string
var SELF_PORT string
var SELF_ID int32
var SBC data.SyncBlockChain
var Peers data.PeerList
var ifStarted bool
var RD data.RegisterData
var isMiner bool = true // No need?

var HEART_BEAT_INTERVAL float32 = 5

// var HEART_BEAT_INTERVAL float32 = 3

// var NONCE_HASH_PREFIX string = "000000"

var NONCE_HASH_PREFIX string = "00000"
var currTxIdForPoW string = ""
var stopPoW bool = false

var rentalIdList mapset.Set = mapset.NewSet()
var queue chan data.Transaction

// Initialization.
// This function will be executed before everything else.
func init() {
	fmt.Println("init() start.")
	fmt.Println("HELLO!")

	// Initialize PeerList.
	Peers = data.NewPeerList(0, 32)
}

// Register ID, download BlockChain, start HeartBeat.
func Start(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Miner Node " + SELF_PORT + " launched.")
	SELF_ADDR = "localhost:" + SELF_PORT

	// Get an ID.
	Register()

	// Get a BlockChain.
	fmt.Println("FIRST_NODE_ADDR: ", FIRST_NODE_ADDR)
	fmt.Println("SELF_ADDR: ", SELF_ADDR)
	if SELF_ADDR == FIRST_NODE_ADDR {
		fmt.Println("Initiate the blockchain.")
		// Initiate the blockchain.
		SBC = data.NewBlockChain()
	} else {
		Download()
	}

	// Start Heart Beat on one thread.
	go StartHeartBeat()

	// Start finding Nonce on another thread. For project 4.
	// go StartTryingNonces()
	go StartDoingPoW()

}

// (User) Register ID, download BlockChain, start HeartBeat.
func StartUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User Node " + SELF_PORT + " launched.")

	SELF_ADDR = "localhost:" + SELF_PORT
	isMiner = false

	// Get an ID.
	Register()

	// Get a BlockChain.
	fmt.Println("FIRST_NODE_ADDR: ", FIRST_NODE_ADDR)
	fmt.Println("SELF_ADDR: ", SELF_ADDR)
	if SELF_ADDR == FIRST_NODE_ADDR {
		fmt.Println("Initiate the blockchain.")
		// Initiate the blockchain.
		SBC = data.NewBlockChain()
	} else {
		Download()
	}

	// Start Heart Beat on one thread.
	go StartHeartBeat()

}

// For testing, Server, doPost().
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\n")
}

// Display peerList and sbc.
func Show(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
	// fmt.Println(Peers.Show())
	fmt.Println(SBC.Show())
	// fmt.Println(SBC.BlockChainToJson())
}

// Register to TA's server (Changed?), get an ID (nodeId).
func Register() {
	id, _ := strconv.Atoi(SELF_PORT)
	// Since ID server is down, this is commented out temporarily.
	// id, _ := getId()
	fmt.Println("id:", id)
	SELF_ID = int32(id)
	RD = data.NewRegisterData(int32(id), "")
	RD.PortNum = SELF_PORT
	Peers.Register(int32(id))
}

// (Client) HTTP Get Request to get an ID from TA's server.
func getId() (int32, error) {
	url := fmt.Sprintf(TA_SERVER + "/peer")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, err
	}
	body, err := doRequest(req)
	if err != nil {
		return -1, err
	}
	id, _ := strconv.Atoi(string(body))

	return int32(id), nil
}

// For fetching id, do HTTP Request and get Response Body.
func doRequest(req *http.Request) ([]byte, error) {
	// fmt.Println("doRequest() start.")

	client := &http.Client{}

	// Do
	// fmt.Println("Do() start.")
	// fmt.Println("req:", req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// HTTP Response.
	// Header.
	// fmt.Println("--- HTTP Response ----")
	// fmt.Printf("[status] %d\n", resp.StatusCode)
	// for k, v := range resp.Header {
	// 	fmt.Print("[header] " + k)
	// 	fmt.Println(": " + strings.Join(v, ","))
	// }

	// Body.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Response body] " + string(body))

	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	// fmt.Println("doRequest() end.")
	return body, nil
}

// (Client) HTTP POST Request with my id, and
// download blockchain from the first node.
func Download() error {
	fmt.Println("Download() start.")
	apiUrl := fmt.Sprintf("http://" + FIRST_NODE_ADDR + "/upload")
	fmt.Println("apiUrl: ", apiUrl)

	// Make data.
	values, err := json.Marshal(RD)
	// selfId := RD.AssignedId
	// data := url.Values{}
	// data.Add("senderId", strconv.Itoa(int(selfId)))
	// data.Encode()

	// Create a request.
	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(values))
	// req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	// Send POST Request with id, and get blockChainJson.
	// res, err := http.PostForm(apiUrl, data)
	body, err := doRequest(req)
	if err != nil {
		return err
	}

	blockChainJson := string(body)
	fmt.Println("UpdateEntireBlockChain() start.")
	fmt.Println("blockChainJson:", blockChainJson)
	if SELF_ADDR == SECOND_NODE_ADDR {
		SBC = data.NewBlockChain()
	} else {
		SBC = data.NewBlockChain()
		SBC.UpdateEntireBlockChain(blockChainJson)
	}

	fmt.Println("Download() end.")
	return err
}

// (Server) doPost() to get the caller's address and id.
// Upload blockchain to whoever called this method, return BlockChain's jsonStr. (/upload)
func Upload(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Upload() start.")

	// HTTP Request.
	// Header.
	fmt.Println("--- HTTP Request ----")
	method := req.Method
	fmt.Println("[method] " + method)
	for k, v := range req.Header {
		fmt.Print("[header] " + k)
		fmt.Println(": " + strings.Join(v, ","))
	}

	// POST
	if method == "POST" {
		defer req.Body.Close()
		// Get addr and id from client.
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("[request body] " + string(body))

		// Unmarshal
		var rd data.RegisterData
		err = json.Unmarshal(body, &rd)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		// Add addr and id to PeerList.
		// clientId, _ := strconv.Atoi(string(body))
		clientId := rd.AssignedId
		clientAddr := "localhost:" + rd.PortNum
		fmt.Println("clientId in Uplopd():", clientId)
		fmt.Println("clientAddr in Uplopd():", clientAddr)
		fmt.Println("Put addr and id in PeerList.")
		Peers.Add(clientAddr, int32(clientId))

		// Send blockChainJson to client.
		blockChainJson, err := SBC.BlockChainToJson()
		if err != nil {
			log.Fatal(err)
			// data.PrintError(err, "Upload")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// if err := json.NewEncoder(w).Encode(blockChainJson); err != nil {
		// 	panic(err)
		// }
		fmt.Fprint(w, blockChainJson)
		// fmt.Fprint(w, "Recieved Post(json) request!!")

	}
}

// (Server) doGet() to send the specified block.
// Upload a block to whoever called this method, return jsonStr. (/block/{height}/{hash})
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UploadBlock() start.")

	// HTTP Request.
	// Header.
	// fmt.Println("--- HTTP Request ----")
	method := r.Method
	// fmt.Println("[method] " + method)
	// for k, v := range r.Header {
	// 	fmt.Print("[header] " + k)
	// 	fmt.Println(": " + strings.Join(v, ","))
	// }

	// GET
	if method == "GET" {
		vars := mux.Vars(r)
		height, _ := strconv.Atoi(vars["height"])
		hash := vars["hash"]
		fmt.Println("height:", height)
		fmt.Println("hash:", hash)

		fmt.Println("GetBlock() start.")
		block, _ := SBC.GetBlock(int32(height), hash)
		// fmt.Println("EncodeToJson() start.")
		jsonBlockStr, err := p2.EncodeToJson(block)
		fmt.Println("Block ready to be uploaded.")
		fmt.Println("Uploading block:", jsonBlockStr)

		if err != nil {
			log.Fatal(err)
			// data.PrintError(err, "Upload")
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, jsonBlockStr)
	}
}

// Receive a heartbeat. (/heartbeat/receive)
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	fmt.Println(SELF_ADDR + " HeartBeatReceive() start.")
	// Get a HBD from a peer.
	hbd := getHBD(w, r)

	// Add info to PeerList.
	if hbd.Addr != SELF_ADDR {
		Peers.Add(hbd.Addr, hbd.Id)
	}
	Peers.InjectPeerMapJson(hbd.PeerMapJson, SELF_ADDR)

	// If the HBD has a new block,
	// verify nonce and insert the new block to the current BlockChain.
	if hbd.IfNewBlock {
		fmt.Println(SELF_ADDR + " received the new block!")
		jsonBlockStr := hbd.BlockJson
		fmt.Println("The new block:", jsonBlockStr)
		block := p2.DecodeFromJson(jsonBlockStr)

		// Verify the new block.
		isValid := verifyBlock(block)
		if isValid {
			parentHash := block.Header.ParentHash

			// Stop calculating nonce if TX ID matches.
			if hbd.TxId == currTxIdForPoW {
				fmt.Println("Someone found a nonce first for the same transaction!")
				stopPoW = true
				fmt.Println("Changed stopPoW to true.")
			}

			// Get all the ancestors.
			if !SBC.CheckParentHash(block) {
				// There is no parent block.
				fmt.Println("No parent!")
				AskForBlock(block.Header.Height-1, parentHash)
			} else {
				fmt.Println("Parent block exists.")
			}

			// Update MPT.
			blockMptUpdated := SBC.UpdateMPT(block)

			// Insert the new block.
			b := SBC.Insert(blockMptUpdated)
			if b {
				fmt.Println("New block has been inserted!")
				fmt.Println("New block:", jsonBlockStr)
			} else {
				fmt.Println("New block has not been inserted.")
			}

		} else {
			// The new block is invalid. Ignore it.
			fmt.Println("New block is invalid!")
		}
	}

	// TODO
	if hbd.IfNewTx && hbd.TxId != currTxIdForPoW {
		fmt.Println(SELF_ADDR + " received the new TX!")
		jTx := hbd.TxJson
		fmt.Println("The TX received:", jTx)
		tx := data.DecodeFromJson(jTx)

		// TODO: Pick up the one TX with max TX fee.

		// Miner can get TX fee.
		fmt.Printf("Miner earned Tx fee %d!\n", tx.TxFee)

		// TODO: Verify Rental ID.

		// TODO: Verify payment.

		// TODO: Verify signature.

		// Put the TX in the pool.
		fmt.Println("Put the TX into the queue.")
		queue <- tx

	}

	// Forward HBD.
	hbd.Hops--
	if hbd.Hops > 0 {
		ForwardHeartBeat(hbd)
	}

}

// (User) Receive a heartbeat. (/heartbeat/receive)
func HeartBeatReceiveForUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(SELF_ADDR + " HeartBeatReceiveForUser() start.")
	// Get a HBD from a peer.
	hbd := getHBD(w, r)

	// Add info to PeerList.
	if hbd.Addr != SELF_ADDR {
		Peers.Add(hbd.Addr, hbd.Id)
	}
	Peers.InjectPeerMapJson(hbd.PeerMapJson, SELF_ADDR)

	// If the HBD has a new block,
	// verify nonce and insert the new block to the current BlockChain.
	if hbd.IfNewBlock {
		fmt.Println(SELF_ADDR + " received the new block!")
		jsonBlockStr := hbd.BlockJson
		fmt.Println("The new block:", jsonBlockStr)
		block := p2.DecodeFromJson(jsonBlockStr)

		// Verify the new block.
		isValid := verifyBlock(block)
		if isValid {
			parentHash := block.Header.ParentHash

			// Get all the ancestors.
			if !SBC.CheckParentHash(block) {
				// There is no parent block.
				fmt.Println("No parent!")
				AskForBlock(block.Header.Height-1, parentHash)
			} else {
				fmt.Println("Parent block exists.")
			}

			// Update MPT.
			blockMptUpdated := SBC.UpdateMPT(block)

			// Insert the new block.
			b := SBC.Insert(blockMptUpdated)
			if b {
				fmt.Println("New block has been inserted!")
				fmt.Println("New block:", jsonBlockStr)
			} else {
				fmt.Println("New block has not been inserted.")
			}

			// Check Rental ID.
			jHistoryMap := block.Value.ValueDb["history"]
			historyMap, _ := p1.DecodeJsonToHistoryMap(jHistoryMap)
			fmt.Println("Check historyMap:", historyMap)
			for _, v := range historyMap {
				if v.State == 0 {
					// Store my own rentalId. (After Publish())
					jRentalMap := block.Value.ValueDb["ads"]
					rentalMap, _ := p1.DecodeJsonToRentalMap(jRentalMap)
					fmt.Println("Check rentalMap:", rentalMap)
					for k, v := range rentalMap {
						if v.LenderId == SELF_ID {
							rentalIdList.Add(k)
							fmt.Println("My own rentalId has been stored in the list.")
						}
					}
					break
				}
				if v.State == 1 {
					// Start Request has been confirmed. (After Start Request submitted)
					fmt.Println("Start Request for my own rental ID!")
					Permit(v)
					break
				}
				if v.State == 2 {
					// Start Time has been confirmed. (After Permit())
					fmt.Println("Start Time has been confirmed!")
					fmt.Println("Start Time:", v.StartTime)
					fmt.Println("Lender " + SELF_ADDR + " sent the KEY to use the asset to the Borrower " + strconv.Itoa(int(v.BorrowerId)) + "'s smartphone.")
					break
				}
				if v.State == 3 {
					// Stop Request has been confirmed. (After Stop Request submitted)
					fmt.Println("Stop Request for my own rental ID!")
					Checkout(v)
					break
				}
				if v.State == 4 {
					// End Time has been confirmed. (After Checkout())
					fmt.Println("End Time has been confirmed!")
					fmt.Println("End Time:", v.EndTime)
					fmt.Println("HELLO!!")
					fmt.Println("Lender " + SELF_ADDR + " sent the REMAINING MONEY " + strconv.Itoa(int(v.Deposit)) + " back to the Borrower " + strconv.Itoa(int(v.BorrowerId)) + "'s smartphone.")
					InitializeHistory(v.RentalId)
					break
				}

			} // end of for loop

		} else {
			// The new block is invalid. Ignore it.
			fmt.Println("New block is invalid!")
		}
	}

	// Forward HBD.
	hbd.Hops--
	if hbd.Hops > 0 {
		ForwardHeartBeat(hbd)
	}
}

// (Server) doPost() to get a HBD.
func getHBD(w http.ResponseWriter, r *http.Request) data.HeartBeatData {
	var hbd data.HeartBeatData

	// HTTP Request.
	// Header.
	// fmt.Println("--- HTTP Request ----")
	method := r.Method
	// fmt.Println("[method] " + method)
	// for k, v := range r.Header {
	// 	fmt.Print("[header] " + k)
	// 	fmt.Println(": " + strings.Join(v, ","))
	// }

	// POST
	if method == "POST" {
		defer r.Body.Close()
		// Get the HeartBeatData.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("[request body](HBD): " + string(body))

		// Unmarshal
		err = json.Unmarshal(body, &hbd)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
	}

	return hbd
}

// Verify the block to see if the nonceHash meets the requirement.
func verifyBlock(block p2.Block) bool {
	isValid := false
	parentHash := block.Header.ParentHash
	nonce := block.Header.Nonce
	mptRoot := block.Value.Root
	fmt.Println("parentHash:", parentHash)
	fmt.Println("nonce:", nonce)
	fmt.Println("mptRoot:", mptRoot)

	nonceHash := calcNonceHash(parentHash, nonce, mptRoot)
	if strings.HasPrefix(nonceHash, NONCE_HASH_PREFIX) {
		isValid = true
		fmt.Println("Verification success!")
		fmt.Println("nonceHash:", nonceHash)
	} else {
		fmt.Println("nonceHash:", nonceHash)
	}

	return isValid
}

// Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string) {
	fmt.Println("AskForBlock() start.")
	peerMap := Peers.Copy()
	for peerAddr, _ := range peerMap {
		b, _ := askPeerForBlock(peerAddr, height, hash)
		if b {
			break
		}
	}
	// We assume that at least one node has the parent block.
}

// (Client) HTTP GET Request to get the block.
func askPeerForBlock(peerAddr string, height int32, hash string) (bool, error) {

	fmt.Println("Start asking parent block.")
	apiUrl := fmt.Sprintf("http://" + peerAddr + "/block/" + strconv.Itoa(int(height)) + "/" + hash)
	fmt.Println("ask " + apiUrl + " for block.")
	fmt.Println("height:", height)
	fmt.Println("hash:", hash)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return false, err
	}

	body, err := doRequest(req)
	if err != nil {
		return false, err
	}
	jsonBlockStr := string(body)

	// TODO
	// Insert the parent block.
	parentBlock := p2.DecodeFromJson(jsonBlockStr)
	b := SBC.Insert(parentBlock)
	// For testing.
	if b {
		fmt.Println("Parent block has been inserted!")
		fmt.Println("Parent block:", jsonBlockStr)
	}

	// Recursively check if the parent exists.
	if !SBC.CheckParentHash(parentBlock) {
		fmt.Println("No parent again!")
		AskForBlock(parentBlock.Header.Height-1, parentBlock.Header.ParentHash)
	}

	return b, nil
}

// Foward the HBD.
func ForwardHeartBeat(heartBeatData data.HeartBeatData) {
	fmt.Println(SELF_ADDR + " ForwardHeartBeat() start.")
	peerMap := Peers.Copy()

	// Send the HeartBeatData.
	for peerAddr, _ := range peerMap {
		sendHeartBeat(peerAddr, heartBeatData)
	}

}

// Start Heart Beat.
func StartHeartBeat() {

	for {
		interval := HEART_BEAT_INTERVAL + rand.Float32()*HEART_BEAT_INTERVAL
		// fmt.Println("interval: ", interval)
		fmt.Println("Heart Beat!")

		// Initialize HeartBeatData
		selfId := RD.AssignedId
		peerMapJson, _ := Peers.PeerMapToJson()
		heartBeatData := data.PrepareHeartBeatData(&SBC, selfId, peerMapJson, SELF_ADDR)

		// Send the HeartBeatData.
		closestAddresses := Peers.Rebalance()
		for _, addr := range closestAddresses {
			sendHeartBeat(addr, heartBeatData)
		}

		// Sleep for a certain amout of time.
		time.Sleep(time.Duration(interval) * time.Second)
	}

	// // Another version of doing something periodically.
	// t := time.NewTicker(hearBeatInterval * time.Second) // notify every heartBeatInterval sec.
	// for {
	// 	select {
	// 	case <-t.C:
	// 		// do something periodically here.
	// 		fmt.Println("Heart Beat ")
	// 	}
	// }
	// t.Stop() // Stop the timer.
}

// (Client) HTTP POST Request. Send the HeartBeatData.
func sendHeartBeat(addr string, hbd data.HeartBeatData) error {
	apiUrl := fmt.Sprintf("http://" + addr + "/heartbeat/receive")
	fmt.Println("Send HBD from " + SELF_ADDR + " to " + apiUrl + ".")
	jHBD, err := json.Marshal(hbd)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jHBD))
	if err != nil {
		return err
	}
	_, err = doRequest(req)

	return err
}

// Take TX, do PoW, and send the block.
// func StartDoingPoW(tx data.Transaction) {
func StartDoingPoW() {
	fmt.Println("StartDoingPoW() start. ")
	fmt.Println("Create a queue.")
	queue = make(chan data.Transaction, 100)
	defer close(queue)

	for tx := range queue {
		fmt.Println("TX taken from the queue.")

		stopPoW = false

		// Create MPT.
		mpt := p1.MerklePatriciaTrie{}
		mpt.Initial()

		// Look into the TX and store RentalInfo and History in MPT.
		// RentalInfo in TX.
		rentalMap := make(map[string]p1.RentalInfo)
		if tx.RentalInfoJson != "" {
			rentalInfo := p1.DecodeJsonToRentalInfo(tx.RentalInfoJson)
			rentalMap[rentalInfo.RentalId] = rentalInfo
		}
		jRentalMap, _ := p1.EncodeRentalMapToJson(rentalMap)
		mpt.Insert("ads", jRentalMap)

		// History in TX.
		historyMap := make(map[string]p1.History)
		if tx.HistoryJson != "" {
			history := p1.DecodeJsonToHistory(tx.HistoryJson)
			historyMap[history.HistoryId] = history
		}
		jHistoryMap, _ := p1.EncodeHistoryMapToJson(historyMap)
		mpt.Insert("history", jHistoryMap)

		// Do Proof of Work.
		listBlock := SBC.GetLatestBlocks()
		parentHash := listBlock[0].Header.Hash
		parentHeight := listBlock[0].Header.Height
		currTxIdForPoW = tx.TxId

		nonce := doProofOfWork(parentHash, mpt.Root)
		if stopPoW {
			fmt.Println("Stop creating block.")
			continue
		}

		// Create a new block.
		timeStamp := time.Now().UnixNano()
		block := p2.NewBlock(parentHeight+1, timeStamp, parentHash, mpt, nonce)
		jsonBlockStr, _ := p2.EncodeToJson(block)

		// Update MPT.
		blockMptUpdated := SBC.UpdateMPT(block)

		// Insert the new block.
		b := SBC.Insert(blockMptUpdated)
		if b {
			fmt.Println("New block has been inserted!")
			fmt.Println("New block:", jsonBlockStr)
		} else {
			fmt.Println("New block has not been inserted.")
		}

		// Create a new HeartBeatData.
		addNewBlock := true
		addNewTx := false
		jTx := ""
		selfId := RD.AssignedId
		peerMapJson, _ := Peers.PeerMapToJson()
		heartBeatData := data.NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, SELF_ADDR)
		heartBeatData.TxId = tx.TxId
		fmt.Println("Added a new block to HBD")
		fmt.Println("block added: ", jsonBlockStr)

		// Send the HeartBeatData to P2P network.
		closestAddresses := Peers.Rebalance()
		for _, addr := range closestAddresses {
			sendHeartBeat(addr, heartBeatData)
		}

	}

}

// // (For project 4) Try to find a nonce to add a new block on another thread.
// func StartTryingNonces() {
// 	fmt.Println("StartTryingNonces() start.")

// 	selfId := RD.AssignedId
// 	peerMapJson, _ := Peers.PeerMapToJson()

// 	for {
// 		fmt.Println("Start trying to add a new block.")
// 		stopPoW = false

// 		// Prepare HeartBeatData.
// 		heartBeatData := prepareHeartBeatDataWithBlock(selfId, peerMapJson, SELF_ADDR)
// 		if stopPoW {
// 			fmt.Println("Stop trying to send HBD.")
// 			continue
// 		}

// 		// Send the HeartBeatData.
// 		closestAddresses := Peers.Rebalance()
// 		for _, addr := range closestAddresses {
// 			sendHeartBeat(addr, heartBeatData)
// 		}
// 	}
// }

// // (For project 4) Add a new block to HBD if a nonce found doing Proof of Work.
// func prepareHeartBeatDataWithBlock(selfId int32, peerMapJson string, addr string) data.HeartBeatData {

// 	mpt := p1.MerklePatriciaTrie{}
// 	mpt.Initial()
// 	listBlock := SBC.GetLatestBlocks()
// 	parentHash := listBlock[0].Header.Hash
// 	// TODO: change to TX id.
// 	currParentHashForPoW = parentHash
// 	parentHeight := listBlock[0].Header.Height

// 	// Do Proof of Work.
// 	nonce := doProofOfWork(parentHash, mpt.Root)
// 	if stopPoW {
// 		fmt.Println("Stop adding block.")
// 		return data.HeartBeatData{}
// 	}

// 	// Add a new block.
// 	timeStamp := time.Now().UnixNano()
// 	block := p2.NewBlock(parentHeight+1, timeStamp, parentHash, mpt, nonce)
// 	jsonBlockStr, _ := p2.EncodeToJson(block)
// 	fmt.Println("Added a new block.")
// 	fmt.Println("added block: ", jsonBlockStr)

// 	// Create a new HeartBeatData.
// 	addNewBlock := true
// 	heartBeatData := data.NewHeartBeatData(addNewBlock, selfId, jsonBlockStr, peerMapJson, addr)

// 	return heartBeatData
// }

// Do Proof of Work. Return the nonce that meets the PoW.
func doProofOfWork(parentHash string, mptRoot string) string {
	fmt.Println("doProofOfWork() start.")
	var nonceHash, nonce string
	for {
		// Check the stop flag.
		if stopPoW {
			fmt.Println("Stop calculating nonce.")
			break
		}

		// Calculate a nonceHash and check if it meets the requirement.
		nonce = getRandHexStr(16)
		nonceHash = calcNonceHash(parentHash, nonce, mptRoot)
		if strings.HasPrefix(nonceHash, NONCE_HASH_PREFIX) {
			fmt.Println("Nonce found! nonceHash:", nonceHash)
			fmt.Println("parentHash:", parentHash)
			fmt.Println("nonce:", nonce)
			fmt.Println("mptRoot:", mptRoot)
			break
		}
	}

	return nonce
}

// Calculate and return a nonceHash.
func calcNonceHash(parentHash string, nonce string, mptRoot string) string {
	str := parentHash + nonce + mptRoot
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

// Return a random hexadecimal string of length n.
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

// (Server) doGet() to print the current canonical chain. (/canonical)
func Canonical(w http.ResponseWriter, r *http.Request) {
	highestList := SBC.GetLatestBlocks()
	chainNum := 0
	for _, b := range highestList {
		chainNum++
		fmt.Printf("Chain #%d:\n", chainNum)
		DisplayAncestors(b)
		fmt.Println()
	}
}

// Display all the blocks from the specified block up to the first block.
func DisplayAncestors(b p2.Block) {
	for height := b.Header.Height; height > 0; height-- {
		jsonBlockStr, _ := p2.EncodeToJson(b)
		fmt.Printf("height=%d, %s\n", height, jsonBlockStr)
		parentBlock, _ := SBC.GetParentBlock(b)
		b = parentBlock
	}
}

// (User)(Server) doGet() to send the HeartBeatData. (/publish)
func Publish(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Publish() start.")

	// Prepare HeartBeatData with RentalInfo.
	peerMapJson, _ := Peers.PeerMapToJson()
	heartBeatData := prepareHbdWithRentalInfo(SELF_ID, peerMapJson, SELF_ADDR)

	// Send the HeartBeatData.
	closestAddresses := Peers.Rebalance()
	for _, addr := range closestAddresses {
		sendHeartBeat(addr, heartBeatData)
	}
}

// (User) Add a new RentalInfo as a transaction to HBD.
func prepareHbdWithRentalInfo(selfId int32, peerMapJson string, addr string) data.HeartBeatData {
	// Create a new RentalInfo.
	rentalInfo := p1.NewRentalInfo(SELF_ID, 50)
	jRentalInfo, _ := rentalInfo.EncodeToJson()

	// Create a new History.
	history := p1.NewHistory(rentalInfo.RentalId, 0, 0)
	jHistory, _ := history.EncodeToJson()
	fmt.Println("jHistory:", jHistory)

	// Create Transaction with the RentalInfo and TX fee.
	tx := data.NewTx(jRentalInfo, jHistory, 5)
	fmt.Println("TX created!")
	jTx, _ := tx.EncodeToJson()
	fmt.Println("TX created:", jTx)

	// Create a new HeartBeatData.
	addNewBlock := false
	addNewTx := true
	jsonBlockStr := ""
	heartBeatData := data.NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, addr)
	heartBeatData.TxId = tx.TxId

	return heartBeatData
}

// (User)(Server) doGet() to print the Ads. (/displayAds)
func DisplayAds(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DisplayAds() start.")

	highestBlockList := SBC.GetLatestBlocks()
	if highestBlockList == nil {
		fmt.Print("No blocks in the blockchain.\n")
		return
	}
	latestBlock := highestBlockList[0]

	// Test
	mptDb := latestBlock.Value.ValueDb
	for k, v := range mptDb {
		fmt.Println("k:", k)
		fmt.Println("v:", v)
	}

	// Display Ads.
	jRentalMap := mptDb["ads"]
	rentalMap, _ := p1.DecodeJsonToRentalMap(jRentalMap)
	for k, v := range rentalMap {
		fmt.Println("Rental ID:", k)
		fmt.Printf("Lender Id: %d, Asking Price: %d, Availability: %v\n", v.LenderId, v.AskingPrice, v.IsAvailable)
	}
}

// (User)(Server) doPost() to send StartRequest. (/sendStartRequest)
func SendStartRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SendStartRequest() start.")

	// Get CommandLineArgs from command line.
	args := getArgs(w, r)
	rentalId := args.RentalId
	depositStr := args.Deposit
	fmt.Println("rentalId to start:", rentalId)
	fmt.Println("deposit:", depositStr)

	// Prepare HeartBeatData with StartRequest.
	peerMapJson, _ := Peers.PeerMapToJson()
	heartBeatData := prepareHbdWithStartReqeust(rentalId, depositStr, SELF_ID, peerMapJson, SELF_ADDR)

	// Send the HeartBeatData.
	closestAddresses := Peers.Rebalance()
	for _, addr := range closestAddresses {
		sendHeartBeat(addr, heartBeatData)
	}
}

type CommandLineArgs struct {
	RentalId string `json:"rentalId"`
	Deposit  string `json:"deposit"`
}

// (Server) doPost() to get command line arguments.
func getArgs(w http.ResponseWriter, r *http.Request) CommandLineArgs {
	// fmt.Println("getArgs() start.")

	var args CommandLineArgs

	// HTTP Request.
	// Header.
	// fmt.Println("--- HTTP Request ----")
	method := r.Method
	// fmt.Println("[method] " + method)
	// for k, v := range r.Header {
	// 	fmt.Print("[header] " + k)
	// 	fmt.Println(": " + strings.Join(v, ","))
	// }

	// POST
	if method == "POST" {
		// fmt.Println("doPost() start.")
		defer r.Body.Close()
		// Get the CommandLineArgs.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("[request body](CommandLineArgs): " + string(body))

		// Unmarshal
		err = json.Unmarshal(body, &args)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
	}

	return args
}

// (User) Update the History, and put it in a transaction, and add it to HBD.
func prepareHbdWithStartReqeust(rentalId string, depositStr string, selfId int32, peerMapJson string, addr string) data.HeartBeatData {

	// Update the History.
	deposit, _ := strconv.Atoi(depositStr)
	history := retrieveHistory(rentalId, 0)
	history.BorrowerId = SELF_ID
	history.Deposit = int32(deposit)
	history.State = 1
	jHistory, _ := history.EncodeToJson()
	fmt.Println("jHistory:", jHistory)

	// Create Transaction with the StartRequest and TX fee.
	tx := data.NewTx("", jHistory, 5)
	fmt.Println("TX created!")
	jTx, _ := tx.EncodeToJson()
	fmt.Println("TX created:", jTx)

	// Create a new HeartBeatData.
	addNewBlock := false
	addNewTx := true
	jsonBlockStr := ""
	heartBeatData := data.NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, addr)
	heartBeatData.TxId = tx.TxId

	return heartBeatData
}

// (User)(Server) doGet() to print the history. (/displayHistory)
func DisplayHistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DisplayHistory() start.")

	highestBlockList := SBC.GetLatestBlocks()
	if highestBlockList == nil {
		fmt.Println("No block in the blockchain.")
		return
	}
	latestBlock := highestBlockList[0]

	// Display history.
	jHistoryMap := latestBlock.Value.ValueDb["history"]
	fmt.Println("jHistoryMap:", jHistoryMap)

	historyMap, _ := p1.DecodeJsonToHistoryMap(jHistoryMap)
	for k, v := range historyMap {
		fmt.Println("History ID:", k)
		fmt.Println("Rental ID:", v.RentalId)
		fmt.Println("Borrower ID:", v.BorrowerId)
		fmt.Println("Start Time:", v.StartTime)
		fmt.Println("End Time:", v.EndTime)
		fmt.Println("Deposit:", v.Deposit)
	}
}

// (User) Send updated information.
func Permit(history p1.History) {
	fmt.Println("Permit() start.")

	// Prepare HeartBeatData with updated information.
	peerMapJson, _ := Peers.PeerMapToJson()
	heartBeatData := prepareHbdWithStartTime(history, SELF_ID, peerMapJson, SELF_ADDR)

	// Send the HeartBeatData.
	closestAddresses := Peers.Rebalance()
	for _, addr := range closestAddresses {
		sendHeartBeat(addr, heartBeatData)
	}
}

// (User) Add updated information as a transaction to HBD.
func prepareHbdWithStartTime(history p1.History, selfId int32, peerMapJson string, addr string) data.HeartBeatData {
	// Set start time in the History.
	start := time.Now()
	// Network latency is taken into consideration for users.
	start = start.Add(time.Duration(5) * time.Second)
	history.StartTime = start
	history.State = 2
	jHistory, _ := history.EncodeToJson()

	// Set AvailabilityFlag to false.
	rentalInfo := setAvailabilityFlag(history.RentalId, false)
	jRentalInfo, _ := rentalInfo.EncodeToJson()

	// Create Transaction with the updated information and TX fee.
	tx := data.NewTx(jRentalInfo, jHistory, 5)
	fmt.Println("TX created!")
	jTx, _ := tx.EncodeToJson()
	fmt.Println("TX created:", jTx)

	// Create a new HeartBeatData.
	addNewBlock := false
	addNewTx := true
	jsonBlockStr := ""
	heartBeatData := data.NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, addr)
	heartBeatData.TxId = tx.TxId

	return heartBeatData
}

// Set the Availability Flag in RentalInfo.
func setAvailabilityFlag(rentalId string, b bool) p1.RentalInfo {
	fmt.Println("setAvailabilityFlag() start.")

	highestBlockList := SBC.GetLatestBlocks()
	if highestBlockList == nil {
		fmt.Print("No block in the blockchain.\n")
		return p1.RentalInfo{}
	}
	latestBlock := highestBlockList[0]
	jRentalMap := latestBlock.Value.ValueDb["ads"]
	fmt.Println("jRentalMap:", jRentalMap)

	rentalMap, _ := p1.DecodeJsonToRentalMap(jRentalMap)
	rentalInfo := rentalMap[rentalId]
	rentalInfo.IsAvailable = b

	return rentalInfo
}

// (User)(Server) doPost() to send StopRequest. (/sendStopRequest)
func SendStopRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SendStopRequest() start.")

	// Get CommandLineArgs from command line.
	args := getArgs(w, r)
	rentalId := args.RentalId
	fmt.Println("rentalId to stop:", rentalId)

	// Prepare HeartBeatData with StopRequest.
	peerMapJson, _ := Peers.PeerMapToJson()
	heartBeatData := prepareHbdWithStopReqeust(rentalId, SELF_ID, peerMapJson, SELF_ADDR)

	// Send the HeartBeatData.
	closestAddresses := Peers.Rebalance()
	for _, addr := range closestAddresses {
		sendHeartBeat(addr, heartBeatData)
	}
}

// (User) Update the History, and put it in a transaction, and add it to HBD.
func prepareHbdWithStopReqeust(rentalId string, selfId int32, peerMapJson string, addr string) data.HeartBeatData {
	// Update the History.
	history := retrieveHistory(rentalId, SELF_ID)
	history.State = 3
	jHistory, _ := history.EncodeToJson()
	fmt.Println("jHistory:", jHistory)

	// Create Transaction with the StopRequest and TX fee.
	tx := data.NewTx("", jHistory, 5)
	fmt.Println("TX created!")
	jTx, _ := tx.EncodeToJson()
	fmt.Println("TX created:", jTx)

	// Create a new HeartBeatData.
	addNewBlock := false
	addNewTx := true
	jsonBlockStr := ""
	heartBeatData := data.NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, addr)
	heartBeatData.TxId = tx.TxId

	return heartBeatData
}

// Retrieve the History from the blockchain.
func retrieveHistory(rentalId string, selfid int32) p1.History {
	fmt.Println("RetrieveHistory() start.")

	highestBlockList := SBC.GetLatestBlocks()
	if highestBlockList == nil {
		fmt.Print("No block in the blockchain.\n")
		return p1.History{}
	}
	latestBlock := highestBlockList[0]
	jHistoryMap := latestBlock.Value.ValueDb["history"]
	fmt.Println("jHistoryMap:", jHistoryMap)

	historyMap, _ := p1.DecodeJsonToHistoryMap(jHistoryMap)
	var history p1.History
	for k, v := range historyMap {
		if v.RentalId == rentalId && v.BorrowerId == selfid {
			history = historyMap[k]
		}
	}

	return history
}

// (User) Send updated information.
func Checkout(history p1.History) {
	fmt.Println("Checkout() start.")
	fmt.Println("HELLO!")

	// Prepare HeartBeatData with updated information.
	peerMapJson, _ := Peers.PeerMapToJson()
	heartBeatData := prepareHbdWithStopTime(history, SELF_ID, peerMapJson, SELF_ADDR)

	// Send the HeartBeatData.
	closestAddresses := Peers.Rebalance()
	for _, addr := range closestAddresses {
		sendHeartBeat(addr, heartBeatData)
	}
}

// (User) Add updated information as a transaction to HBD.
func prepareHbdWithStopTime(history p1.History, selfId int32, peerMapJson string, addr string) data.HeartBeatData {
	// Set AvailabilityFlag to true.
	rentalInfo := setAvailabilityFlag(history.RentalId, true)
	jRentalInfo, _ := rentalInfo.EncodeToJson()

	// Set start time in the History.
	stop := time.Now()
	// Network latency is taken into consideration for users.
	stop = stop.Add(time.Duration(-5) * time.Second)
	history.EndTime = stop
	history.State = 4

	// Calculate the service fee.
	elapse := int32((history.EndTime.Sub(history.StartTime)).Seconds())
	history.Deposit = history.Deposit - rentalInfo.AskingPrice*elapse
	jHistory, _ := history.EncodeToJson()

	// Create Transaction with the updated information and TX fee.
	tx := data.NewTx(jRentalInfo, jHistory, 5)
	fmt.Println("TX created!")
	jTx, _ := tx.EncodeToJson()
	fmt.Println("TX created:", jTx)

	// Create a new HeartBeatData.
	addNewBlock := false
	addNewTx := true
	jsonBlockStr := ""
	heartBeatData := data.NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, addr)
	heartBeatData.TxId = tx.TxId

	return heartBeatData
}

// Initialize the History.
func InitializeHistory(rentalId string) {
	fmt.Println("InitializeHistory() start.")

	// Prepare HeartBeatData with updated information.
	peerMapJson, _ := Peers.PeerMapToJson()
	heartBeatData := prepareHbdWithIniHistory(rentalId, SELF_ID, peerMapJson, SELF_ADDR)

	// Send the HeartBeatData.
	closestAddresses := Peers.Rebalance()
	for _, addr := range closestAddresses {
		sendHeartBeat(addr, heartBeatData)
	}
}

func prepareHbdWithIniHistory(rentalId string, selfId int32, peerMapJson string, addr string) data.HeartBeatData {
	// Create a new History.
	history := p1.NewHistory(rentalId, 0, 0)
	jHistory, _ := history.EncodeToJson()
	fmt.Println("jHistory:", jHistory)

	// Create Transaction with the RentalInfo and TX fee.
	tx := data.NewTx("", jHistory, 5)
	fmt.Println("TX created!")
	jTx, _ := tx.EncodeToJson()
	fmt.Println("TX created:", jTx)

	// Create a new HeartBeatData.
	addNewBlock := false
	addNewTx := true
	jsonBlockStr := ""
	heartBeatData := data.NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, addr)
	heartBeatData.TxId = tx.TxId

	return heartBeatData
}
