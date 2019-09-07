package repository

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/garyburd/redigo/redis"
	"go-postgres-redis/pkg/model"
)

const (
	userDeviceAlertFormatKey = "user_id_%d_device_id_%d_alert"
)

type DeviceAlert struct {
	db    *sql.DB
	cache *redis.Pool
	sb    squirrel.StatementBuilderType
}

func NewDeviceAlert(db *sql.DB, cache *redis.Pool) *DeviceAlert {
	return &DeviceAlert{
		db,
		cache,
		squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}

func (r *DeviceAlert) Put(m *model.DeviceAlert) (int, error) {
	q := r.sb.Insert(
		"device_alerts",
	).Columns(
		"device_id", "message",
	).Values(
		m.DeviceID, m.Message,
	).Suffix(
		"RETURNING \"id\"",
	)
	var ID int
	return ID, q.QueryRow().Scan(&ID)
}

func (r *DeviceAlert) PutToCache(userID int, m *model.DeviceAlert) error {
	c := r.cache.Get()
	defer c.Close()
	key := fmt.Sprintf(userDeviceAlertFormatKey, userID, m.DeviceID)
	_, err := c.Do("SET", key, m.Message)
	return err
}
