package scheduler_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/RabotaemActivno/pulse/internal/domain"
	"github.com/RabotaemActivno/pulse/internal/scheduler"
	"github.com/RabotaemActivno/pulse/pkg/checker"
	"github.com/google/uuid"
)

var config = scheduler.Config{Workers: 5, RetryBase: 2 * time.Second}
var httpClient = http.DefaultClient
var ch = checker.New(httpClient)

type TestCheckCreator struct {
}

func (cr TestCheckCreator) CreateCheck(ctx context.Context, m domain.Monitor, res checker.CheckResult) error {
	return nil
}

func Test_StartScheduler_Success(t *testing.T) {
	tcr := TestCheckCreator{}
	s := scheduler.New(config, ch, tcr)
	monitor := domain.Monitor{}
	s.Start()
	err := s.Submit(monitor)
	if err != nil {
		t.Errorf("err submit: %v", err)
	}

	<-time.After(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}

func Test_SubmitMonitor(t *testing.T) {
	tcr := TestCheckCreator{}
	m := domain.Monitor{
		ID:             uuid.New(),
		URL:            "https://httpbin.org/status/200",
		Method:         http.MethodGet,
		ExpectedStatus: http.StatusOK,
	}

	s := scheduler.New(config, ch, tcr)
	s.Start()
	err := s.Submit(m)
	if err != nil {
		t.Errorf("err submit: %v", err)
	}

	<-time.After(1 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	s.Shutdown(ctx)
}
