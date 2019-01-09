package boardtools

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
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
		return err
	}
	err = ioutil.WriteFile(name, byteSlice, 0644)
	return err
}

func scoreInSlice(s Score, slice []Score) bool {
	for _, m := range slice {
		if s == m {
			log.Tracef("Score %v is in slice", s)
			return true
		}
	}
	return false
}

func (b *Board) Filter(f *Board) error {
	if b.LeaderboardId != f.LeaderboardId {
		return errors.New("Board Ids do not match")
	}
	log.Tracef("Filtering %v with filter %v", b.LeaderboardId, f.LeaderboardId)
	filtered := b.Scores[:0]
	for _, s := range b.Scores {
		if !scoreInSlice(s, f.Scores) {
			filtered = append(filtered, s)
		}
	}
	log.Tracef("Filtered slice: %v", filtered)
	b.Scores = filtered
	return nil
}

func (b *BoardData) Filter(f *BoardData) {
	for idxb := range b.LeaderboardsData {
		for idxf := range f.LeaderboardsData {
			b.LeaderboardsData[idxb].Filter(&f.LeaderboardsData[idxf])

		}
	}
}

func (b *Board) Add(a *Board) error {
	if b.LeaderboardId != a.LeaderboardId {
		err := errors.New("Board Ids do not match")
		return err
	}
	log.Tracef("Adding %v with  %v", b.LeaderboardId, a.LeaderboardId)
	for idxa := range a.Scores {
		exists := false
		for idxb := range b.Scores {
			if a.Scores[idxa] == b.Scores[idxb] {
				exists = true
				log.Tracef("Not adding score %v duplicate %v", a.Scores[idxa], b.Scores[idxb])
				continue
			}
		}
		if !exists {
			log.Tracef("Adding new score %v", a.Scores[idxa])
			b.Scores = append(b.Scores, a.Scores[idxa])
		}
	}
	return nil
}

func (b *BoardData) Add(a *BoardData) {
	for idxa := range a.LeaderboardsData {
		boardFound := false
		for idxb := range b.LeaderboardsData {
			err := b.LeaderboardsData[idxb].Add(&a.LeaderboardsData[idxa])
			if err == nil {
				boardFound = true
			}
		}
		if !boardFound {
			log.Tracef("Adding new board %v", a.LeaderboardsData[idxa])
			b.LeaderboardsData = append(b.LeaderboardsData, a.LeaderboardsData[idxa])
		}
	}
}

func (b *BoardData) Marshal() string {
	var bld strings.Builder
	for _, board := range b.LeaderboardsData {
		if len(board.Scores) > 0 {
			// byteArray, err := json.MarshalIndent(board.Scores, board.LeaderboardId, " ")
			/*
				byteArray, err := json.Marshal(board.Scores)
				if err != nil {
					log.Error("Unalbe to Marshal board: %v", err)
				}
				bld.Write(byteArray)
				bld.WriteString("\n")
			*/
			values := fmt.Sprintf("%v: \n\t", board.LeaderboardId)
			bld.WriteString(values)
			for _, score := range board.Scores {
				score_info := fmt.Sprintf("%v:%v, ", score.PlayerName, score.Score)
				bld.WriteString(score_info)
			}
			bld.WriteString("\n")
		}
	}
	return bld.String()
}
