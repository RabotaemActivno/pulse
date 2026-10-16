package checker

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type StatResult struct {
	HttpStatus int
	Status     string
}

type CheckResult struct {
	CheckedAt  time.Time
	HTTPStatus int
	Latency    time.Duration
	Err        error
}

type Checker struct {
	client *http.Client
	logger *slog.Logger
}

func New(c *http.Client) *Checker {
	return &Checker{
		client: c,
	}
}

func (c *Checker) Check(ctx context.Context, url string, method string) (CheckResult, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return CheckResult{}, err
	}

	start := time.Now()

	res, err := c.client.Do(req)
	if err != nil {
		return CheckResult{Latency: time.Since(start), Err: err}, nil
	}

	defer res.Body.Close()

	_, _ = io.Copy(io.Discard, res.Body)

	return CheckResult{
		Latency:    time.Since(start),
		HTTPStatus: res.StatusCode,
		CheckedAt:  time.Now(),
		Err:        err,
	}, nil
}
