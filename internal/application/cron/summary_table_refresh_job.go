package cron

import (
	"abt-dashboard-api/internal/domain/repository"
	"context"
	"database/sql"
	"log"

	"github.com/robfig/cron/v3"
)

const logPrefixCronJob = `abt-dashboard-api.internal.application.cron.summary_table_refresh_job`

func StartSummaryTableRefreshJob(db *sql.DB) *cron.Cron {

	if err := RunSummaryTableRefresh(db); err != nil {
		log.Printf("ERROR [%s]: Initial summary refresh failed: %v", logPrefixCronJob, err)
	}

	c := cron.New()

	_, err := c.AddFunc("0 0 * * *", func() {
		log.Printf("INFO [%s]: Running scheduled summary table refresh job...", logPrefixCronJob)

		if err := RunSummaryTableRefresh(db); err != nil {
			log.Printf("ERROR [%s]: Scheduled summary refresh failed: %v", logPrefixCronJob, err)
			return
		}

		log.Printf("INFO [%s]: Scheduled summary tables refreshed successfully.", logPrefixCronJob)
	})

	if err != nil {
		log.Fatalf("FATAL [%s]: Error scheduling cron job: %v", logPrefixCronJob, err)
	}

	c.Start()
	log.Printf("INFO [%s]: Cron job for summary table refresh started.", logPrefixCronJob)

	return c
}

func RunSummaryTableRefresh(db *sql.DB) error {
	log.Printf("INFO [%s]: Running summary table refresh job...", logPrefixCronJob)

	repo := repository.NewTransactionRepository(db)
	ctx := context.Background()

	if err := repo.RefreshSummaryTables(ctx); err != nil {
		log.Printf("ERROR [%s]: Failed to refresh summary tables: %v", logPrefixCronJob, err)
		return err
	}

	log.Printf("INFO [%s]: Summary tables refreshed successfully.", logPrefixCronJob)
	return nil
}
