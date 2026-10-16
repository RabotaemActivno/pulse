package scheduler

import (
	"container/heap"
	"time"
)

type TimeHeap []*Task

func (h TimeHeap) Len() int {
	return len(h)
}

func (h TimeHeap) Less(i, j int) bool {
	return h[i].monitor.NextCheckAt.Before(h[j].monitor.NextCheckAt)
}

func (h TimeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *TimeHeap) Push(x any) {
	*h = append(*h, x.(*Task))
}

func (h *TimeHeap) Pop() any {
	old := *h
	n := len(old)
	t := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return t
}

func (s *Scheduler) dispatcher() {
	defer close(s.out)

	var th TimeHeap
	in := s.in

	for {

		if in == nil && th.Len() == 0 {
			return
		}

		var wait <-chan time.Time
		var timer *time.Timer
		if th.Len() > 0 {
			d := time.Until(th[0].monitor.NextCheckAt)
			if d < 0 {
				d = 0
			}
			timer = time.NewTimer(d)
			wait = timer.C
		}

		select {
		case <-s.workCtx.Done():
			stopTimer(timer)
			return
		case t, ok := <-in:
			stopTimer(timer)
			if !ok {
				in = nil
				continue
			}
			heap.Push(&th, t)
		case <-wait:
			ready := heap.Pop(&th).(*Task)
			select {
			case s.out <- ready:
			case <-s.workCtx.Done():
				return
			}
		}
	}
}

func stopTimer(t *time.Timer) {
	if t != nil {
		t.Stop()
	}
}
