package main

import (
	"fmt"
	// bt "github.com/chris-sanders/leaderboards/boardtools"
	"github.com/chris-sanders/leaderboards/internal/cfg"
	"github.com/chris-sanders/leaderboards/internal/sftp"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	logFile, err := os.OpenFile("leaderboards.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	config, err := (&cfg.Config{}).New()
	if err != nil {
		fmt.Printf("Error loading config: %v \n", err)
		log.Fatalf("Error loading config: %v \n", err)
	}

	fmt.Println("Downlading remote files")
	log.Info("Downloading remote files")

	sftp.InitClient(config)
	defer sftp.CloseClient()
	sftp.GetRemoteFiles()
}
