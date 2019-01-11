package boardtools

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"testing"
)

var data1 = Data1()
var data2 = Data2()

func Data1() *BoardData {
	data := &BoardData{
		LeaderboardsData: []Board{
			{
				LeaderboardId: "PopStarsEasy",
				Scores: []Score{
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

func Data2() *BoardData {
	data := Data1()
	scores := data.LeaderboardsData[0].Scores
	data.LeaderboardsData[0].Scores = scores[0 : len(scores)-1]
	data.File = "test2.dat"
	return data
}

func TestBoardDataSave(t *testing.T) {
	// Log.SetLevel(Log.TraceLevel)
	fmt.Println("Running tests")
	err := data1.Save("test1.dat")
	if err != nil {
		t.FailNow()
	}
	err = data2.Save("test2.dat")
	if err != nil {
		t.FailNow()
	}
	err = data1.Save("")
	if err == nil {
		t.Error("Open should have returned an error")
	}
}

func TestBoardDataLoad(t *testing.T) {
	test_data1 := &BoardData{}
	err := test_data1.Load("test1.dat")
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(test_data1, data1) {
		t.Error("Loaded data does not match")
		t.Error(cmp.Diff(test_data1, data1))
	}
	test_data2 := &BoardData{}
	err = test_data2.Load("test2.dat")
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(test_data2, data2) {
		t.Error("Loaded data does not match")
		t.Error(cmp.Diff(test_data2, data2))
	}
	err = test_data1.Load("")
	if err == nil {
		t.Error("Load should have returned an error")
	}
}

func TestBoardFilter(t *testing.T) {
	// Filter board1 with board2
	board1 := &Data1().LeaderboardsData[0]
	board2 := &data2.LeaderboardsData[0]
	board1.Filter(board2)
	expect := 1
	got := len(board1.Scores)
	if got != expect {
		t.Errorf("Expected %v scores got %v", expect, got)
	}
	if board1.Scores[0] == board2.Scores[0] {
		t.Errorf("Scores are the same after filtering")
	}
	// Filter board 1 with self
	board1 = &Data1().LeaderboardsData[0]
	board1.Filter(board1)
	expect = 0
	got = len(board1.Scores)
	if len(board1.Scores) > 0 {
		t.Errorf("Expected %v scores got %v", expect, got)
		t.Error(board1)
	}
	// Filter with an empty board
	err := board1.Filter(board1)
	if err != nil {
		t.Error("Empty boards should filter with out error")
		t.Error(err)
	}
	err = board1.Filter(&Board{})
	if err == nil {
		t.Error("Boards should not filter if Id does not match")
	}
}

func TestBoardDataFilter(t *testing.T) {
	data1 := Data1()
	data1.Filter(data2)
	expect := 1
	got := len(data1.LeaderboardsData[0].Scores)
	if got != expect {
		t.Errorf("Expected %v scores got %v", expect, got)
	}
	data1.Filter(data2)
	expect = 1
	got = len(data1.LeaderboardsData[0].Scores)
	if got != expect {
		t.Errorf("Expected %v scores got %v", expect, got)
	}
	data1.Filter(data1)
	expect = 0
	got = len(data1.LeaderboardsData)
	if got != expect {
		t.Errorf("Expected %v scores got %v", expect, got)
	}
}

func TestBoardAdd(t *testing.T) {
	board1 := &Data1().LeaderboardsData[0]
	board2 := &Data2().LeaderboardsData[0]
	board1.Add(board2)
	expect := 2
	got := len(board1.Scores)
	if expect != got {
		t.Error("Added duplicate score")
	}
	board2.Add(board1)
	got = len(board2.Scores)
	if expect != got {
		t.Error("Didn't add missing score")
	}
	if !cmp.Equal(board1, board2) {
		t.Error("Boards don't match after adding")
		t.Errorf(cmp.Diff(board1, board2))
	}
	board3 := &Board{}
	err := board3.Add(board1)
	if err == nil {
		t.Error("Boards should not add if Id does not match")
	}
	board3.LeaderboardId = board1.LeaderboardId
	err = board3.Add(board1)
	if err != nil {
		t.Error("Empty board does not add")
	}
	if !cmp.Equal(board3, board1) {
		t.Error("Boards don't match after adding")
		t.Errorf(cmp.Diff(board3, board1))
	}
	files, err := filepath.Glob("*.dat")
	if err != nil {
		t.Errorf("Error listing dat files: %v", err)
	}
	for _, f := range files {
		if err = os.Remove(f); err != nil {
			t.Errorf("Error removing dat files: %v", err)
		}
	}
}

func TestBoardDataAdd(t *testing.T) {
	// data1 = &Data1()
	data2 = &BoardData{}
	data2.Add(data1)
	if !cmp.Equal(data1.LeaderboardsData, data2.LeaderboardsData) {
		t.Errorf("BoardData add didn't add new data to empty board:\n%v",
			cmp.Diff(data1.LeaderboardsData,
				data2.LeaderboardsData))
	}
}

func TestBoardFilterScores(t *testing.T) {
	board1 := &Data1().LeaderboardsData[0]
	board1.FilterScores(3)
	expect := 2
	got := len(board1.Scores)
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, board1.Scores)
	}
	board1.FilterScores(2)
	expect = 2
	got = len(board1.Scores)
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, board1.Scores)
	}
	board1.FilterScores(1)
	expect = 1
	got = len(board1.Scores)
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, board1.Scores)
	}
	board1.FilterScores(0)
	expect = 0
	got = len(board1.Scores)
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, board1.Scores)
	}
}

