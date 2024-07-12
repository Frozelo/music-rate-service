package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

type Postgres struct {
	conn         *pgx.Conn
	connAttempts int
	connTimeout  time.Duration
}

func New(connString string) (*Postgres, error) {
	pg := &Postgres{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	var err error
	for pg.connAttempts > 0 {
		pg.conn, err = pgx.Connect(context.Background(), connString)
		if err == nil {
			break
		}
		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts-1)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (s *Postgres) Close() {
	s.conn.Close(context.Background())
}
