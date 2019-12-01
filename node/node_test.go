package node

import (
	"net"
	"testing"
	"time"
)

var (
	tNode *Node
)

func TestNewNode(t *testing.T) {
	tNode = NewNode()
	n1 := NewNode()
	if n1 == tNode {
		t.Fatal("Node ID is not unique", n1)
	}
	tNode.IP = net.IPv4(192, 168, 10, 100)
	tNode.Port = 9001
	tNode.MaxTransactions = 20
}

func TestNode_UpdateTime(t *testing.T) {
	durTotal := int64(0)
	for cnt, dur := range []int64{12, 8, 30, 15} {
		tNode.UpdateTime(time.Duration(dur) * time.Nanosecond)
		time.Sleep(1 * time.Millisecond)
		if tNode.TransactionCount() != int64(cnt+1) {
			t.Fatal("transaction count not incrementing")
		}
		durTotal += dur
		durAvg := durTotal / int64(cnt+1)
		if tNode.AverageTransactionTime() != time.Duration(durAvg) {
			t.Fatal("Average transaction time not computed correctly", tNode.AverageTransactionTime())
		}
	}
	min, max := tNode.TransactionTimeRange()
	t.Log("min=", min, ", max=", max)
	if min != 8 || max != 30 {
		t.Fatal("transaction min/max not set correctly")
	}
}

func TestNode_Reset(t *testing.T) {
	//t.Logf("Before reset %#v\n", *tNode)
	tNode.Reset()
	//t.Logf("After reset %#v\n", *tNode)
	if tNode.TransactionTime() != 0 {
		t.Fatal("total transaction time is not zero")
	}
	if tNode.TransactionCount() != 0 {
		t.Fatal("total transaction count is not zero")
	}
	if tNode.AverageTransactionTime() != 0 {
		t.Fatal("Average transaction time is not zero")
	}
	min, max := tNode.TransactionTimeRange()
	if min != 0 {
		t.Fatal("min transaction time is not zero")
	}
	if max != 0 {
		t.Fatal("max transaction time is not zero")
	}

}
