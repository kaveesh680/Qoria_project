package repository

import (
	"abt-dashboard-api/internal/domain/boundary"
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"database/sql"
	"log"
)

const logPrefixTransactionRepository = "abt-dashboard-api.internal.domain.repository.get_country_revenue"

type transactionRepo struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) boundary.TransactionsRepositoryInterface {
	return &transactionRepo{
		db: db,
	}
}

func (r *transactionRepo) GetCountryRevenue(ctx context.Context, limit int, offset int) (response *[]entity.CountryRevenueResponse, err error) {

	query := `
		SELECT country, product_name, total_revenue, transaction_count
		FROM transaction_summary
		ORDER BY total_revenue DESC
		LIMIT ? OFFSET ?;
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Printf("ERROR [%s]: r.db.QueryContext.Error: %v", logPrefixTransactionRepository, err)
		return nil, err
	}
	defer rows.Close()

	var results []entity.CountryRevenueResponse
	for rows.Next() {
		var row entity.CountryRevenueResponse
		if err := rows.Scan(&row.Country, &row.ProductName, &row.TotalRevenue, &row.TransactionCount); err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	return &results, nil
}

func (r *transactionRepo) GetTopProducts(ctx context.Context, limit int) (response *[]entity.TopProduct, err error) {

	query := `
	SELECT 
    product_id,
    product_name,
    total_purchased,
    stock_quantity
	FROM product_purchase_summary
	ORDER BY total_purchased DESC
		LIMIT ?;
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		log.Printf("ERROR [repository.top_products]: Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []entity.TopProduct
	for rows.Next() {
		var row entity.TopProduct
		if err := rows.Scan(&row.ProductId, &row.ProductName, &row.PurchaseCount, &row.AvailableStock); err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	return &results, nil
}

func (r *transactionRepo) RefreshSummaryTables(ctx context.Context) error {
	// Step 1: Refresh transaction_summary
	if _, err := r.db.ExecContext(ctx, `DROP TABLE IF EXISTS transaction_summary`); err != nil {
		log.Printf("ERROR [%s]: Failed to drop transaction_summary: %v", logPrefixTransactionRepository, err)
		return err
	}

	if _, err := r.db.ExecContext(ctx, `
		CREATE TABLE transaction_summary AS
		SELECT
			country,
			product_name,
			SUM(price * quantity) AS total_revenue,
			COUNT(*) AS transaction_count
		FROM transactions
		GROUP BY country, product_name
	`); err != nil {
		log.Printf("ERROR [%s]: Failed to create transaction_summary: %v", logPrefixTransactionRepository, err)
		return err
	}

	// Step 2: Refresh product_purchase_summary
	if _, err := r.db.ExecContext(ctx, `DROP TABLE IF EXISTS product_purchase_summary`); err != nil {
		log.Printf("ERROR [%s]: Failed to drop product_purchase_summary: %v", logPrefixTransactionRepository, err)
		return err
	}

	if _, err := r.db.ExecContext(ctx, `
		CREATE TABLE product_purchase_summary AS
		SELECT 
			product_id,
			product_name,
			SUM(quantity) AS total_purchased,
			MAX(stock_quantity) AS stock_quantity
		FROM transactions
		GROUP BY product_id, product_name
	`); err != nil {
		log.Printf("ERROR [%s]: Failed to create product_purchase_summary: %v", logPrefixTransactionRepository, err)
		return err
	}

	log.Printf("INFO [%s]: Summary tables refreshed successfully (via DROP + CREATE)", logPrefixTransactionRepository)
	return nil
}

func (r *transactionRepo) GetMonthlySalesVolume(ctx context.Context, limit int) (*[]entity.MonthlySalesVolume, error) {
	query := `
		SELECT DATE_FORMAT(transaction_date, '%Y-%m') AS month, SUM(quantity) AS total_sales
		FROM transactions
		GROUP BY month
		ORDER BY total_sales DESC
		LIMIT ?;
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		log.Printf("ERROR [%s]: Failed to fetch monthly sales: %v", logPrefixTransactionRepository, err)
		return nil, err
	}
	defer rows.Close()

	var results []entity.MonthlySalesVolume
	for rows.Next() {
		var row entity.MonthlySalesVolume
		if err := rows.Scan(&row.Month, &row.TotalSales); err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	return &results, nil
}

func (r *transactionRepo) GetTopRegions(ctx context.Context, limit int) (*[]entity.RegionRevenue, error) {
	query := `
	SELECT 
		region,
		SUM(price * quantity) AS total_revenue,
		SUM(quantity) AS items_sold
	FROM transactions
	GROUP BY region
	ORDER BY total_revenue DESC
	LIMIT ?;
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		log.Printf("ERROR [repository.top_regions]: Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []entity.RegionRevenue
	for rows.Next() {
		var row entity.RegionRevenue
		if err := rows.Scan(&row.Region, &row.TotalRevenue, &row.ItemsSold); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return &results, nil
}
