/*
	Main function for Project 3 and 4.
	Create a node and make it stand-by.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
	Note  : The first node address must be "localhost:6686".
			The second node address must be "localhost:6687".
			The nodes after those must be any address except the first's, the second's and
			"localhost:6688".
*/
package main

import (
	"cs686/cs686-blockchain-p3-kayfuku/p3"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Create a node and make it stand-by.
// Run the API {URL}/start to launch the stand-by node.
func main() {
	fmt.Println("HELLO!!")

	// Create a node.
	var localPort string
	var router *mux.Router
	var nodeType string

	if len(os.Args) == 3 {
		localPort = os.Args[1]

		if os.Args[2] == "m" {
			// For miners.
			nodeType = "Miner"
			router = p3.NewRouter(p3.RoutesMiner)

		} else if os.Args[2] == "u" {
			// For users.
			nodeType = "User"
			router = p3.NewRouter(p3.RoutesUser)

		}

	} else if len(os.Args) == 1 {
		localPort = "7000"
		nodeType = "User"
		router = p3.NewRouter(p3.RoutesUser)
	} else {
		fmt.Println("Error. Command line args needed. ")
		fmt.Println("Usage: <program name> <port number> {m/u}")
		fmt.Println("Example: ./main 6686 m")
		return
	}
	p3.SELF_PORT = localPort

	// Start listening.
	fmt.Println(nodeType + " Node " + localPort + " stand by.")
	log.Fatal(http.ListenAndServe(":"+localPort, router))

}
