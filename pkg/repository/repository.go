package repository

import (
	"database/sql"
	"github.com/garyburd/redigo/redis"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/logrusadapter"
	"github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
	"time"
)

type PostgresConfig struct {
	URI string
}

type RedisConfig struct {
	Addr        string
	MaxIdle     int
	IdleTimeout time.Duration
}

func ConnectPostgres(logger *logrus.Logger, cfg *PostgresConfig) (*sql.DB, error) {
	c, err := pgx.ParseURI(cfg.URI)
	if err != nil {
		return nil, err
	}
	c.Logger = logrusadapter.NewLogger(logger)
	var db *sql.DB
	for {
		db = stdlib.OpenDB(c)
		if err = db.Ping(); err != nil {
			logger.Error(err)
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}
	return db, nil
}

func ConnectRedis(logger *logrus.Logger, cfg *RedisConfig) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		IdleTimeout: cfg.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", cfg.Addr)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return pool, nil
}

func IsNotFound(err error) bool {
	return err == sql.ErrNoRows
}
