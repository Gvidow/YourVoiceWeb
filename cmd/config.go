package main

import (
	"fmt"

	"go.uber.org/config"
)

type ServerConfig struct {
	port string
	host string
}

type YandexCloudConfig struct {
	oAuthToken string
	folderId   string
}

type ChatGPTConfig struct {
	token string
}

func ReadAllConfig() (*config.YAML, error) {
	optServer := config.File("configs/server/server.yml")
	optYandexCloud := config.File("configs/api/yandex_cloud.yml")
	optGPT := config.File("configs/api/gpt.yml")
	cfg, err := config.NewYAML(optServer, optYandexCloud, optGPT)
	if err != nil {
		return nil, fmt.Errorf("reading all config: %w", err)
	}
	return cfg, nil
}

func GetServerConfig(cfg *config.YAML) ServerConfig {
	return ServerConfig{
		port: cfg.Get("server.port").String(),
		host: cfg.Get("server.host").String(),
	}
}

func GetYandexCloudConfig(cfg *config.YAML) YandexCloudConfig {
	return YandexCloudConfig{
		oAuthToken: cfg.Get("yandexCloud.oAuthToken").String(),
		folderId:   cfg.Get("yandexCloud.folderId").String(),
	}
}

func GetChatGPTConfig(cfg *config.YAML) ChatGPTConfig {
	return ChatGPTConfig{token: cfg.Get("gpt.token").String()}
}
