package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ctrlPathInit() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(false)
	//router = AddRoutes(router)

	if *pDebug {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
	return
}
