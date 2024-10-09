package main

import (
	"fmt"
	"net/http"

	"github.com/juankohler/crypto-bot/bots"
	"github.com/juankohler/crypto-bot/common"
	"github.com/juankohler/crypto-bot/libs/go/http/server"
)

var bootables = []common.Bootable{
	bots.Boot,
}

func boot(cfg *common.Config, deps *common.Dependencies) error {
	for _, bootable := range bootables {
		if err := bootable(cfg, deps); err != nil {
			return err
		}
	}

	fmt.Printf("Server running on port %d\n", cfg.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), deps.Mux)
}

func main() {
	cfg, err := common.GetConfig()
	if err != nil {
		panic(err)
	}

	deps, err := common.BuildDependencies(cfg)
	if err != nil {
		panic(err)
	}

	server.EnableLogging()

	if err := boot(cfg, deps); err != nil {
		panic(err)
	}
}
