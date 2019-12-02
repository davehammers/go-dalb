/*
Copyright (c) 2019 Dave Hammers
*/
package node

import (
	"net"
	"time"
)

type Node struct {
	IP              net.IP
	Port            int
	MaxTransactions int
	statsChan       chan time.Duration
	stat            struct {
		totalTransactions    int64
		totalTransactionTime time.Duration
		minTransactionTime   time.Duration
		maxTransactionTime   time.Duration
	}
}

// Returns a new *Node with the ID initialized to a unique number.
func NewNode() *Node {
	n := &Node{
		statsChan: make(chan time.Duration, 1000),
	}
	// this go routine listens on a node channel for transaction durations
	// it offloads any node statistics updates from the main program path
	go func(n *Node) {
		for duration := range n.statsChan {
			n.stat.totalTransactions++
			n.stat.totalTransactionTime += duration
			if n.stat.minTransactionTime == 0 || n.stat.minTransactionTime > duration {
				n.stat.minTransactionTime = duration
			}
			if n.stat.maxTransactionTime < duration {
				n.stat.maxTransactionTime = duration
			}
		}
	}(n)
	return n
}

//delete a Node by closing its active structures
func (n *Node) Delete() {
	close(n.statsChan)
}

// After a transaction is complete, update the node with the time.Duration it took to process the transaction
func (n *Node) UpdateTime(duration time.Duration) {
	n.statsChan <- duration
}

//
// S T A T I S T I C S
//

// Initialize the node statistics
func (n *Node) Reset() {
	n.stat.totalTransactions = 0
	n.stat.totalTransactionTime = 0
	n.stat.minTransactionTime = 0
	n.stat.maxTransactionTime = 0
}

// returns the average transaction time for this node
func (n *Node) AverageTransactionTime() time.Duration {
	if n.stat.totalTransactions == 0 {
		return 0
	}
	return time.Duration(n.stat.totalTransactionTime.Nanoseconds() / n.stat.totalTransactions)
}

// Returns the number of transactions processed by a node
func (n *Node) TransactionCount() int64 {
	return n.stat.totalTransactions
}

// returns the total time.Duration for all transactions processed by a node
func (n *Node) TransactionTime() time.Duration {
	return n.stat.totalTransactionTime
}

// returns the minimum and maximum time.Duration for all transactions processed by a node
func (n *Node) TransactionTimeRange() (time.Duration, time.Duration) {
	return n.stat.minTransactionTime, n.stat.maxTransactionTime
}
