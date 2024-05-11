package main

import (
	"chipin/cmd/base"
	"fmt"

	_ "chipin/ent/runtime"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	client, _, _, _, _, logger, err2 := base.InitAll()
	if client != nil {
		defer client.Close()
	}
	if err2 != nil {
		logger.Fatal().Msg(fmt.Sprintf("Failed to init: %v", err2))
	}

}
