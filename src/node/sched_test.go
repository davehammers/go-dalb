package node

import (
	"testing"
	"time"
)

var (
	tSched *Scheduler
)

func TestNewScheduler(t *testing.T) {
	tSched = NewScheduler(0)
	s1 := NewScheduler(0)
	if s1 == tSched {
		t.Fatal("Sched ID is not unique", s1)
	}
}

func TestScheduler_SchedAddNode(t *testing.T) {
	tSched = NewScheduler(0)
	for i := 0; i < 5; i++ {
		n := NewNode()
		n.MaxTransactions = 1
		tSched.SchedAddNode(n)
	}
	if len(tSched.SchedNodeMap) != 5 {
		t.Fatal("scheduler Node count is not correct")
	}
}

func TestScheduler_SchedGetNode(t *testing.T) {
	t.Log("Node channel len", len(tSched.nodeChannel))
	n := tSched.SchedGetNode()
	if len(tSched.nodeChannel) != 4 {
		t.Fatal("scheduler channel is not correct")
	}
	tSched.SchedRescheduleNode(n)
	if len(tSched.nodeChannel) != 5 {
		t.Fatal("scheduler channel is not correct")
	}
}

func TestScheduler_SchedDeleteNode(t *testing.T) {
	for n := range tSched.SchedNodeMap {
		tSched.SchedDeleteNode(n)
	}
	if len(tSched.SchedNodeMap) != 0 {
		t.Fatal("scheduler Node count is not correct after delete")
	}
}

func TestSched_UpdateTime(t *testing.T) {
	durTotal := int64(0)
	for cnt, dur := range []int64{12, 8, 30, 15} {
		tSched.UpdateTime(time.Duration(dur) * time.Nanosecond)
		time.Sleep(1 * time.Millisecond)
		if tSched.TransactionCount() != int64(cnt+1) {
			t.Fatal("transaction count not incrementing")
		}
		durTotal += dur
		durAvg := durTotal / int64(cnt+1)
		if tSched.AverageTransactionTime() != time.Duration(durAvg) {
			t.Fatal("Average transaction time not computed correctly", tSched.AverageTransactionTime())
		}
	}
	min, max := tSched.TransactionTimeRange()
	t.Log("min=", min, ", max=", max)
	if min != 8 || max != 30 {
		t.Fatal("transaction min/max not set correctly")
	}
}

func TestSched_Reset(t *testing.T) {
	//t.Logf("Before reset %#v\n", *tSched)
	tSched.Reset()
	//t.Logf("After reset %#v\n", *tSched)
	if tSched.TransactionTime() != 0 {
		t.Fatal("total transaction time is not zero")
	}
	if tSched.TransactionCount() != 0 {
		t.Fatal("total transaction count is not zero")
	}
	if tSched.AverageTransactionTime() != 0 {
		t.Fatal("Average transaction time is not zero")
	}
	min, max := tSched.TransactionTimeRange()
	if min != 0 {
		t.Fatal("min transaction time is not zero")
	}
	if max != 0 {
		t.Fatal("max transaction time is not zero")
	}

}
func TestScheduler_Delete(t *testing.T) {
	tSched.Delete()
}
