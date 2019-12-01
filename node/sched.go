package node

import (
	"sync"
	"time"
)

const (
	//schedule length determines how many node transactions can be scheduled at a time.
	//E.g. node.MaxTransactions * number of worker nodes = schedule len
	DefaultScheduleLen = 1000
	//The schedule rebalancer examines the performance of the worker nodes periodically.
	DefaultRebalanceMinutes = 15
)

type schedNodeMapType map[*Node]bool
type schedChannel chan *Node
type Scheduler struct {
	SchedNodeMap    schedNodeMapType
	lock            sync.Mutex
	nodeChannel     schedChannel
	statsChan       chan time.Duration
	rebalanceTicker *time.Ticker

	stat struct {
		totalTransactions    int64
		totalTransactionTime time.Duration
		minTransactionTime   time.Duration
		maxTransactionTime   time.Duration
	}
}

//Return a new scheduler used to schedule traffic to Nodes
func NewScheduler(schedLen int) *Scheduler {
	if schedLen == 0 {
		schedLen = DefaultScheduleLen
	}
	s := &Scheduler{
		SchedNodeMap:    make(schedNodeMapType),
		lock:            sync.Mutex{},
		statsChan:       make(chan time.Duration, 1000),
		nodeChannel:     make(schedChannel, schedLen),
		rebalanceTicker: time.NewTicker(time.Minute * time.Duration(DefaultRebalanceMinutes)),
	}
	// this go routine listens on a scheduler channel for transaction durations
	// it offloads any scheduler statistics updates from the main program path
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

//delete a scheduler by closing it's active channels
func (s *Scheduler) Delete() {
	// stop the rebalancer ticker
	s.rebalanceTicker.Stop()
	// close the channel used to update the statistics
	close(s.statsChan)
	// close the channel used to schedule worker nodes
	close(s.nodeChannel)
	// init the schedule map to release any references to *Node(s)
	s.SchedNodeMap = nil
}

//add node to the distribution schedule n.MaxTransactions times
// initially this will cause the node to be scheduled back to back. Over time, as transactions are processed
// this will distribute itself into the schedule with the other nodes.
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
		// it can be deleted if the node is removed from service or the scheduler is rebalanced.
		_, ok := s.SchedNodeMap[n]
		if ok {
			return n
		}
		// fall thru means the node has been deleted and should not be used anymore
		// get the next one
	}
	return nil
}

//re-adds the *Node to the end of the schedule
func (s *Scheduler) SchedRescheduleNode(n *Node) {
	s.nodeChannel <- n
}

//deletes a node from the scheduler map. This will eventually remove the node from the nodeChannel
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

// After a transaction is complete, update the scheduler with the time.Duration it took to process the transaction
func (s *Scheduler) UpdateTime(duration time.Duration) {
	s.statsChan <- duration
}

// Initialize the scheduler statistics
func (s *Scheduler) Reset() {
	s.stat.totalTransactions = 0
	s.stat.totalTransactionTime = 0
	s.stat.minTransactionTime = 0
	s.stat.maxTransactionTime = 0
}

// returns the average transaction time for this scheduler
func (s *Scheduler) AverageTransactionTime() time.Duration {
	if s.stat.totalTransactions == 0 {
		return 0
	}
	return time.Duration(s.stat.totalTransactionTime.Nanoseconds() / s.stat.totalTransactions)
}

// Returns the number of transactions processed by a scheduler
func (s *Scheduler) TransactionCount() int64 {
	return s.stat.totalTransactions
}

// returns the total time.Duration for all transactions processed by a scheduler
func (s *Scheduler) TransactionTime() time.Duration {
	return s.stat.totalTransactionTime
}

// returns the minimum and maximum time.Duration for all transactions processed by a scheduler
func (s *Scheduler) TransactionTimeRange() (time.Duration, time.Duration) {
	return s.stat.minTransactionTime, s.stat.maxTransactionTime
}
