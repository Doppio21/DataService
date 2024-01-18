package main

import (
	"context"
	"dataservice/internal/api"
	"dataservice/internal/api/ageapi"
	"dataservice/internal/api/genderapi"
	"dataservice/internal/api/nationalizeapi"
	"dataservice/internal/manager"
	"dataservice/internal/pgxprovider"
	"dataservice/internal/server"
	"dataservice/internal/userdb/db"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log, _ := zap.NewDevelopment()

	if err := godotenv.Load(); err != nil {
		log.Error("error loading .env file:", zap.Error(err))
		return
	}
	defer log.Sync()

	pgxp, err := pgxprovider.New(pgxprovider.Config{
		URL: os.Getenv("POSTGRES_URL"),
	})
	if err != nil {
		log.Error("failed postges:", zap.Error(err))
		return
	}
	defer pgxp.Close(context.Background())

	db := db.New(
		db.Config{
			QueryTimeout: 1 * time.Second,
		},
		db.Dependencies{
			Log: log,
			PGX: pgxp,
		},
	)

	client := http.Client{}

	ageapi := ageapi.NewAgify(
		ageapi.Config{
			URI: os.Getenv("AGIFY_URI"),
		},
		ageapi.Dependencies{
			Client: &client,
			Log:    log,
		},
	)

	genderapi := genderapi.NewGenderize(
		genderapi.Config{
			URI: os.Getenv("GENDERIZE_URI"),
		},
		genderapi.Dependencies{
			Client: &client,
			Log:    log,
		},
	)

	nationalizeapi := nationalizeapi.NewNationalize(
		nationalizeapi.Config{
			URI: os.Getenv("NATIONALIZE_URI"),
		},
		nationalizeapi.Dependencies{
			Client: &client,
			Log:    log,
		},
	)

	api := api.NewAPI(
		api.Dependencies{
			Age:         ageapi,
			Gender:      genderapi,
			Nationalize: nationalizeapi,
		},
	)

	manager := manager.New(
		manager.Config{
			Timeout: time.Second,
		},
		manager.Dependencies{
			API: api,
			DB:  db,
			Log: log,
		},
	)

	server := server.New(
		server.Config{
			Address: os.Getenv("SERVER_ADDR"),
		},
		server.Dependencies{
			Manager: *manager,
			Log:     log,
		},
	)
	if err = server.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Error("failed to run server:", zap.Error(err))
	}
}
