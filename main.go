package main

import (
	_ "study-rabbitmq/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"study-rabbitmq/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
