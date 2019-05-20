package main

import (
	"github.com/spf13/pflag"
	"github.com/vaksi/messaging/cmd"
	"github.com/vaksi/messaging/configs"
	"os"
	"runtime"

	_ "github.com/pressly/goose"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // optimize runtime
	var filename string
	root := cmd.RootCmd()
	fs := pflag.NewFlagSet("Root", pflag.ContinueOnError)
	fs.StringVarP(&filename,
		"file",
		"f",
		"",
		"Custom configuration filename",
	)
	root.Flags().AddFlagSet(fs)
	configuration := configs.New(filename, cmd.ConfigPath...)
	root.AddCommand(
		cmd.NewHttpCmd(
			configuration,
		).BaseCmd,
	)
	root.AddCommand(
		cmd.NewConsumerCmd(
			configuration,
		).BaseCmd,
	)
	if err := root.Execute(); err != nil {
		panic(err.Error())
		os.Exit(1)
	}
}