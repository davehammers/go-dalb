/*
Copyright (c) 2019 Dave Hammers
*/
package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go-dalb/cors"
)

// CONSTANTS
const (
	DefaultDataPort    = "8080"
	DefaultControlPort = "8081"
)

var (
	// command line flags
	pDebug    *bool
	pDataPort *string
	pCtrlPort *string
	pHttp     *bool
)

func commandLineInit() {
	// get command line parameters
	if pDebug == nil {
		pDebug = flag.Bool("d", false, "enable debug logging output")
		pDataPort = flag.String("data", DefaultDataPort, "HTTP listens on this port for datapath requests")
		pCtrlPort = flag.String("ctrl", DefaultControlPort, "HTTP listens on this port for control requests")
		pHttp = flag.Bool("http", false, "Use HTTP instead of HTTPS")
	}
	flag.Parse()
}

func main() {
	commandLineInit()
	// start the control HTTP server
	go func() {
		router := ctrlPathInit()
		if *pHttp {
			log.Debug("Server started at http://localhost:", *pCtrlPort)
			cors.StartCORSHandler(*pCtrlPort, router)
		} else {
			log.Debug("Server started at https://localhost:", *pCtrlPort)
			cors.StartCORSHandlerHTTPS(*pCtrlPort, router)
		}
	}()

	//start data path server
	proxy := dataPathInit()
	if *pHttp {
		log.Debug("Server started at http://localhost:", *pDataPort)
		cors.StartCORSHandler(*pDataPort, proxy.router)
	} else {
		log.Debug("Server started at https://localhost:", *pDataPort)
		cors.StartCORSHandlerHTTPS(*pDataPort, proxy.router)
	}

}
