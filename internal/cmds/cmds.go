package cmds

import (
	"fmt"
	bt "github.com/chris-sanders/leaderboards/boardtools"
	"github.com/chris-sanders/leaderboards/internal/cfg"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
)

func UpdateLocalDb(config *cfg.Config) {
	fmt.Println("Loading Local Game Data")
	log.Info("Loading Local Game Data")
	game_data := &bt.BoardData{}
	err := game_data.Load(config.Global.Local_file)
	if err != nil {
		fmt.Println(err)
		log.Panic(err)
	}

	data_file := fmt.Sprintf("%v-db.dat", config.Global.Account)
	db_data := &bt.BoardData{}
	paths, err := filepath.Glob("*-db.dat")
	for _, file := range paths {
		log.Infof("Filtering known scores from %v", file)
		data := &bt.BoardData{}
		err := data.Load(file)
		if err != nil {
			fmt.Printf("Error loading file: %v", err)
			log.Panic("Error loading file: %v", err)
		}
		game_data.Filter(data)
		if strings.HasPrefix(file, config.Global.Account) {
			db_data = data
		}
	}

	new_scores := game_data.Marshal()
	if len(new_scores) > 0 {
		fmt.Printf("Adding new scores:\n %v \n", new_scores)
		log.Infof("Adding new scores:\n %v \n", new_scores)
		fmt.Println("Writing database file")
		log.Info("Writing database file")
		db_data.Add(game_data)
		db_data.Save(data_file)
	} else {
		fmt.Println("No new local data")
		log.Info("No new local data")
	}

}

func WriteGameFile(config *cfg.Config) {
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
