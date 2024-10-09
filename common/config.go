package common

import (
	"github.com/juankohler/crypto-bot/libs/go/config"
	"github.com/juankohler/crypto-bot/libs/go/restclient"
)

type Config struct {
	Env         config.Env
	Port        int
	Database    string
	BinanceRepo restclient.Config
}

func GetConfig() (*Config, error) {
	env, err := config.NewEnv(config.GetEnv("ENV", "DEV"))
	if err != nil {
		return nil, err
	}

	timeOut := 9000

	return &Config{
		Env:      env,
		Port:     config.GetEnvAsInt("PORT", 8080),
		Database: config.GetEnv("DATABASE", "database/local.db"),
		BinanceRepo: restclient.Config{
			BaseUrl:   "https://api.binance.com/api",
			Retries:   1,
			TimeoutMs: &timeOut,
		},
	}, nil
}
