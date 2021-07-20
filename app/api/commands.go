package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/foundation/logs"
)

func APICommand() *cli.Command {
	return &cli.Command{
		Name:  "api",
		Usage: "start api server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log",
				Value: "prod",
				Usage: "set log level (prod, debug)",
				EnvVars: []string{"LOG_LEVEL"},
			},
		},
		Action: func(c *cli.Context) error {
			// setup logger
			var logger *zap.Logger
			var err error
			switch c.String("log") {
			case "debug":
				logger, err = logs.CreateLoggerDebug()
				if err != nil {
					return err
				}
			case "prod":
				logger, err = logs.CreateLogger()
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("unknown log level %q", c.String("log"))
			}
			defer logger.Sync()

			logger.Info("service is starting", zap.String("app-name", c.App.Name), zap.String("app-version", c.App.Version))
			return RunAPIServer(logger)
		},
	}
}
