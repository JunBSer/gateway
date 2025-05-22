package main

import (
	"github.com/JunBSer/gateway/internal/app"
	"github.com/JunBSer/gateway/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.MustRun(cfg)
}
