package postgres

import "github.com/RabotaemActivno/pulse/pkg/postgres"

type Postgres struct {
	PgPool *postgres.Pool
}

func New(pool *postgres.Pool) *Postgres {
	return &Postgres{PgPool: pool}
}
