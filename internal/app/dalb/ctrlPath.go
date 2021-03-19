/*
Copyright (c) 2019 Dave Hammers
*/
package dalb

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"dalb/internal/node"

	"github.com/gorilla/mux"
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
		SchedStatsGet,
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

func CtrlPathInit() (Router *mux.Router) {
	Router = mux.NewRouter().StrictSlash(false)
	Router = AddRoutes(Router)

	return
}

func AddRoutes(Router *mux.Router) *mux.Router {
	for _, route := range dalbRoutes {
		r := Router.NewRoute()
		r.Methods(route.Method)
		r.Path(route.Pattern)
		r.Handler(route.HandlerFunc)
	}
	return Router
}

// OBSERVABILITY
// stats to see how our Scheduller and individual worker nodes are doing
type schedulerStats struct {
	Path                           string  `json:"path"`
	TransactionCount               int64   `json:"transactionCount"`
	AverageTransactionTimeMilliSec float64 `json:"averageTransactionTimeMilliSec"`
	MinimumTransactionTimeMilliSec float64 `json:"minimumTransactionTimeMilliSec"`
	MaximumTransactionTimeMilliSec float64 `json:"maximumTransactionTimeMilliSec"`
}

func SchedStatsGet(w http.ResponseWriter, r *http.Request) {
	min, max := Proxy.Sched.TransactionTimeRange()
	stat := schedulerStats{
		Path:                           Proxy.path,
		TransactionCount:               Proxy.Sched.TransactionCount(),
		AverageTransactionTimeMilliSec: float64(Proxy.Sched.AverageTransactionTime() / time.Millisecond),
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
	for n := range Proxy.Sched.SchedNodeMap {
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
	Proxy.Sched.SchedAddNode(n)
}
