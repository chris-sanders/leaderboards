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
	File             string  `json:"-"`
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

func (b *Board) Filter(f *Board) {
	if b.LeaderboardId != f.LeaderboardId {
		return
	}
	filtered := b.Scores
	for i, s := range b.Scores {
		for _, m := range f.Scores {
			if s == m {
				filtered = append(b.Scores[:i], b.Scores[i+1:]...)
			}
		}
	}
	b.Scores = filtered
}

func (b *BoardData) Filter(f *BoardData) {
	for idxb := range b.LeaderboardsData {
		for idxf := range f.LeaderboardsData {
			b.LeaderboardsData[idxb].Filter(&f.LeaderboardsData[idxf])
		}
	}
}

func (b *Board) Add(a *Board) {
	for idxa := range a.Scores {
		exists := false
		for idxb := range b.Scores {
			if a.Scores[idxa] == b.Scores[idxb] {
				exists = true
				continue
			}
		}
		if !exists {
			b.Scores = append(b.Scores, a.Scores[idxa])
		}
	}
}

func main() {
	fmt.Println("Loading file")
	game_data := BoardData{}
	err := game_data.Load("LocalLeaderboards.dat")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Writing file")
	game_data.Save("zarek-db.dat")
}
