/*
Copyright (c) 2019 Dave Hammers
*/
package dalb

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"dalb/internal/node"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type DataPathProxy struct {
	path   string
	Proxy  *httputil.ReverseProxy
	Router *mux.Router
	Sched  *node.Scheduler
}

var (
	Proxy *DataPathProxy
)

func DataPathInit(path string) *DataPathProxy {
	//create a reverse Proxy that distributes the requests to the worker nodes
	dpProxy := &DataPathProxy{
		path: path,
	}
	dpProxy.Proxy = &httputil.ReverseProxy{Director: dpProxy.dataPathDirector}
	//allocate a load balancer scheduler for the data path
	dpProxy.Sched = node.NewScheduler(0)
	//load any pre-configured worker node definitions
	//TODO

	dpProxy.Router = mux.NewRouter().StrictSlash(false)
	dpProxy.Router.HandleFunc(path, dpProxy.dataPathForward)
	Proxy = dpProxy
	return dpProxy
}

//direct the request to the next available worker node
//STUB - work is done in dataPathForward
func (p *DataPathProxy) dataPathDirector(r *http.Request) {
}

func (p *DataPathProxy) dataPathForward(w http.ResponseWriter, r *http.Request) {
	n := p.Sched.SchedGetNode()
	if n != nil {
		r.URL.Host = fmt.Sprintf("%s:%d", n.IP.String(), n.Port)
	} else {
		r.URL.Host = ""
		log.Error("Cannot get a worker node for request")
	}
	tStart := time.Now()
	p.Proxy.ServeHTTP(w, r)
	// make the node available for another request
	p.Sched.SchedReScheduleNode(n)
	// compute how long the worker node took to complete the transaction
	tDur := time.Since(tStart)
	//update node stats
	n.UpdateTime(tDur)
	//update scheduler stats
	p.Sched.UpdateTime(tDur)

}
