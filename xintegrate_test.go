/*
Copyright (c) 2019 Dave Hammers
*/
package main

import (
	"go-dalb/node"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// an integation test that:
// creates <n> HTTP worker nodes
// starts the DALB
// sends <y> HTTP requests
// reports the results
func TestIntegration(t *testing.T) {
	//initialize the data path server
	proxy := dataPathInit("/")
	// create 10 worker nodes
	ipList, err := net.LookupIP("localhost")
	if err != nil {
		t.Fatal("Cannot lookup localhost")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	for port := 9000; port < 9010; port++ {
		go func(port int) {
			log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
		}(port)
		n := node.NewNode()
		n.IP = ipList[0]
		n.Port = port
		n.MaxTransactions = 10
		proxy.sched.SchedAddNode(n)
	}

	// there should be 10 worker nodes running. They have all been added to the scheduler.
	// Send transactions as fast as possible to see how this works
	r := httptest.NewRequest("GET", "http://localhost/", nil)
	w := httptest.NewRecorder()
	tStart := time.Now()
	for cnt := 0; cnt < 10000; cnt++ {
		proxy.router.ServeHTTP(w, r)
	}
	t.Log("Elapse time", time.Since(tStart))
	t.Log("Number of transactions", proxy.sched.TransactionCount())
	t.Log("Average transaction time", proxy.sched.AverageTransactionTime())
	min, max := proxy.sched.TransactionTimeRange()
	t.Log("Min transaction time", min)
	t.Log("Max transaction time", max)
	for n := range proxy.sched.SchedNodeMap {
		t.Log("======================================")
		t.Log("Node:", n.IP.String(), ":", n.Port)
		t.Log("Number of transactions", n.TransactionCount())
		t.Log("Average transaction time", n.AverageTransactionTime())
		min, max := n.TransactionTimeRange()
		t.Log("Min transaction time", min)
		t.Log("Max transaction time", max)
	}
}
