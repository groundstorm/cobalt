package main

import (
	"fmt"

	"github.com/groundstorm/cobalt/apps/go/src/smashgg"
	logging "github.com/op/go-logging"
)

func main() {
	logging.SetLevel(logging.INFO, "smashgg")
	tourney, err := smashgg.ImportTournament("evo2018")
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Printf("%v\n", tourney)
}
