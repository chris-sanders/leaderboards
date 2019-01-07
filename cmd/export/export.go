package main

import (
	"fmt"
	bt "github.com/chris-sanders/leaderboards/boardtools"
	"github.com/chris-sanders/leaderboards/internal/cfg"
)

func main() {
	fmt.Println("Loading file")
	game_data := bt.BoardData{}
	err := game_data.Load("LocalLeaderboards.dat")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Writing file")
	game_data.Save("zarek-db.dat")
}
