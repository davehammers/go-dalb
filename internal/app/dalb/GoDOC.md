package dalb // import "dalb/internal/app/dalb"

Copyright (c) 2019 Dave Hammers

Copyright (c) 2019 Dave Hammers

FUNCTIONS

func AddRoutes(Router *mux.Router) *mux.Router
func CtrlPathInit() (Router *mux.Router)
func SchedStatsGet(w http.ResponseWriter, r *http.Request)

TYPES

type AddNode struct {
	Path            string `json:"path"`
	Address         string `json:"address"`
	Port            int    `json:"port"`
	MaxTransactions int    `json:"maxTransactions"`
}

type DataPathProxy struct {
	Proxy  *httputil.ReverseProxy
	Router *mux.Router
	Sched  *node.Scheduler
	// Has unexported fields.
}

var (
	Proxy *DataPathProxy
)
func DataPathInit(path string) *DataPathProxy

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

