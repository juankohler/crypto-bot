package bots

import (
	"github.com/juankohler/crypto-bot/common"
)

func Boot(cfg *common.Config, commonDeps *common.Dependencies) error {
	_, err := BuildDependencies(cfg, commonDeps)
	if err != nil {
		return err
	}

	// handlers := NewHandlers(cfg, deps)

	// mux := commonDeps.Mux

	// mux.HandleFunc("POST /v1/init", handlers.Init)

	return nil
}
