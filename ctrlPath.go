/*
Copyright (c) 2019 Dave Hammers
*/
package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-dalb/node"
	"net"
	"net/http"
	"time"
)

type route struct {
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type routes []route

var dalbRoutes = routes{
	route{
		"GET",
		"/scheduler",
		schedStatsGet,
	},
	route{
		"GET",
		"/node",
		nodeStatsGet,
	},
	route{
		"POST",
		"/node",
		nodePost,
	},
}

func ctrlPathInit() (router *mux.Router) {
	router = mux.NewRouter().StrictSlash(false)
	router = AddRoutes(router)

	if *pDebug {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}
	return
}

func AddRoutes(router *mux.Router) *mux.Router {
	for _, route := range dalbRoutes {
		r := router.NewRoute()
		r.Methods(route.Method)
		r.Path(route.Pattern)
		r.Handler(route.HandlerFunc)
	}
	return router
}

// OBSERVABILITY
// stats to see how our scheduller and individual worker nodes are doing
type SchedulerStats struct {
	Path                           string  `json:"path"`
	TransactionCount               int64   `json:"transactionCount"`
	AverageTransactionTimeMilliSec float64 `json:"averageTransactionTimeMilliSec"`
	MinimumTransactionTimeMilliSec float64 `json:"minimumTransactionTimeMilliSec"`
	MaximumTransactionTimeMilliSec float64 `json:"maximumTransactionTimeMilliSec"`
}

func schedStatsGet(w http.ResponseWriter, r *http.Request) {
	min, max := proxy.sched.TransactionTimeRange()
	stat := SchedulerStats{
		Path:                           proxy.path,
		TransactionCount:               proxy.sched.TransactionCount(),
		AverageTransactionTimeMilliSec: float64(proxy.sched.AverageTransactionTime() / time.Millisecond),
		MinimumTransactionTimeMilliSec: float64(min / time.Millisecond),
		MaximumTransactionTimeMilliSec: float64((max / time.Millisecond)),
	}
	json.NewEncoder(w).Encode(stat)
}

type NodeStats struct {
	Nodes []Nodes `json:"nodes"`
}
type Nodes struct {
	Address                        string  `json:"address"`
	Port                           int     `json:"port"`
	MaxTransactions                int     `json:"maxTransactions"`
	TransactionCount               int64   `json:"transactionCount"`
	AverageTransactionTimeMilliSec float64 `json:"averageTransactionTimeMilliSec"`
	MinimumTransactionTimeMilliSec float64 `json:"minimumTransactionTimeMilliSec"`
	MaximumTransactionTimeMilliSec float64 `json:"maximumTransactionTimeMilliSec"`
}

func nodeStatsGet(w http.ResponseWriter, r *http.Request) {
	stats := NodeStats{
		Nodes: make([]Nodes, 0),
	}
	for n := range proxy.sched.SchedNodeMap {
		min, max := n.TransactionTimeRange()
		node := Nodes{
			Address:                        n.IP.String(),
			Port:                           n.Port,
			MaxTransactions:                n.MaxTransactions,
			TransactionCount:               n.TransactionCount(),
			AverageTransactionTimeMilliSec: float64(n.AverageTransactionTime() / time.Millisecond),
			MinimumTransactionTimeMilliSec: float64(min / time.Millisecond),
			MaximumTransactionTimeMilliSec: float64((max / time.Millisecond)),
		}
		stats.Nodes = append(stats.Nodes, node)
	}
	json.NewEncoder(w).Encode(stats)
}

type AddNode struct {
	Path            string `json:"path"`
	Address         string `json:"address"`
	Port            int    `json:"port"`
	MaxTransactions int    `json:"maxTransactions"`
}

//Add a worker node to the <path> scheduler
func nodePost(w http.ResponseWriter, r *http.Request) {
	newNode := &AddNode{}
	err := json.NewDecoder(r.Body).Decode(newNode)
	if err != nil {
		http.Error(w, "JSON format error", http.StatusBadRequest)
		return
	}
	ipList, err := net.LookupIP(newNode.Address)
	if err != nil {
		http.Error(w, "invalid IP address", http.StatusBadRequest)
		return
	}
	n := node.NewNode()
	n.IP = ipList[0]
	n.Port = newNode.Port
	n.MaxTransactions = newNode.MaxTransactions
	proxy.sched.SchedAddNode(n)
}
