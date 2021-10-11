package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"giautm.dev/awesome/ent"
	"giautm.dev/awesome/ent/migrate"
	_ "giautm.dev/awesome/ent/runtime"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewEntClientFx(lc fx.Lifecycle, logger *zap.Logger) (*ent.Client, error) {
	cfg := Config{
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	db := cfg.OpenDB()

	drv := entsql.OpenDB(dialect.MySQL, db)
	client := ent.NewClient(ent.Driver(drv))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Run the auto migration tool.
			err := client.Schema.Create(ctx,
				migrate.WithGlobalUniqueID(true),
			)
			if err != nil {
				logger.Error(fmt.Sprintf("failed creating schema resources: %v", err))
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Close Ent client.")
			return client.Close()
		},
	})

	return client, nil
}
