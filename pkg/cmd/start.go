package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func Start(callback func(string, string) error) *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "start server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "configs directory",
				Value: "configs",
			},
			&cli.StringFlag{
				Name:  "env",
				Usage: "environment configuration files or directory",
				Value: "dev",
			},
			&cli.BoolFlag{
				Name:  "daemon",
				Usage: "start as daemon",
			},
		},
		Action: func(c *cli.Context) error {
			configs := c.String("config")
			env := c.String("env")

			if c.Bool("daemon") {
				bin, err := filepath.Abs(os.Args[0])
				if err != nil {
					return err
				}

				args := []string{"start"}
				args = append(args, "--config", configs)
				args = append(args, "--env", env)
				fmt.Printf(">>>>>>>>>> Execute Command: %s %s\n", bin, strings.Join(args, " "))

				command := exec.Command(bin, args...)
				err = command.Start()
				if err != nil {
					return err
				}

				pid := command.Process.Pid
				err = os.WriteFile(fmt.Sprintf("%s.lock", c.App.Name), []byte(fmt.Sprintf("%d", pid)), 0666)
				if err != nil {
					return err
				}

				fmt.Printf(">>>>>>>>>> Start Server With Daemon: %d\n", pid)
				os.Exit(0)
			}

			if err := callback(configs, env); err != nil {
				return err
			}

			return nil
		},
	}
}
