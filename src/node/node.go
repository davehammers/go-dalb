package node

import (
	"net"
	"time"
)

type Node struct {
	ID     uint64
	IP     net.IP
	Port   int
	Weight int
	stat   struct {
		totalTransactions    int64
		totalTransactionTime time.Duration
		minTransactionTime   time.Duration
		maxTransactionTime   time.Duration
	}
}

var (
	nextID uint64
)

// Returns a new *Node with the ID initialized to a unique number.
func NewNode() *Node {
	nextID++
	return &Node{ID: nextID}
}

// After a transaction is complete, update the node with the time.Duration it took to process the transaction
func (n *Node) UpdateTime(duration time.Duration) {
	n.stat.totalTransactions++
	n.stat.totalTransactionTime += duration
	if n.stat.minTransactionTime == 0 || n.stat.minTransactionTime > duration {
		n.stat.minTransactionTime = duration
	}
	if n.stat.maxTransactionTime < duration {
		n.stat.maxTransactionTime = duration
	}
}

// Initialize the node statistics
func (n *Node) Reset() {
	n.stat.totalTransactions = 0
	n.stat.totalTransactionTime = 0
	n.stat.minTransactionTime = 0
	n.stat.maxTransactionTime = 0
}

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
