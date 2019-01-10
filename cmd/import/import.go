package main

import (
	"fmt"
	bt "github.com/chris-sanders/leaderboards/boardtools"
	"github.com/chris-sanders/leaderboards/internal/cfg"
	"github.com/chris-sanders/leaderboards/internal/cmds"
	"github.com/chris-sanders/leaderboards/internal/sftp"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Setup logging
	logFile, err := os.OpenFile("leaderboards.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Get config
	config, err := (&cfg.Config{}).New()
	if err != nil {
		fmt.Printf("Error loading config: %v \n", err)
		log.Fatalf("Error loading config: %v \n", err)
	}

	// Download remote files
	fmt.Println("Downlading remote files")
	log.Info("Downloading remote files")
	sftp.InitClient(config)
	defer sftp.CloseClient()
	sftp.GetRemoteFiles()

	// Merge with local game file
	cmds.UpdateLocalDb(config)
	merge_data := &bt.BoardData{}
	paths, err := filepath.Glob("*-db.dat")
	for _, file := range paths {
		fmt.Printf("Processing scores in %v\n", file)
		log.Infof("Processing scores in %v", file)
		if strings.HasPrefix(file, config.Global.Account) {
			local_data := &bt.BoardData{}
			err := local_data.Load(file)
			if err != nil {
				fmt.Printf("Error loading file: %v", err)
				log.Panic("Error loading file: %v", err)
			}
			local_data.FilterScores(config.Import.Local_limit)
			merge_data.Add(local_data)
		} else {
			remote_data := &bt.BoardData{}
			err := remote_data.Load(file)
			if err != nil {
				fmt.Printf("Error loading file: %v\n", err)
				log.Warnf("Error loading file: %v", err)
			}
			remote_data.FilterScores(config.Import.Remote_limit)
			merge_data.Add(remote_data)
		}
	}
	merge_data.TruncateScores(10)
	fmt.Println("Writing new game data file")
	log.Info("Writing new game data file")
	err = merge_data.Save(config.Global.Local_file)
	if err != nil {
		fmt.Printf("Error writing file: %v", err)
		log.Panic("Error writing file: %v", err)
	}
}
