package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
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
	//data1 := Data1()
	err := data1.Save("test1.dat")
	if err != nil {
		t.FailNow()
	}
	//data2 := Data2()
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
	board1 = &Data2().LeaderboardsData[0]
	board1.Filter(board2)
	expect = 0
	got = len(board1.Scores)
	if len(board1.Scores) > 0 {
		t.Errorf("Expected %v scores got %v", expect, got)
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
	got = len(data1.LeaderboardsData[0].Scores)
	if got != expect {
		t.Errorf("Expected %v scores got %v", expect, got)
	}
}

func TestBoardAdd(t *testing.T) {
	board1 := &Data1().LeaderboardsData[0]
	board2 := &Data2().LeaderboardsData[0]
	board1.Add(board2)
	fmt.Printf("board1: %v\n", board1)
	fmt.Printf("board2: %v\n", board2)
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
	fmt.Printf("board1: %v\n", board1)
	fmt.Printf("board2: %v\n", board2)
	if !cmp.Equal(board1, board2) {
		t.Error("Boards don't match after adding")
		t.Errorf(cmp.Diff(board1, board2))
	}
}
