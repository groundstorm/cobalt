package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/groundstorm/cobalt/apps/go/src/cobalt/app"
	logging "github.com/op/go-logging"
)

var (
	verboseFlag      = flag.Bool("v", false, "verbose logging")
	fetchRegsCommand = flag.NewFlagSet("fetch-regs", flag.ExitOnError)
	showRegsCommand  = flag.NewFlagSet("show-regs", flag.ExitOnError)
)

func init() {
	flag.Parse()

	logLevel := logging.INFO
	if *verboseFlag {
		logLevel = logging.DEBUG
	}
	logging.SetLevel(logLevel, "cmd")
	logging.SetLevel(logLevel, "app")
	logging.SetLevel(logLevel, "storage")
	logging.SetLevel(logLevel, "smashgg")
}

func usage(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	args := flag.Args()
	if len(args) < 1 {
		usage("not enough arguments")
	}

	app, err := app.New()
	if err != nil {
		log.Fatalf("error starting application: %s", err)
	}
	defer app.Close()

	switch args[0] {
	case "fetch-regs":
		fetchRegsCommand.Parse(args[1:])
		args2 := fetchRegsCommand.Args()
		if len(args2) != 1 {
			usage("fetch-regs [tournament slug]")
		}
		if err := app.FetchRegs(args2[0]); err != nil {
			log.Fatalf("error importing: %s", err)
		}
	case "show-regs":
		showRegsCommand.Parse(args[1:])
		args2 := showRegsCommand.Args()
		if len(args2) != 1 {
			usage("show-regs [tournament slug]")
		}
		a, err := app.LoadRegs(args2[0])
		if err != nil {
			log.Fatalf("error loading: %s", err)
		}
		for _, p := range a.Participants {
			fmt.Printf("%v\n", p)
		}
		fmt.Printf("%d total registrations", len(a.Participants))

	default:
		usage("unknown command")
		return
	}
}
