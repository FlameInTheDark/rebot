package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/foundation/logs"
)

// RunCommand is a "run" command for the command line
func RunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "start weather command service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log",
				Value:   "prod",
				Usage:   "set log level (prod, debug)",
				EnvVars: []string{"LOG_LEVEL"},
			},
		},
		Action: func(c *cli.Context) error {
			// setup logger
			var logger *zap.Logger
			switch c.String("log") {
			case "debug":
				logger = logs.CreateLoggerDebug()
			case "prod":
				logger = logs.CreateLogger()
			default:
				return fmt.Errorf("unknown log level %q", c.String("log"))
			}
			defer func(logger *zap.Logger) {
				err := logger.Sync()
				if err != nil {
				}
			}(logger)

			logger.Info("service is starting", zap.String("app-name", c.App.Name), zap.String("app-version", c.App.Version))
			return RunWeatherService(logger)
		},
	}
}
