package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fatih/color"
	"github.com/fujiwara/logutils"
	"github.com/handlename/ssmwrap"
	mackerelnullbridge "github.com/mashiike/mackerel-null-bridge"
	"github.com/urfave/cli/v2"
)

var (
	Version    = "current"
	ssmwrapErr error
)

func main() {
	ssmwrapPaths := os.Getenv("SSMWRAP_PATHS")
	paths := strings.Split(ssmwrapPaths, ",")
	if ssmwrapPaths != "" && len(paths) > 0 {
		ssmwrapErr = ssmwrap.Export(ssmwrap.ExportOptions{
			Paths:   paths,
			Retries: 3,
		})
	}

	cliApp := &cli.App{
		Name:  "mackerel-null-bridge",
		Usage: "A command line tool for filling missing metric values on Mackerel.",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "config file path, can set multiple",
				EnvVars:  []string{"CONFIG_FILE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:        "apikey",
				Aliases:     []string{"k"},
				Usage:       "for access mackerel API",
				DefaultText: "*********",
				EnvVars:     []string{"MACKEREL_APIKEY"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "output log level",
				DefaultText: "info",
				EnvVars:     []string{"LOG_LEVEL"},
			},
			&cli.BoolFlag{
				Name:  "deploy",
				Usage: "deploy flag (cli only)",
			},
			&cli.BoolFlag{
				Name:    "dry-run",
				Usage:   "dry-run flag (lambda only)",
				EnvVars: []string{"DRY_RUN"},
			},
		},
		UsageText: "mackerel-null-bridge --config <config file> --apikey <Mackerel APIKEY>",
		Action: func(c *cli.Context) error {
			if ssmwrapErr != nil {
				return fmt.Errorf("ssmwrap.Export failed: %w", ssmwrapErr)
			}
			cfg := mackerelnullbridge.NewDefaultConfig()
			if err := cfg.Load(c.StringSlice("config")...); err != nil {
				return err
			}
			if err := cfg.ValidateVersion(Version); err != nil {
				return err
			}
			if strings.HasPrefix(os.Getenv("AWS_EXECUTION_ENV"), "AWS_Lambda") || os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
				app := mackerelnullbridge.New(cfg, c.String("apikey"), !c.Bool("dry-run"))
				handler := func(ctx context.Context) error {
					return app.Run(ctx)
				}
				lambda.Start(handler)
				return nil

			}
			app := mackerelnullbridge.New(cfg, c.String("apikey"), c.Bool("deploy"))
			return app.Run(c.Context)
		},
	}
	sort.Sort(cli.FlagsByName(cliApp.Flags))
	cliApp.Version = Version
	cliApp.EnableBashCompletion = true
	cliApp.Before = func(c *cli.Context) error {
		filter := &logutils.LevelFilter{
			Levels: []logutils.LogLevel{"debug", "info", "notice", "warn", "error"},
			ModifierFuncs: []logutils.ModifierFunc{
				nil,
				logutils.Color(color.FgWhite),
				logutils.Color(color.FgHiBlue),
				logutils.Color(color.FgYellow),
				logutils.Color(color.FgRed, color.BgBlack),
			},
			MinLevel: logutils.LogLevel(strings.ToLower(c.String("log-level"))),
			Writer:   os.Stderr,
		}
		log.SetOutput(filter)
		return nil
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	defer cancel()
	if err := cliApp.RunContext(ctx, os.Args); err != nil {
		log.Printf("[error] %s", err)
	}
}
