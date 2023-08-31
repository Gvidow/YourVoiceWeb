package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) BadPath(c *gin.Context) {
	s.log.Info("request on bad url path")
	c.Redirect(http.StatusPermanentRedirect, "/main")
}
