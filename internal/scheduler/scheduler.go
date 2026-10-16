package scheduler

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/pkg/checker"
)

var (
	ErrQueueFull = errors.New("scheduler: queue full")
)

type CheckCreator interface {
	CreateCheck(ctx context.Context, m domain.Monitor, res checker.CheckResult) error
}

type Config struct {
	Workers   int
	RetryBase time.Duration
}

type Task struct {
	monitor  domain.Monitor
	attempts int
}

type Scheduler struct {
	cfg Config

	creator CheckCreator
	c       *checker.Checker

	in  chan *Task
	out chan *Task

	mu sync.Mutex

	submitted atomic.Int64
	pending   atomic.Int64
	failed    atomic.Int64
	panicked  atomic.Int64
	retried   atomic.Int64
	dropped   atomic.Int64

	wg         sync.WaitGroup
	workCtx    context.Context
	cancelWork context.CancelFunc
}

func New(cfg Config, c *checker.Checker, creator CheckCreator) *Scheduler {
	s := Scheduler{
		cfg:     cfg,
		c:       c,
		creator: creator,
	}

	s.workCtx, s.cancelWork = context.WithCancel(context.Background())

	s.in = make(chan *Task, 64)
	s.out = make(chan *Task)

	return &s
}

func (s *Scheduler) Start() {
	for i := 0; i < s.cfg.Workers; i++ {
		s.wg.Add(1)
		go s.worker()
	}
	go s.dispatcher()
}

func (s *Scheduler) Submit(m domain.Monitor) error {
	t := Task{monitor: m}
	select {
	case s.in <- &t:
		s.submitted.Add(1)
		s.pending.Add(1)
	default:
		s.dropped.Add(1)
		return ErrQueueFull
	}

	return nil
}

func (s *Scheduler) worker() {
	fmt.Println("worker started")
	defer s.wg.Done()
	for {
		select {
		case <-s.workCtx.Done():
			return
		case t, ok := <-s.out:
			if !ok {
				return
			}
			err := s.safeProcess(s.workCtx, *t)
			if err != nil {
				s.failed.Add(1)
				continue
			}
		}
	}
}

func (s *Scheduler) retry(ctx context.Context, t Task) error {
	var err error
	attempts := t.monitor.ConsecutiveFailures
	if attempts <= 0 {
		attempts = 1
	}
	for i := 0; i < attempts; i++ {
		if err = s.check(ctx, t); err == nil {
			return nil
		}
		s.retried.Add(1)
		if i == t.monitor.ConsecutiveFailures-1 {
			break
		}
		t.attempts++
		d := s.cfg.RetryBase * time.Duration(1<<i)
		d += time.Duration(rand.Int64N(int64(d / 2)))
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return err
}

func (s *Scheduler) check(ctx context.Context, t Task) error {
	res, err := s.c.Check(ctx, t.monitor.URL, t.monitor.Method)
	if err != nil {
		return err
	}
	err = s.creator.CreateCheck(ctx, t.monitor, res)
	if err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) safeProcess(ctx context.Context, t Task) (err error) {
	fmt.Println("safe process called")
	defer func() {
		if r := recover(); r != nil {
			s.panicked.Add(1)
			err = fmt.Errorf("panicked: %v", r)
		}
	}()
	return s.retry(ctx, t)
}

func (s *Scheduler) Shutdown(ctx context.Context) {

	close(s.in)

	wait := make(chan interface{})
	go func() {
		s.wg.Wait()
		wait <- ""
	}()

	select {
	case <-ctx.Done():
		s.cancelWork()
	case <-wait:
		fmt.Println("stopped gracefully")
	}
}
