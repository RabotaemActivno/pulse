package grpc

import (
	"context"
	"time"

	ssov1 "github.com/RabotaemActivno/sso-grpc/gen/go/sso"
)

type Client struct {
	api ssov1.AuthClient
}

func New(
	ctx context.Context,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {

	return &Client{}, nil
}
