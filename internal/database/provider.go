package database

import (
	"context"
	"fmt"
	"os"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gocloud.dev/server/health/sqlhealth"

	"giautm.dev/awesome/ent"
)

// Module exports the database module.
var Module = fx.Options(
	fx.Provide(NewEntClientFx),
	fx.Provide(func(e *ent.Client) *sqlhealth.Checker {
		return e.HealthCheck()
	}),
)

// NewEntClientFx returns a new ent.Client.
func NewEntClientFx(lc fx.Lifecycle, logger *zap.Logger) (*ent.Client, error) {
	cfg := Config{
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	db := cfg.OpenDB()

	drv := sql.OpenDB(dialect.MySQL, db)
	client := ent.NewClient(ent.Driver(drv))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Run the auto migration tool.
			logger.Info("Ent Migrating.")
			err := client.Schema.Create(ctx, schema.WithAtlas(true))
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
