package main

import (
	"database/sql"
	"log"

	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/merefield/buncheck/app"
	"github.com/merefield/buncheck/migrations"
	bun "github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/migrate"

	grpclog "google.golang.org/grpc/grpclog"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "buncheck",
		Commands: []*cli.Command{
			serverCommand,
			newDBCommand(migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var serverCommand = &cli.Command{
	Name:  "runtest",
	Usage: "start User API server",

	Action: func(c *cli.Context) error {
		_, buncheckapp, err := app.Start(c.Context, "api", c.String("env"))

		cfg := buncheckapp.Cfg

		sqldb, err := sql.Open("pgx", cfg.DB.Dev.PSN)

		_ = bun.NewDB(sqldb, pgdialect.New())

		return err
	},
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Init(ctx, app.DB())
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Migrate(ctx, app.DB())
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Rollback(ctx, app.DB())
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Lock(ctx, app.DB())
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Unlock(ctx, app.DB())
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.CreateGo(ctx, app.DB(), c.Args().Get(0))
				},
			},
			{
				Name:  "create_sql",
				Usage: "create SQL migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.CreateSQL(ctx, app.DB(), c.Args().Get(0))
				},
			},
		},
	}
}

func checkErr(log grpclog.LoggerV2, err error) {
	if err != nil {
		log.Fatal(err)
	}
}
