package service

import (
	"go-postgres-redis/pkg/model"
	"go-postgres-redis/pkg/repository"
)

type Device struct {
	repository *repository.Device
}

func NewDevice(r *repository.Device) *Device {
	return &Device{r}
}

func (s *Device) DeviceInfo(ID int) (*model.DeviceInfo, error) {
	return s.repository.DeviceInfo(ID)
}
