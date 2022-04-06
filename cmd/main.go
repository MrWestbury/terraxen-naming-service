package main

import (
	"github.com/MrWestbury/terraxen-naming-service/internals/apis"
	"github.com/MrWestbury/terraxen-naming-service/internals/config"
)

func main() {
	cfg := config.GetConfig()
	api := apis.NewApi(cfg)
	api.Run(":7070")
}
