package cors

// Copyright (c) 2019 by Extreme Networks Inc.

import (
	"crypto/tls"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func StartCORSHandler(port string, router *mux.Router) {
	headersOk := handlers.AllowedHeaders([]string{
		"*",
		"Authorization",
		"X-Requested-With",
		"Content-Type",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{
		"GET",
		"HEAD",
		"PATCH",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS"})

	log.Fatal(http.ListenAndServe(":"+port,
		handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}

func StartCORSHandlerHTTPS(port string, router *mux.Router) {
	headersOk := handlers.AllowedHeaders([]string{
		"*",
		"Authorization",
		"X-Requested-With",
		"Content-Type",
	})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{
		"GET",
		"HEAD",
		"PATCH",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS"})

	// Disable security check for HTTPS
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// get certificates for HTTPS
	certificate, privkey, err := CertKeys()
	defer os.Remove(certificate) // clean up
	defer os.Remove(privkey)     // clean up
	if err != nil {
		log.Fatal("Cannot locate certificates for HTTPS")
	}
	log.Fatal(http.ListenAndServeTLS(":"+port,
		certificate, privkey,
		handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}
