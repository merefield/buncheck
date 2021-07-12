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
		} else {
			db.NewDelete().
				Model(new(model.User)).
				Where("username LIKE '%user %'").
				Exec(c.Context)

			_, err = db.NewUpdate().
				Model(new(model.User)).
				Set("deleted_at = NULL").
				WhereAllWithDeleted().
				Exec(c.Context)

			if err != nil {
				return err
			}
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

					migrator := migrate.NewMigrator(app.DB(), migrations)

					return migrator.Init(ctx)
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

					migrator := migrate.NewMigrator(app.DB(), migrations)

					_, err = migrator.Migrate(ctx)

					return err
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

					migrator := migrate.NewMigrator(app.DB(), migrations)

					_, err = migrator.Rollback(ctx)

					return err
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

					migrator := migrate.NewMigrator(app.DB(), migrations)

					return migrator.Lock(ctx)
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

					migrator := migrate.NewMigrator(app.DB(), migrations)

					return migrator.Unlock(ctx)
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

					migrator := migrate.NewMigrator(app.DB(), migrations)

					_, err = migrator.CreateGo(ctx, c.Args().Get(0))

					return err
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

					migrator := migrate.NewMigrator(app.DB(), migrations)

					_, err = migrator.CreateSQL(ctx, c.Args().Get(0))

					return err
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
