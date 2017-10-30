package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	getopt "github.com/kesselborn/go-getopt"
	"gohive/internal/config"
	"gohive/server/game"
)

var (
	_VERSION_   = "unknown"
	_BUILDDATE_ = "unknown"

	pwd string
)

func init() {
	pwd, _ = os.Getwd()
}

func main() {
	optionDefinition := getopt.Options{
		"description",
		getopt.Definitions{
			{"config|c", "config file", getopt.IsConfigFile | getopt.ExampleIsDefault, filepath.Join(config.PWD, "conf/gamed.json")},
			{"version|v", "show version", getopt.Optional | getopt.Flag, nil},
		},
	}

	options, _, _, e := optionDefinition.ParseCommandLine()
	help, wantsHelp := options["help"]
	if e != nil || wantsHelp {
		exit_code := 0
		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Println("**** Error: ", e.Error(), "\n", optionDefinition.Help())
			exit_code = e.ErrorCode
		}
		os.Exit(exit_code)
	}
	version, showVersion := options["version"]
	if showVersion && version.Bool {
		fmt.Printf("server version %s\n%s\n", _VERSION_, _BUILDDATE_)
		os.Exit(0)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	cfg, err := config.Load(options["config"].String)
	if err != nil {
		panic(fmt.Sprintf("load config failed: %s", err))
	}
	game.Run(sc, cfg)
}
