package cmd

import "github.com/urfave/cli/v2"

type Commands []*cli.Command

type Info struct {
	Version  string
	Name     string
	Commands Commands
}

func Run(args []string, info Info) error {
	app := cli.NewApp()
	app.Version = info.Version
	app.Name = info.Name
	app.Commands = info.Commands
	return app.Run(args)
}
