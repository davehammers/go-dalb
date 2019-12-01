package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-dalb/node"
	"net/http"
	"net/http/httputil"
	"time"
)

type dataPathProxy struct {
	proxy  *httputil.ReverseProxy
	sched  *node.Scheduler
	router *mux.Router
}

func dataPathInit() *dataPathProxy {
	//create a reverse proxy that distributes the requests to the worker nodes
	dpProxy := &dataPathProxy{}
	dpProxy.proxy = &httputil.ReverseProxy{Director: dpProxy.dataPathDirector}
	//allocate a load balancer scheduler for the data path
	dpProxy.sched = node.NewScheduler(0)
	//load any pre-configured worker node definitions
	//TODO

	dpProxy.router = mux.NewRouter().StrictSlash(false)
	dpProxy.router.HandleFunc("/{path:.*}", dpProxy.dataPathForward)
	return dpProxy
}

//direct the request to the next available worker node
//STUB - work is done in dataPathForward
func (p *dataPathProxy) dataPathDirector(r *http.Request) {
}

func (p *dataPathProxy) dataPathForward(w http.ResponseWriter, r *http.Request) {
	n := p.sched.SchedGetNode()
	if n != nil {
		r.URL.Host = fmt.Sprintf("%s:%d", n.IP.String(), n.Port)
	} else {
		r.URL.Host = ""
		log.Error("Cannot get a worker node for request")
	}
	tStart := time.Now()
	p.proxy.ServeHTTP(w, r)
	// make the node available for another request
	p.sched.SchedRescheduleNode(n)
	// compute how long the worker node took to complete the transaction
	tDur := time.Since(tStart)
	//update node stats
	n.UpdateTime(tDur)
	//update scheduler stats
	p.sched.UpdateTime(tDur)

}
