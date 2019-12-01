/*
Copyright (c) 2019 Dave Hammers
*/
package main

import (
	"flag"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-dalb/src/cors"
	"os"
)

// CONSTANTS
const (
	DefaultDataPort    = "8080"
	DefaultControlPort = "8081"
)

// Environment variable oneDEBUG=1 sets the package DEBUG=true
const EnvDebug = "dalbDEBUG"

var (
	// command line flags
	debugPtr *bool
	portPtr  *string
	httpPtr  *bool
	// env variables
	DEBUG bool
)

func init() {
	envSetup()
}
func envSetup() {
	if d, ok := os.LookupEnv(EnvDebug); ok {
		switch d {
		case "0":
			DEBUG = false
		case "1":
			DEBUG = true
		}
	}
}

func mainStart() (port string, router *mux.Router) {
	// environment variables
	defaultPort := DefaultDataPort
	envPort, ok := os.LookupEnv("PORT")
	if ok {
		defaultPort = envPort
	}

	// get command line parameters
	debugPtr = flag.Bool("d", false, "standalone debug development (no container)")
	portPtr = flag.String("p", defaultPort, "HTTPS listens on this port")
	httpPtr = flag.Bool("http", false, "Use HTTP instead of HTTPS")
	flag.Parse()

	port = *portPtr

	if *debugPtr {
		log.SetLevel(log.DebugLevel)
	}

	// some parts of GCP want the applications prefix while others don't
	// register both GCP and standalone versions of the URL
	// save setting
	router = mux.NewRouter().StrictSlash(false)
	//router = AddRoutes(router)

	if DEBUG {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
	log.Println("Server started at http://localhost:", port)
	return
}

func main() {
	port, router := mainStart()
	if *httpPtr {
		cors.StartCORSHandler(port, router)
	} else {
		cors.StartCORSHandlerHTTPS(port, router)
	}
}
