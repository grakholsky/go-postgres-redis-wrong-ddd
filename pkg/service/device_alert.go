package service

import (
	"go-postgres-redis/pkg/model"
	"go-postgres-redis/pkg/repository"
)

type DeviceAlert struct {
	repository *repository.DeviceAlert
}

func NewDeviceAlert(r *repository.DeviceAlert) *DeviceAlert {
	return &DeviceAlert{r}
}

func (s *DeviceAlert) Put(m model.DeviceAlert) (ID int, err error) {
	return s.repository.Put(&m)
}

func (s *DeviceAlert) PutToCache(userID int, m model.DeviceAlert) error {
	return s.repository.PutToCache(userID, &m)
}
