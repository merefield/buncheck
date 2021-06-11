package main

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	model "github.com/merefield/buncheck/model"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/merefield/buncheck/app"
	"github.com/merefield/buncheck/migrations"
	bun "github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/migrate"

	grpclog "google.golang.org/grpc/grpclog"

	"github.com/uptrace/bun/extra/bundebug"
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

		//cfg.PreferSimpleProtocol = true

		sqldb, err := sql.Open("pgx", cfg.DB.Dev.PSN)

		db := bun.NewDB(sqldb, pgdialect.New())

		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose()))

		// if err := createSchema(c, db); err != nil {
		// 	panic(err)
		// }

		users := []*model.User{
			{Username: "user 1", Email: "bob"},
			{Username: "user 2", Email: "jim"},
		}
		if _, err := db.NewInsert().Model(&users).Exec(c.Context); err != nil {
			return err
		}

		user := new(model.User)
		if err := db.NewSelect().
			Model(user).
			//Column("user.*").
			// Relation("OwnerOfGroups", func(q *bun.SelectQuery) *bun.SelectQuery {
			// 	return q.Where("active IS TRUE")
			// }).
			// OrderExpr("user.id ASC").
			// Limit(1).
			Scan(c.Context); err != nil {
			panic(err)
		}
		fmt.Println(user.ID, user.Username) //, user.OwnerOfGroups[0])
		// Output: 1 user 1 &{1 en true 1} &{2 ru true 1}

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
			{
				Name:  "fixtures",
				Usage: "load fixtures",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					cfg := app.Cfg

					sqldb, err := sql.Open("pgx", cfg.DB.Dev.PSN)
					if err != nil {
						return err
					}

					db := bun.NewDB(sqldb, pgdialect.New())

					// Let the db know about the models.
					models := []interface{}{
						(*model.UserGroup)(nil),
						(*model.User)(nil),
					}

					for _, this_model := range models {
						db.RegisterModel(this_model)
					}

					fixture := dbfixture.New(db, dbfixture.WithTruncateTables())

					return fixture.Load(ctx, os.DirFS("fixtures"), "fixtures.yaml")
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
