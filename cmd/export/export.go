package main

import (
	"fmt"
	bt "github.com/chris-sanders/leaderboards/boardtools"
	"github.com/chris-sanders/leaderboards/internal/cfg"
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

	fmt.Println("Loading Local Game Data")
	log.Info("Loading Local Game Data")
	game_data := &bt.BoardData{}
	err = game_data.Load(config.Global.Local_file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Writing database file")
	log.Info("Writing database file")
	data_file := fmt.Sprintf("%v-db.dat", config.Global.Account)
	db_data := &bt.BoardData{}
	err = db_data.Load(data_file)
	if err != nil {
		log.Warnf("Error loading db file: %v \n", err)
	}
	game_data.Filter(db_data)
	new_scores := game_data.Marshal()
	if len(new_scores) > 0 {
		fmt.Printf("Adding new scores:\n%v \n", new_scores)
		log.Infof("Adding new scores:\n%v \n", new_scores)
	}
	db_data.Add(game_data)
	db_data.Save(db_data.File)
}
