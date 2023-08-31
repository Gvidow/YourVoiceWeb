package routing

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gvidow/YourVoiceWeb/logger"
	"go.uber.org/config"
)

const (
	configName = "server"
	staticDir  = "static"
)

type router struct {
	*gin.Engine
	cfg config.Value
	log *logger.Logger
}

func New(cfg *config.YAML, log *logger.Logger) router {
	val := cfg.Get(configName)
	gin.SetMode(val.Get("mode").String())
	return router{gin.New(), val, log}
}

func (r *router) Run() error {
	addr := fmt.Sprintf("%s:%s", r.cfg.Get("host"), r.cfg.Get("port"))
	r.log.Sugar().Infof("start server on %s", addr)
	return r.Engine.Run(addr)
}
