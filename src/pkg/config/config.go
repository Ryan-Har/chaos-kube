package config

import (
	"fmt"
	"github.com/Ryan-Har/chaos-kube/pkg/streams"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	RedisConfig  RedisConfig
	RedisStreams ReadRedisStreams
	Hostname     string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type ReadRedisStreams struct {
	ConsumerGroup   string
	ConsumerStreams []string
}

func init() {
	logLevel := getLogLevel()
	var addSource bool
	if logLevel == slog.LevelDebug {
		addSource = true
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: addSource,
		Level:     logLevel,
	})

	slog.SetDefault(slog.New(jsonHandler))

	slog.Info("logger initialised", "level", logLevel.String())
}

func Load(serviceName string) (*Config, error) {
	slog.Info(fmt.Sprintf("loading config for %s", serviceName))
	var config Config
	config.Hostname = getHostname()
	rc, err := loadRedisConfig()
	if err != nil {
		return &config, err
	}
	config.RedisConfig = rc

	streams, err := loadRedisStreams(serviceName)
	if err != nil {
		return &config, err
	}
	config.RedisStreams = streams

	return &config, nil
}

func getLogLevel() slog.Level {
	level := getEnvWithDefault("LOG_LEVEL", "info")
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// checks if an env var exists, defaults to fallback if it does not
func getEnvWithDefault(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		slog.Debug("failed to load environment variable so using default", "key", key)
		return fallback
	}
	slog.Debug("successfully loaded environment variable", "key", key)
	return value
}

// Returns the hostname if it can get it, otherwise an empty string
func getHostname() string {
	s, _ := os.Hostname()
	return s
}

func loadRedisConfig() (RedisConfig, error) {
	var rc RedisConfig
	redisHost := getEnvWithDefault("REDIS_HOST", "localhost")
	redisPort := getEnvWithDefault("REDIS_PORT", "6379")
	rc.Addr = fmt.Sprintf("%s:%s", redisHost, redisPort)
	rc.Password = getEnvWithDefault("REDIS_PASSWORD", "")
	db, err := strconv.Atoi(getEnvWithDefault("REDIS_DB", "0"))
	if err == nil {
		rc.DB = db
	}
	return rc, err
}

func loadRedisStreams(serviceName string) (ReadRedisStreams, error) {
	var rs ReadRedisStreams
	rs.ConsumerGroup = serviceName
	switch serviceName {
	case "controller":
		rs.ConsumerStreams = []string{
			streams.JobControl.String(),
			streams.ExperimentControl.String(),
			streams.ConfigControl.String()}
	case "executor":
		rs.ConsumerStreams = []string{
			streams.ExperimentControl.String()}
	default:
		return rs, fmt.Errorf("unable to determine redis streams for %s", serviceName)
	}
	return rs, nil
}
