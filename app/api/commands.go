package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

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
			srv, err := API(logger)

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			go func() {
				for range ch {
					// sig is a ^C, handle it
					logger.Info("Service shutting down..")

					// create context with timeout
					ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
					defer cancel()

					// start http shutdown
					srv.Shutdown(ctx)

					// verify, in worst case call cancel via defer
					select {
					case <-time.After(21 * time.Second):
						logger.Info("Not all connections done")
					case <-ctx.Done():

					}
				}
			}()
			return srv.ListenAndServe()
		},
	}
}
