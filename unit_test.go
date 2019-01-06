package main

import (
	//	"fmt"
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
