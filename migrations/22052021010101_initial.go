package migrations

import (
	"context"
	"fmt"

	model "github.com/merefield/buncheck/model"
	"github.com/uptrace/bun"
)

func init() {

	// Drop and create tables.
	models := []interface{}{
		(*model.UserGroup)(nil),
		(*model.User)(nil),
	}

	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")

		for _, model := range models {
			_, err := db.NewDropTable().Model(model).IfExists().Exec(ctx)
			if err != nil {
				panic(err)
			}

			_, err = db.NewCreateTable().Model(model).Exec(ctx)
			if err != nil {
				panic(err)
			}
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		for _, model := range models {
			_, err := db.NewDropTable().Model(model).IfExists().Exec(ctx)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})

}
