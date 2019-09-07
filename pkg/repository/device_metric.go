package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"go-postgres-redis/pkg/model"
	"time"
)

type DeviceMetric struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewDeviceMetric(db *sql.DB) *DeviceMetric {
	return &DeviceMetric{
		db,
		squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}

func (r *DeviceMetric) Find(limit, offset int, from time.Time) ([]*model.DeviceMetric, error) {
	q := r.sb.Select(
		"id", "device_id", "metric_1", "metric_2", "metric_3", "metric_4", "metric_5",
		"local_time", "server_time",
	).From(
		"device_metrics",
	).OrderBy(
		"server_time ASC",
	).Where(
		squirrel.GtOrEq{
			"server_time": from,
		},
	).Limit(uint64(limit)).Offset(uint64(offset))

	rows, err := q.Query()
	if err != nil {
		return nil, err
	}

	var mm []*model.DeviceMetric
	for rows.Next() {
		m := new(model.DeviceMetric)
		err = rows.Scan(
			&m.ID, &m.DeviceID, &m.Metric1, &m.Metric2, &m.Metric3, &m.Metric4, &m.Metric5,
			&m.LocalTime, &m.ServerTime,
		)
		if err != nil {
			return nil, err
		}
		mm = append(mm, m)
	}

	return mm, rows.Close()
}

func (r *DeviceMetric) Count(from time.Time) (int, error) {
	q := r.sb.Select(
		"count(id)",
	).From(
		"device_metrics",
	).Where(
		squirrel.GtOrEq{
			"server_time": from,
		},
	)
	var count int
	return count, q.QueryRow().Scan(&count)
}
