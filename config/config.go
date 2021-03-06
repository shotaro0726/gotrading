package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"time"
)

type ConfigList struct {
	ApiKey    string
	ApiSecret string
	LogFile   string

	TradeDuration time.Duration
	Durations     map[string]time.Duration
	DbName        string
	SQLDriver     string
	Port          int
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	durations := map[string]time.Duration {
		"1s": time.Second,
		"1m": time.Minute,
		"1h": time.Hour,
	}

	Config = ConfigList{
		ApiKey:    cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret: cfg.Section("bitflyer").Key("api_secret").String(),
		LogFile:   cfg.Section("bitflyer").Key("log_file").String(),
		Durations: durations,
		TradeDuration: durations[cfg.Section("gotrading").Key("trade_duration").String()],
		DbName: cfg.Section("db").Key("driver").String(),
		SQLDriver: cfg.Section("db").Key("driver").String(),
		Port: cfg.Section("web").Key("port").MustInt(),
	}
}
