package main

import (
	"context"
	"log"

	"github.com/DimTur/lp_learning_platform/cmd/lp"
)

func main() {
	ctx := context.Background()

	cmd := lp.NewServeCmd()
	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("smth went wrong: %s", err)
	}
}
