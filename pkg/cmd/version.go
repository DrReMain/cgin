package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Version(v string) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show the version",
		Action: func(c *cli.Context) error {
			fmt.Printf(">>>>>>>>>> Version: %s\n", v)
			return nil
		},
	}
}
