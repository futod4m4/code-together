package main

import (
	"github.com/futod4m4/m/config"
	"github.com/futod4m4/m/internal/server"
	"github.com/futod4m4/m/pkg/logger"
	"github.com/futod4m4/m/pkg/utils"
	"log"
	"os"
)

func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("Parse Config %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	s := server.NewServer(cfg, appLogger)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
