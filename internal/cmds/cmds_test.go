package cmds

import (
	"fmt"
	bt "github.com/chris-sanders/leaderboards/boardtools"
	"github.com/chris-sanders/leaderboards/internal/cfg"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func Data1() *bt.BoardData {
	data := &bt.BoardData{
		LeaderboardsData: []bt.Board{
			{
				LeaderboardId: "PopStarsEasy",
				Scores: []bt.Score{
					{
						Score:      888888,
						PlayerName: "TST",
						Timestamp:  1545962334,
						FullCombo:  false,
					},
					{
						Score:      999,
						PlayerName: "TST",
						Timestamp:  1545962334,
						FullCombo:  false,
					},
				},
			},
		},
		File: "test1.dat",
	}
	return data
}

func Data2() *bt.BoardData {
	data2 := Data1()
	data2.LeaderboardsData[0].Scores = data2.LeaderboardsData[0].Scores[0:1]
	return data2
}

func Data3() *bt.BoardData {
	data := &bt.BoardData{
		LeaderboardsData: []bt.Board{
			{
				LeaderboardId: "PopStarsEasy",
				Scores: []bt.Score{
					{
						Score:      42,
						PlayerName: "TST",
						Timestamp:  1545962334,
						FullCombo:  false,
					},
					{
						Score:      1,
						PlayerName: "TST",
						Timestamp:  1545962334,
						FullCombo:  false,
					},
				},
			},
		},
		File: "test1.dat",
	}
	return data
}

func TestUpdateLocalDb(t *testing.T) {
	// Turn off logging
	log.SetLevel(log.PanicLevel)
	Stdout := os.Stdout
	os.Stdout = nil
	// Create config
	config, _ := (&cfg.Config{}).New()
	config, _ = config.New()
	// Write out a Game Data file
	data1 := Data1()
	data1.Save(config.Global.Local_file)
	// With no local database it should be created
	UpdateLocalDb(config)
	local_info, err := os.Stat(config.Global.Local_file)
	db_path := fmt.Sprintf("%v-db.dat", config.Global.Account)
	db_info, err := os.Stat(db_path)
	if err != nil {
		t.Error("Failed to create new db file:", err.Error())
	}
	if local_info.Size() != db_info.Size() {
		t.Errorf("Db size %v does not match local size %v", local_info.Size(), db_info.Size())
	}
	// Database should not be modified on 2nd run
	UpdateLocalDb(config)
	db_info2, err := os.Stat(db_path)
	if err != nil {
		t.Error("Failed to stat db file:", err.Error())
	}
	if db_info.ModTime() != db_info2.ModTime() {
		t.Error("Db modified during no-op")
	}
	// Database should not be modified for known values
	data3 := Data3()
	data3.Save(config.Global.Local_file)
	data3.Save("other-db.dat")
	UpdateLocalDb(config)
	db_info3, err := os.Stat(db_path)
	if err != nil {
		t.Error("Failed to stat db file:", err.Error())
	}
	if db_info.ModTime() != db_info3.ModTime() {
		t.Error("Db modified during no-op")
	}
	// New values should be added while known values are skipped
	// Local: 88888, 999
	// db: 42, 1
	// other: 88888
	// expected result: 42, 1, 999
	data3.Save(db_path)
	data1.Save(config.Global.Local_file)
	data2 := Data2()
	data2.Save("other-db.dat")
	UpdateLocalDb(config)
	test_data := bt.BoardData{}
	test_data.Load(db_path)
	os.Remove("LocalLeaderboards.dat")
	os.Remove("other-db.dat")
	os.Remove("unique-name-db.dat")
	expect := []int{42, 1, 999}
	got := []int{test_data.LeaderboardsData[0].Scores[0].Score,
		test_data.LeaderboardsData[0].Scores[1].Score,
		test_data.LeaderboardsData[0].Scores[2].Score}
	for i, _ := range expect {
		if expect[i] != got[i] {
			t.Errorf("test_data %v doesn't match %v", got, expect)
		}
	}
	// Enable standout
	os.Stdout = Stdout
}
