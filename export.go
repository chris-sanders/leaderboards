package main

import (
	"encoding/json"
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"os"
)

type Score struct {
	Score      int    `json:"_score"`
	PlayerName string `json:"_playerName"`
	FullCombo  bool   `json:"_fullCombo"`
	Timestamp  int    `json:"_timestamp"`
}

type Board struct {
	LeaderboardId string  `json:"_leaderboardId"`
	Scores        []Score `json:"_scores"`
}

type BoardData struct {
	LeaderboardsData []Board `json:"_leaderboardsData"`
	File             string
}

func (b *BoardData) Load(name string) error {
	b.File = name
	jsonFile, err := os.Open(name)
	defer jsonFile.Close()
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(byteValue), b)
	return err
}

func (b *BoardData) Save(name string) error {
	byteSlice, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(name, byteSlice, 0644)
	return err
}

func main() {
	fmt.Println("Loading file")
	game_data := BoardData{}
	err := game_data.Load("LocalLeaderboards.dat")
	if err != nil {
		fmt.Println(err)
	}
	// spew.Config.DisableCapacities = true
	// spew.Config.DisableMethods = true
	// spew.Dump(game_data)

	fmt.Println("Writing file")
	game_data.Save("zarek-db.dat")
}
