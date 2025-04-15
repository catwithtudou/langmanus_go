package main

import (
	"catwithtudou/langmanus_go/config"
	"catwithtudou/langmanus_go/internal/llm"
	"catwithtudou/langmanus_go/log"
)

func main() {
	log.InitLogger()
	defer log.GetLogger().Sync()
	config.LoadConfig()
	llm.InitLLMClient()
}