func TestBoardDataFilterScores(t *testing.T) {
	data1 := Data1()
	data1.LeaderboardsData = append(data1.LeaderboardsData, data1.LeaderboardsData[0])
	data1.FilterScores(3)
	expect := 2
	got := len(data1.LeaderboardsData[0].Scores)
	if len(data1.LeaderboardsData[0].Scores) != len(data1.LeaderboardsData[1].Scores) {
		t.Errorf("Boards not filtered the same:\n%v\n%v", data1.LeaderboardsData[0].Scores, data1.LeaderboardsData[1].Scores)
	}
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, data1)
	}
	data1.FilterScores(2)
	expect = 2
	got = len(data1.LeaderboardsData[0].Scores)
	if len(data1.LeaderboardsData[0].Scores) != len(data1.LeaderboardsData[1].Scores) {
		t.Errorf("Boards not filtered the same:\n%v\n%v", data1.LeaderboardsData[0].Scores, data1.LeaderboardsData[1].Scores)
	}
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, data1)
	}
	data1.FilterScores(1)
	expect = 1
	got = len(data1.LeaderboardsData[0].Scores)
	if len(data1.LeaderboardsData[0].Scores) != len(data1.LeaderboardsData[1].Scores) {
		t.Errorf("Boards not filtered the same:\n%v\n%v", data1.LeaderboardsData[0].Scores, data1.LeaderboardsData[1].Scores)
	}
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, data1)
	}
	data1.FilterScores(0)
	expect = 0
	got = len(data1.LeaderboardsData)
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, data1)
	}
}

func TestBoardTruncateScores(t *testing.T) {
	board1 := &Data1().LeaderboardsData[0]
	board1.TruncateScores(3)
	expect := 2
	got := len(board1.Scores)
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, board1.Scores)
	}
	board1.TruncateScores(2)
	expect = 2
	got = len(board1.Scores)
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, board1.Scores)
	}
	board1.TruncateScores(1)
	expect = 1
	got = len(board1.Scores)
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, board1.Scores)
	}
	board1.TruncateScores(0)
	expect = 0
	got = len(board1.Scores)
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, board1.Scores)
	}
}

func TestBoardDataTruncateScores(t *testing.T) {
	data1 := Data1()
	data1.LeaderboardsData = append(data1.LeaderboardsData, data1.LeaderboardsData[0])
	data1.TruncateScores(3)
	expect := 2
	got := len(data1.LeaderboardsData[0].Scores)
	if len(data1.LeaderboardsData[0].Scores) != len(data1.LeaderboardsData[1].Scores) {
		t.Errorf("Boards not filtered the same:\n%v\n%v", data1.LeaderboardsData[0].Scores, data1.LeaderboardsData[1].Scores)
	}
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, data1)
	}
	data1.TruncateScores(2)
	expect = 2
	got = len(data1.LeaderboardsData[0].Scores)
	if len(data1.LeaderboardsData[0].Scores) != len(data1.LeaderboardsData[1].Scores) {
		t.Errorf("Boards not filtered the same:\n%v\n%v", data1.LeaderboardsData[0].Scores, data1.LeaderboardsData[1].Scores)
	}
	if expect != got {
		t.Errorf("Filtered to many scores, expected %v got %v\n%v", expect, got, data1)
	}
	data1.TruncateScores(1)
	expect = 1
	got = len(data1.LeaderboardsData[0].Scores)
	if len(data1.LeaderboardsData[0].Scores) != len(data1.LeaderboardsData[1].Scores) {
		t.Errorf("Boards not filtered the same:\n%v\n%v", data1.LeaderboardsData[0].Scores, data1.LeaderboardsData[1].Scores)
	}
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, data1)
	}
	data1.TruncateScores(0)
	expect = 0
	got = len(data1.LeaderboardsData[0].Scores)
	if len(data1.LeaderboardsData[0].Scores) != len(data1.LeaderboardsData[1].Scores) {
		t.Errorf("Boards not filtered the same:\n%v\n%v", data1.LeaderboardsData[0].Scores, data1.LeaderboardsData[1].Scores)
	}
	if expect != got {
		t.Errorf("Scores not filtered correctly, expected %v got %v\n%v", expect, got, data1)
	}
}
