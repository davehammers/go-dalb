# dalb/internal/node package

This package manages the worker nodes and scheduler for an Application Load Balancer instance.

```sh
package node // import "dalb/internal/node"

Copyright (c) 2019 Dave Hammers

CONSTANTS

const (
	//Schedule length determines how many node transactions can be Scheduled at a time.
	//E.g. node.MaxTransactions * number of worker nodes = Schedule len
	DefaultScheduleLen = 1000
	//The Schedule rebalancer examines the performance of the worker nodes periodically.
	DefaultRebalanceMinutes = 15
)

TYPES

type Node struct {
	IP              net.IP
	Port            int
	MaxTransactions int

	// Has unexported fields.
}

func NewNode() *Node
    Returns a new *Node with the ID initialized to a unique number.

func (n *Node) AverageTransactionTime() time.Duration
    returns the average transaction time for this node

func (n *Node) Delete()
    delete a Node by closing its active structures

func (n *Node) Reset()
    Initialize the node statistics

func (n *Node) TransactionCount() int64
    Returns the number of transactions processed by a node

func (n *Node) TransactionTime() time.Duration
    returns the total time.Duration for all transactions processed by a node

func (n *Node) TransactionTimeRange() (time.Duration, time.Duration)
    returns the minimum and maximum time.Duration for all transactions processed
    by a node

func (n *Node) UpdateTime(duration time.Duration)
    After a transaction is complete, update the node with the time.Duration it
    took to process the transaction

type scheduler struct {
	SchedNodeMap SchedNodeMapType

	// Has unexported fields.
}

func NewScheduler(SchedLen int) *scheduler
    Return a new scheduler used to Schedule traffic to Nodes

func (s *scheduler) AverageTransactionTime() time.Duration
    returns the average transaction time for this scheduler

func (s *scheduler) Delete()
    delete a scheduler by closing it's active channels

func (s *scheduler) Reset()
    Initialize the scheduler statistics

func (s *scheduler) SchedAddNode(n *Node)
    add node to the distribution Schedule n.MaxTransactions times initially this
    will cause the node to be Scheduled back to back. Over time, as transactions
    are processed this will distribute itself into the Schedule with the other
    nodes.

func (s *scheduler) SchedDeleteNode(n *Node)
    deletes a node from the scheduler map. This will eventually remove the node
    from the nodeChannel

func (s *scheduler) SchedGetNode() *Node
    returns the next *Node that should be used for a reverse proxy request

func (s *scheduler) SchedRebalance()
    Periodically examine the the performance of each worker node to see if some
    nodes are out performing others. For the nodes that are underperforming
    shift the workloads to other faster nodes by deleting the slower node and
    re-adding it with a lower MaxTransactions value.

func (s *scheduler) SchedReScheduleNode(n *Node)
    re-adds the *Node to the end of the Schedule

func (s *scheduler) TransactionCount() int64
    Returns the number of transactions processed by a scheduler

func (s *scheduler) TransactionTime() time.Duration
    returns the total time.Duration for all transactions processed by a
    scheduler

func (s *scheduler) TransactionTimeRange() (time.Duration, time.Duration)
    returns the minimum and maximum time.Duration for all transactions processed
    by a scheduler

func (s *scheduler) UpdateTime(duration time.Duration)
    After a transaction is complete, update the scheduler with the time.Duration
    it took to process the transaction
```
