package bots

import (
	"context"

	"github.com/juankohler/crypto-bot/bots/application"
	"github.com/juankohler/crypto-bot/bots/infrastructure"
	"github.com/juankohler/crypto-bot/common"
	"github.com/juankohler/crypto-bot/libs/go/logs"
)

type Dependencies struct {
}

func BuildDependencies(cfg *common.Config, commonDeps *common.Dependencies) (*Dependencies, error) {
	logs.InitLogger(&logs.Config{
		Format: logs.FormatPretty,
	})

	

	/** Infraestruture dependencies */
	binanceRepo, err := infrastructure.NewBinanceRepo(&cfg.BinanceRepo)
	if err != nil {
		panic(err)
	}

	initService := application.NewInit(binanceRepo)
	initService.Exec(logs.ContextWithLogger(context.Background()), &application.InitInput{})

	return &Dependencies{}, nil
}
