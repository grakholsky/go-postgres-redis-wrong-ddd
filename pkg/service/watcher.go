package service

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-postgres-redis/pkg/model"
	"go-postgres-redis/pkg/repository"
	"time"
)

const limit int = 100

type (
	WatcherConfig struct {
		Interval           time.Duration
		PostgresCfg        *repository.PostgresConfig
		RedisCfg           *repository.RedisConfig
		NotifyConfig       *NotifyConfig
		DeviceMetricConfig *DeviceMetricConfig
	}

	Watcher struct {
		config    *WatcherConfig
		logger    *logrus.Logger
		dvc       *Device
		dvcMetric *DeviceMetric
		dvcAlert  *DeviceAlert
		notify    *Notify
	}
)

func NewWatcher(logger *logrus.Logger, config *WatcherConfig) *Watcher {
	pgDb, err := repository.ConnectPostgres(logger, config.PostgresCfg)
	if err != nil {
		logger.Fatalf("service: connect to postgres error: %s", err)
	}
	redisDb, err := repository.ConnectRedis(logger, config.RedisCfg)
	if err != nil {
		logger.Fatalf("service: connect to redis error: %s", err)
	}
	w := new(Watcher)
	w.config = config
	w.logger = logger
	w.dvc = NewDevice(repository.NewDevice(pgDb))
	w.dvcMetric = NewDeviceMetric(repository.NewDeviceMetric(pgDb), config.DeviceMetricConfig)
	w.dvcAlert = NewDeviceAlert(repository.NewDeviceAlert(pgDb, redisDb))
	w.notify = NewNotify(logger, config.NotifyConfig)
	return w
}

func (s *Watcher) Run() {
	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()
	from := time.Now()
	var offset int
	for {
		select {
		case now := <-ticker.C:
			count, err := s.dvcMetric.Count(from)
			if err != nil {
				s.logger.Errorf("service: watcher device metric count error: %v", err)
				break
			}
			offset = 0
			for ; count > 0; count -= limit {
				mm, err := s.dvcMetric.Find(limit, offset, from)
				if err != nil {
					s.logger.Errorf("service: watcher device metric find error: %v", err)
					break
				}
				for _, m := range mm {
					if err := s.dvcMetric.WarnValidate(m); err != nil {
						if err = s.doAlert(m, err); err != nil {
							s.logger.Errorf("service: watcher sending alert error: %v", err)
						}
					}
				}
				offset = limit
			}
			from = now
		}
	}
}

func (s *Watcher) doAlert(m *model.DeviceMetric, alert error) error {
	mInfo, err := s.dvc.DeviceInfo(m.DeviceID)
	if err != nil {
		s.logger.Errorf("service: watcher retrieving device info error: %v", err)
		return err
	}

	msg := fmt.Sprintf(
		"User name: %s, Device name: %s alert: %s",
		mInfo.UserName, mInfo.Name, alert.Error(),
	)
	mAlert := &model.DeviceAlert{
		DeviceID: m.DeviceID,
		Message:  msg,
	}
	ID, err := s.dvcAlert.Put(*mAlert)
	if err != nil {
		s.logger.Errorf("service: watcher inserting device alert error: %v", err)
		return err
	}

	mAlert.ID = ID
	if err := s.dvcAlert.PutToCache(mInfo.UserID, *mAlert); err != nil {
		s.logger.Errorf("service: watcher inserting to cache device alert error: %v", err)
		return err
	}

	go s.notify.SendMail(mInfo.UserEmail, "Device Alert!", msg)
	return nil
}
