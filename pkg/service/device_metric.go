package service

import (
	"fmt"
	"go-postgres-redis/pkg/model"
	"go-postgres-redis/pkg/repository"
	"strings"
	"time"
)

type (
	DeviceMetricConfig struct {
		MaxMetric1 int
		MaxMetric2 int
		MaxMetric3 int
		MaxMetric4 int
		MaxMetric5 int
	}
	DeviceMetric struct {
		repository *repository.DeviceMetric
		config     *DeviceMetricConfig
	}
)

func NewDeviceMetric(r *repository.DeviceMetric, config *DeviceMetricConfig) *DeviceMetric {
	return &DeviceMetric{r, config}
}

func (s *DeviceMetric) Find(limit, offset int, from time.Time) ([]*model.DeviceMetric, error) {
	return s.repository.Find(limit, offset, from)
}

func (s *DeviceMetric) Count(from time.Time) (int, error) {
	return s.repository.Count(from)
}

func (s *DeviceMetric) WarnValidate(m *model.DeviceMetric) error {
	var errs []string
	if m.Metric1 != nil && *m.Metric1 >= s.config.MaxMetric1 {
		errs = append(errs, fmt.Sprintf("metric1 has been: %d", *m.Metric1))
	}
	if m.Metric2 != nil && *m.Metric2 >= s.config.MaxMetric2 {
		errs = append(errs, fmt.Sprintf("metric2 has been: %d", *m.Metric2))
	}
	if m.Metric3 != nil && *m.Metric3 >= s.config.MaxMetric3 {
		errs = append(errs, fmt.Sprintf("metric3 has been: %d", *m.Metric3))
	}
	if m.Metric4 != nil && *m.Metric4 >= s.config.MaxMetric4 {
		errs = append(errs, fmt.Sprintf("metric4 has been: %d", *m.Metric4))
	}
	if m.Metric5 != nil && *m.Metric5 >= s.config.MaxMetric5 {
		errs = append(errs, fmt.Sprintf("metric5 has been: %d", *m.Metric5))
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}
	return nil
}
