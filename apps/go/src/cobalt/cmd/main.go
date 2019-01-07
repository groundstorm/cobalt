package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/groundstorm/cobalt/apps/go/src/cobalt/app"
	"github.com/groundstorm/cobalt/apps/go/src/smashgg"
	logging "github.com/op/go-logging"
)

var (
	log         = logging.MustGetLogger("cmd")
	verboseFlag = flag.Bool("v", false, "verbose logging")

	fetchRegsCommand    = flag.NewFlagSet("fetch-regs", flag.ExitOnError)
	fetchRegsStdoutFlag = fetchRegsCommand.Bool("stdout", false, "write to stdout rather than the db")

	showRegsCommand = flag.NewFlagSet("show-regs", flag.ExitOnError)
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
		fetchRegs(app, args[1:])
	case "show-regs":
		showRegs(app, args[1:])
	default:
		usage("unknown command")
	}
}

func fetchRegs(app *app.App, args []string) {
	fetchRegsCommand.Parse(args)
	args2 := fetchRegsCommand.Args()
	if len(args2) != 1 {
		usage("fetch-regs [tournament slug]")
	}

	slug := args2[0]
	log.Infof("fetching registrations for %s", slug)

	info, err := smashgg.GetTournamentRegistrationInfo(slug)
	if err != nil {
		log.Fatal("failed to get tournament info: %s", err)
	}
	if *fetchRegsStdoutFlag {
		blob, err := smashgg.FetchAttendees(info)
		if err != nil {
			log.Fatalf("failed to get attendee list: %v", err)
		}
		fmt.Println(string(blob))
		return
	}
	attendees, err := smashgg.LoadAttendees(info)
	if err != nil {
		log.Fatalf("failed to fetch attendee list: %v", err)
	}

	log.Infof("storing %d participants for %s", len(attendees.Registrations), slug)
	if err := app.StoreRegs(slug, attendees); err != nil {
		log.Fatalf("error importing: %s", err)
	}
}
func showRegs(app *app.App, args []string) {
	showRegsCommand.Parse(args)
	args2 := showRegsCommand.Args()
	if len(args2) != 1 {
		usage("show-regs [tournament slug]")
	}
	a, err := app.LoadRegs(args2[0])
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}
	for _, r := range a.Registrations {
		fmt.Printf("%v\n", r)
	}
	fmt.Printf("%d total registrations", len(a.Registrations))
}
