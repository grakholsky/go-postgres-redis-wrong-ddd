package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"go-postgres-redis/pkg/model"
)

type Device struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

func NewDevice(db *sql.DB) *Device {
	return &Device{
		db,
		squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}
}

func (r *Device) DeviceInfo(deviceID int) (*model.DeviceInfo, error) {
	q := r.sb.Select(
		"users.name", "users.email",
		"devices.id", "devices.name", "devices.user_id",
	).From(
		"devices",
	).Join(
		"users ON users.id = devices.user_id",
	).Where(
		squirrel.Eq{
			"devices.id": deviceID,
		},
	)
	m := new(model.DeviceInfo)
	row := q.QueryRow()
	return m, row.Scan(&m.UserName, &m.UserEmail, &m.ID, &m.Name, &m.UserID)
}
