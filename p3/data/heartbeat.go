/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	Hold information for Heart Beat.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package data

// Hold the content of the Heart Beat.
type HeartBeatData struct {
	IfNewBlock bool `json:"ifNewBlock"`
	IfNewTx    bool `json:"ifNewTx"`
	// Sender's Id
	Id        int32  `json:"id"`
	BlockJson string `json:"blockJson"`
	TxJson    string `json:"txJson"`
	TxId      string `json:"txId"`
	// JSON format of PeerList.PeerMap, the output of function PeerList.PeerMapToJson().
	PeerMapJson string `json:"peerMapJson"`
	// Sender's Addr
	Addr string `json:"addr"`
	Hops int32  `json:"hops"`
}

// Create a new HeartBeatData.
func NewHeartBeatData(ifNewBlock bool, ifNewTx bool, id int32, jsonBlockStr string, jTx string, peerMapJson string, addr string) HeartBeatData {
	heartBeatData := HeartBeatData{IfNewBlock: ifNewBlock, IfNewTx: ifNewTx, Id: id, BlockJson: jsonBlockStr,
		TxJson: jTx, PeerMapJson: peerMapJson, Addr: addr, Hops: 2}
	return heartBeatData
}

// PrepareHeartBeatData would first create a new instance of HeartBeatData,
// then decide whether or not you will create a new block and send the new block to other peers.
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapJson string, addr string) HeartBeatData {
	addNewBlock := false
	addNewTx := false
	jsonBlockStr := ""
	jTx := ""

	// // For Project 3, Add a new block randomly.
	// if rand.Float32() < 0.5 {
	// 	addNewBlock = true
	// 	mpt := p1.MerklePatriciaTrie{}
	// 	block := sbc.GenBlock(mpt)
	// 	jsonBlockStr, _ = p2.EncodeToJson(block)
	// 	fmt.Println("Added a new block.")
	// 	fmt.Println("added block: ", jsonBlockStr)
	// }

	heartBeatData := NewHeartBeatData(addNewBlock, addNewTx, selfId, jsonBlockStr, jTx, peerMapJson, addr)

	return heartBeatData
}
