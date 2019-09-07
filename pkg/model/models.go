package model

import "time"

type (
	Device struct {
		ID     int
		Name   string
		UserID int
	}

	DeviceInfo struct {
		UserName  string
		UserEmail string
		Device
	}

	DeviceMetric struct {
		ID         int
		DeviceID   int
		Metric1    *int
		Metric2    *int
		Metric3    *int
		Metric4    *int
		Metric5    *int
		LocalTime  *time.Time
		ServerTime *time.Time
	}

	DeviceAlert struct {
		ID       int
		DeviceID int
		Message  string
	}
)
