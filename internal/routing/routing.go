package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/gvidow/YourVoiceWeb/internal/pkg/service"
)

func (r *router) ProduceRouting(s *service.Service) {
	r.Use(gin.Logger(), gin.Recovery())
	r.Static("/static", staticDir)
	r.GET("/main/*id", s.Main)
	chat := r.Group("/chat")
	{
		chat.POST("/swap", s.SwapChats)
		chat.POST("/add", s.AddNewChat)
		chat.POST("/edit", s.EditChat)
		chat.DELETE("/delete/:id", s.DeleteChat)
		chat.POST("/setting/save/:id", s.SaveSettings)
	}
	r.GET("/asr", s.Asr)
	r.NoRoute(s.BadPath)
}
