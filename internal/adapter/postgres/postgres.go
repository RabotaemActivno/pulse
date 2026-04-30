package postgres

import "github.com/RabotaemActivno/pulse/pkg/postgres"

type Postgres struct {
	pgPool *postgres.Pool
}

func New(pool *postgres.Pool) *Postgres {
	return &Postgres{pgPool: pool}
}
