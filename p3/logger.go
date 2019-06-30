/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	This is for logging.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package p3

import (
	"log"
	"net/http"
	"time"
)

// HTTP Request and Response logger.
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
