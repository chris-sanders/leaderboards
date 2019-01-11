package main

import (
	"fmt"
	"github.com/chris-sanders/leaderboards/internal/cfg"
	"github.com/chris-sanders/leaderboards/internal/cmds"
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

	cmds.UpdateLocalDb(config)
	sftp.InitClient(config)
	fmt.Println("Connected to server")
	defer sftp.CloseClient()
	sftp.UploadDb()
	fmt.Println("File uploaded")
}
