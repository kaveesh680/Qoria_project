package internal

import (
	"abt-dashboard-api/internal/application/cron"
	"abt-dashboard-api/internal/application/database"
	"abt-dashboard-api/internal/application/http"
	pkgDB "abt-dashboard-api/pkg/database"

	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Init @title ABT Dashboard API
// @version 1.0
// @description This is the backend API documentation for ABT Corporation’s dashboard.
// @host localhost:8080
// @BasePath /
func Init() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	os.Setenv("DB_NAME", "abt_dashboard")
	os.Setenv("DB_USERNAME", "root")
	os.Setenv("DB_PASSWORD", "123")

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("DB_NAME environment variable not set")
	}

	dbUserName := os.Getenv("DB_USERNAME")
	if dbUserName == "" {
		panic("DB_USERNAME environment variable not set")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		panic("DB_PASSWORD environment variable not set")
	}

	abtDashboardDBConn, err := database.NewDbConfig(
		dbUserName,
		dbPassword,
		dbName,
		"127.0.0.1",
		"3306",
		"300s",
		"300s",
		"30s",
	).NewDatabaseConnection(pkgDB.DefaultConnector{})
	if err != nil {
		log.Fatal(err)
	}

	defer func(abtDashboardDatabase *sql.DB) {
		err := abtDashboardDatabase.Close()
		if err != nil {
			log.Fatalf("Error occurred while closing Database connection pool  : %v", err)
		}
	}(abtDashboardDBConn)

	cronJob := cron.StartSummaryTableRefreshJob(abtDashboardDBConn)
	defer cronJob.Stop()

	httpServer := http.NewServer(abtDashboardDBConn)
	httpServer.Start(ctx)

	select {
	case <-sigs:
		httpServer.Stop(ctx)
	}
}
