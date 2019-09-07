package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"go-postgres-redis/pkg/repository"
	"go-postgres-redis/pkg/service"
	"time"
)

type EnvConfig struct {
	WatcherInterval    string `envconfig:"WATCHER_INTERVAL"`
	RedisURL           string `envconfig:"REDIS_URL"`
	PostgresURI        string `envconfig:"POSTGRES_URI"`
	EmailPort          int    `envconfig:"EMAIL_PORT"`
	EmailHost          string `envconfig:"EMAIL_HOST"`
	EmailHostUser      string `envconfig:"EMAIL_HOST_USER"`
	EmailHostPassword  string `envconfig:"EMAIL_HOST_PASSWORD"`
	EmailFrom          string `envconfig:"EMAIL_FROM"`
	EmailTlsSkipVerify bool   `envconfig:"EMAIL_TLS_SKIP_VERIFY"`
	MaxMetric1         int    `envconfig:"MAX_METRIC_1"`
	MaxMetric2         int    `envconfig:"MAX_METRIC_2"`
	MaxMetric3         int    `envconfig:"MAX_METRIC_3"`
	MaxMetric4         int    `envconfig:"MAX_METRIC_4"`
	MaxMetric5         int    `envconfig:"MAX_METRIC_5"`
}

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
	}
	var envCfg EnvConfig
	err := envconfig.Process("", &envCfg)
	if err != nil {
		logger.Fatal(err)
	}
	watcherInterval, err := time.ParseDuration(envCfg.WatcherInterval)
	if err != nil {
		logger.Fatalf("parse 'WATCHER_INTERVAL' error: %s", err)
	}
	watcherConfig := &service.WatcherConfig{
		Interval:    watcherInterval,
		PostgresCfg: &repository.PostgresConfig{URI: envCfg.PostgresURI},
		RedisCfg: &repository.RedisConfig{
			Addr:        envCfg.RedisURL,
			MaxIdle:     10,
			IdleTimeout: 10 * time.Second,
		},
		NotifyConfig: &service.NotifyConfig{
			Host:          envCfg.EmailHost,
			HostUser:      envCfg.EmailHostUser,
			HostPwd:       envCfg.EmailHostPassword,
			From:          envCfg.EmailFrom,
			Port:          envCfg.EmailPort,
			TlsSkipVerify: envCfg.EmailTlsSkipVerify,
		},
		DeviceMetricConfig: &service.DeviceMetricConfig{
			MaxMetric1: envCfg.MaxMetric1,
			MaxMetric2: envCfg.MaxMetric2,
			MaxMetric3: envCfg.MaxMetric3,
			MaxMetric4: envCfg.MaxMetric4,
			MaxMetric5: envCfg.MaxMetric5,
		},
	}
	watcher := service.NewWatcher(logger, watcherConfig)
	watcher.Run()
}
