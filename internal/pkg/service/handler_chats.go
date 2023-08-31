package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gvidow/YourVoiceWeb/internal/pkg/repository/chat"
)

func (s *Service) Main(c *gin.Context) {
	chatRepo, err := s.usecase.ChatRepository()
	if err != nil {
		s.log.Panic(err.Error())
	}

	chatList, err := chatRepo.SelectAllByOrder(c.Request.Context())
	if err != nil {
		s.log.Error(err.Error())
	}

	settings := &chat.Setting{}
	var activeID string
	if len(chatList) > 0 {
		activeID = getStringParam(c.Param("id"))
		var find bool
		for _, ch := range chatList {
			if ch.StringID() == activeID {
				find = true
				settings = ch.Settings
				break
			}
		}
		if !find {
			activeID = chatList[0].StringID()
			settings = chatList[0].Settings
		}
	}

	c.Render(http.StatusOK, render.HTML{
		Template: s.tmpl,
		Name:     "index.html",
		Data: map[string]any{
			"Chats":    chatList,
			"Active":   activeID,
			"Settings": settings,
		},
	})
}

func (s *Service) DeleteChat(c *gin.Context) {
	chatRepo, err := s.usecase.ChatRepository()
	if err != nil {
		s.log.Panic(err.Error())
	}
	id := getStringParam(c.Param("id"))
	n, err := chatRepo.DeleteMany(c.Request.Context(), []string{id})
	if err != nil {
		s.log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"code":    "delete_chat",
			"message": "an error occurred when deleting the chat",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"comment": fmt.Sprintf("%d chats were deleted", n),
		})
	}
}

func (s *Service) SaveSettings(c *gin.Context) {
	chatRepo, err := s.usecase.ChatRepository()
	if err != nil {
		s.log.Panic(err.Error())
	}
	id := getStringParam(c.Param("id"))

	settings, err := readChatSettings(c.Request.Body)
	if err != nil {
		s.log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"code":    "bad_request_body",
			"message": "couldn't subtract the necessary settings from the request body",
		})
		return
	}
	c.Request.Body.Close()

	err = chatRepo.SaveSettings(c.Request.Context(), id, settings)
	if err != nil {
		s.log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"code":    "save_settings_chat",
			"message": "couldn't save chat settings",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"comment": "settings saved",
		})
	}
}

func (s *Service) SwapChats(c *gin.Context) {
	chatRepo, err := s.usecase.ChatRepository()
	if err != nil {
		s.log.Panic(err.Error())
	}
	m, err := unmarshalBody(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, responseErr("bad_request_body",
			"couldn't subtract the necessary settings from the request body"))
		return
	}
	c.Request.Body.Close()
	id1, ok1 := m["id1"]
	id2, ok2 := m["id2"]
	if !(ok1 && ok2) {
		c.JSON(http.StatusOK, responseErr("bad_request_body",
			"couldn't subtract the necessary settings from the request body"))
		return
	}
	err = chatRepo.SwapPlaces(c.Request.Context(), id1, id2)
	if err != nil {
		s.log.Error(err.Error())
		c.JSON(http.StatusOK, responseErr("swap_chat",
			"couldn't swap selected chats"))
	} else {
		s.log.Sugar().Infof("successfully swapped chats with id %s and %s", id1, id2)
		c.JSON(http.StatusOK, responseOK("successfully swapped chats"))
	}
}

func (s *Service) AddNewChat(c *gin.Context) {
	chatRepo, err := s.usecase.ChatRepository()
	if err != nil {
		s.log.Panic(err.Error())
	}
	m, err := unmarshalBody(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, responseErr("bad_request_body",
			"couldn't subtract the necessary settings from the request body"))
		return
	}
	c.Request.Body.Close()
	title, ok := m["title"]
	if !ok {
		c.JSON(http.StatusOK, responseErr("bad_request_body",
			"couldn't subtract the necessary settings from the request body"))
		return
	}
	id, err := chatRepo.AddNewChat(c.Request.Context(), title)
	if err != nil {
		s.log.Error(err.Error())
		c.JSON(http.StatusOK, responseErr("add_chat", "fail"))
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"comment": "successfully inserted a new entry",
			"id":      id,
		})
	}
}

func (s *Service) EditChat(c *gin.Context) {
	chatRepo, err := s.usecase.ChatRepository()
	if err != nil {
		s.log.Panic(err.Error())
	}
	m, err := unmarshalBody(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, responseErr("bad_request_body",
			"couldn't subtract the necessary settings from the request body"))
		return
	}
	c.Request.Body.Close()
	title, ok1 := m["title"]
	id, ok2 := m["id"]
	if !(ok1 && ok2) {
		c.JSON(http.StatusOK, responseErr("bad_request_body",
			"couldn't subtract the necessary settings from the request body"))
		return
	}
	err = chatRepo.EditChat(c.Request.Context(), id, title)
	if err != nil {
		s.log.Error(err.Error())
		c.JSON(http.StatusOK, responseErr("edit_chat", "fdsf"))
	} else {
		c.JSON(http.StatusOK, responseOK("fdsf"))
	}
}
