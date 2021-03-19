/*
Copyright (c) 2019 Dave Hammers
*/
package node

import (
	"sync"
	"time"
)

const (
	//Schedule length determines how many node transactions can be Scheduled at a time.
	//E.g. node.MaxTransactions * number of worker nodes = Schedule len
	DefaultScheduleLen = 1000
	//The Schedule rebalancer examines the performance of the worker nodes periodically.
	DefaultRebalanceMinutes = 15
)

type SchedNodeMapType map[*Node]bool
type SchedChannel chan *Node
type Scheduler struct {
	SchedNodeMap    SchedNodeMapType
	lock            sync.Mutex
	nodeChannel     SchedChannel
	statsChan       chan time.Duration
	rebalanceTicker *time.Ticker

	stat struct {
		totalTransactions    int64
		totalTransactionTime time.Duration
		minTransactionTime   time.Duration
		maxTransactionTime   time.Duration
	}
}

//Return a new Scheduler used to Schedule traffic to Nodes
func NewScheduler(SchedLen int) *Scheduler {
	if SchedLen == 0 {
		SchedLen = DefaultScheduleLen
	}
	s := &Scheduler{
		SchedNodeMap:    make(SchedNodeMapType),
		lock:            sync.Mutex{},
		statsChan:       make(chan time.Duration, 1000),
		nodeChannel:     make(SchedChannel, SchedLen),
		rebalanceTicker: time.NewTicker(time.Minute * time.Duration(DefaultRebalanceMinutes)),
	}
	// this go routine listens on a Scheduler channel for transaction durations
	// it offloads any Scheduler statistics updates from the main program path
	go func(s *Scheduler) {
		for duration := range s.statsChan {
			s.stat.totalTransactions++
			s.stat.totalTransactionTime += duration
			if s.stat.minTransactionTime == 0 || s.stat.minTransactionTime > duration {
				s.stat.minTransactionTime = duration
			}
			if s.stat.maxTransactionTime < duration {
				s.stat.maxTransactionTime = duration
			}
		}
	}(s)
	go func(s *Scheduler) {
		for range s.rebalanceTicker.C {
			s.SchedRebalance()
		}
	}(s)

	return s
}

//delete a Scheduler by closing it's active channels
func (s *Scheduler) Delete() {
	// stop the rebalancer ticker
	s.rebalanceTicker.Stop()
	// close the channel used to update the statistics
	close(s.statsChan)
	// close the channel used to Schedule worker nodes
	close(s.nodeChannel)
	// init the Schedule map to release any references to *Node(s)
	s.SchedNodeMap = nil
}

//add node to the distribution Schedule n.MaxTransactions times
// initially this will cause the node to be Scheduled back to back. Over time, as transactions are processed
// this will distribute itself into the Schedule with the other nodes.
func (s *Scheduler) SchedAddNode(n *Node) {
	s.lock.Lock()
	s.SchedNodeMap[n] = true
	s.lock.Unlock()
	for idx := 0; idx < n.MaxTransactions; idx++ {
		s.nodeChannel <- n
	}
}

//returns the next *Node that should be used for a reverse proxy request
func (s *Scheduler) SchedGetNode() *Node {
	// get the next worker node to be used
	for n := range s.nodeChannel {
		// check to verify the node is still valid.
		// it can be deleted if the node is removed from service or the Scheduler is rebalanced.
		_, ok := s.SchedNodeMap[n]
		if ok {
			return n
		}
		// fall thru means the node has been deleted and should not be used anymore
		// get the next one
	}
	return nil
}

//re-adds the *Node to the end of the Schedule
func (s *Scheduler) SchedReScheduleNode(n *Node) {
	s.nodeChannel <- n
}

//deletes a node from the Scheduler map. This will eventually remove the node from the nodeChannel
func (s *Scheduler) SchedDeleteNode(n *Node) {
	s.lock.Lock()
	delete(s.SchedNodeMap, n)
	s.lock.Unlock()
}

//Periodically examine the the performance of each worker node to see if some nodes are
//out performing others. For the nodes that are underperforming shift the workloads to other
//faster nodes by deleting the slower node and re-adding it with a lower MaxTransactions value.
func (s *Scheduler) SchedRebalance() {
	//TODO
}

//
// S T A T I S T I C S
//

// After a transaction is complete, update the Scheduler with the time.Duration it took to process the transaction
func (s *Scheduler) UpdateTime(duration time.Duration) {
	s.statsChan <- duration
}

// Initialize the Scheduler statistics
func (s *Scheduler) Reset() {
	s.stat.totalTransactions = 0
	s.stat.totalTransactionTime = 0
	s.stat.minTransactionTime = 0
	s.stat.maxTransactionTime = 0
}

// returns the average transaction time for this Scheduler
func (s *Scheduler) AverageTransactionTime() time.Duration {
	if s.stat.totalTransactions == 0 {
		return 0
	}
	return time.Duration(s.stat.totalTransactionTime.Nanoseconds() / s.stat.totalTransactions)
}

// Returns the number of transactions processed by a Scheduler
func (s *Scheduler) TransactionCount() int64 {
	return s.stat.totalTransactions
}

// returns the total time.Duration for all transactions processed by a Scheduler
func (s *Scheduler) TransactionTime() time.Duration {
	return s.stat.totalTransactionTime
}

// returns the minimum and maximum time.Duration for all transactions processed by a Scheduler
func (s *Scheduler) TransactionTimeRange() (time.Duration, time.Duration) {
	return s.stat.minTransactionTime, s.stat.maxTransactionTime
}
