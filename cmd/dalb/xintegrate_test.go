/*
Copyright (c) 2019 Dave Hammers
*/
package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"dalb/internal/app/dalb"
	"dalb/internal/node"
)

// an integation test that:
// creates <n> HTTP worker nodes
// starts the DALB
// sends <y> HTTP requests
// reports the results
func TestIntegration(t *testing.T) {
	//initialize the data path server
	proxy := dalb.DataPathInit("/")
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
		n.MaxTransactions = port - 8999
		proxy.Sched.SchedAddNode(n)
	}

	// there should be 10 worker nodes running. They have all been added to the scheduler.
	// Send transactions as fast as possible to see how this works
	r := httptest.NewRequest("GET", "http://localhost/", nil)
	w := httptest.NewRecorder()
	tStart := time.Now()
	for cnt := 0; cnt < 15000; cnt++ {
		proxy.Router.ServeHTTP(w, r)
	}
	t.Log("\nscheduler")
	t.Log("Elapse time", time.Since(tStart))
	t.Log("Number of transactions", proxy.Sched.TransactionCount())
	t.Log("Average transaction time", proxy.Sched.AverageTransactionTime())
	min, max := proxy.Sched.TransactionTimeRange()
	t.Log("Min transaction time", min)
	t.Log("Max transaction time", max)
	t.Log("\nWorker Nodes")
	for n := range proxy.Sched.SchedNodeMap {
		t.Log("======================================")
		t.Log("Node:", n.IP.String(), ":", n.Port)
		t.Log("MaxTransactions", n.MaxTransactions)
		t.Log("Number of transactions", n.TransactionCount())
		t.Log("Average transaction time", n.AverageTransactionTime())
		min, max := n.TransactionTimeRange()
		t.Log("Min transaction time", min)
		t.Log("Max transaction time", max)
	}
}
