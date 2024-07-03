package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/AxelTahmid/golang-starter/config"
)

type Postgres struct {
	DB *pgxpool.Pool
}

var (
	pool   *Postgres
	pgOnce sync.Once
)

func initConfig(conf config.Database) *pgxpool.Config {

	dbConfig, err := pgxpool.ParseConfig(conf.Url)

	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}

	dbConfig.MaxConns = conf.PoolMax
	dbConfig.MinConns = conf.PoolMin
	dbConfig.MaxConnLifetime = conf.MaxConnLifetime
	dbConfig.MaxConnIdleTime = conf.MaxConnIdleTime
	dbConfig.HealthCheckPeriod = conf.HealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = conf.ConnectTimeout

	dbConfig.AfterConnect = setDbTimeZone

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	if os.Getenv("ENV") == "development" {
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	}

	dbConfig.ConnConfig.Tracer = &pgTracer{
		log: &log,
	}

	return dbConfig
}

func CreatePool(ctx context.Context, conf config.Database) (*Postgres, error) {
	var err error

	pgOnce.Do(func() {
		dbPool, dbErr := pgxpool.NewWithConfig(context.Background(), initConfig(conf))
		if dbErr != nil {
			err = dbErr
		}

		pool = &Postgres{dbPool}
	})

	return pool, err
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}

func setDbTimeZone(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, "SET TIME ZONE 'UTC';")

	if err != nil {
		return fmt.Errorf("unable to set timezone: %w", err)
	}
	log.Printf("Timezone set to UTC\n")
	return nil
}
