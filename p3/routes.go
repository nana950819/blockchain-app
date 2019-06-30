/*
	CS686 Project 3 and 4: Build a Gossip Network and a simple PoW.
	Route is a handler in HTTP server.
	Author: Kei Fukutani
	Date  : March 31, Apr 6, 2019
*/
package p3

import "net/http"

// Hold the API information of a handler, such as
// name, HTTP method, path, and the function name.
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Miner's handlers.
var RoutesMiner = Routes{
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"POST",
		"/upload",
		Upload,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"Start",
		"GET",
		"/start",
		Start,
	},
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	// For testing
	Route{
		"Hello",
		"POST",
		"/hello",
		Hello,
	},
}

// User's handlers.
var RoutesUser = Routes{
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
	Route{
		"Upload",
		"POST",
		"/upload",
		Upload,
	},
	Route{
		"UploadBlock",
		"GET",
		"/block/{height}/{hash}",
		UploadBlock,
	},
	Route{
		"HeartBeatReceiveForUser",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceiveForUser,
	},
	Route{
		"StartUser",
		"GET",
		"/start",
		StartUser,
	},
	Route{
		"Canonical",
		"GET",
		"/canonical",
		Canonical,
	},
	Route{
		"Publish",
		"GET",
		"/publish",
		Publish,
	},
	Route{
		"DisplayAds",
		"GET",
		"/displayAds",
		DisplayAds,
	},
	Route{
		"SendStartRequest",
		"POST",
		"/sendStartRequest",
		SendStartRequest,
	},
	Route{
		"DisplayHistory",
		"GET",
		"/displayHistory",
		DisplayHistory,
	},
	Route{
		"SendStopRequest",
		"POST",
		"/sendStopRequest",
		SendStopRequest,
	},
	// For testing
	Route{
		"Hello",
		"POST",
		"/hello",
		Hello,
	},
}
