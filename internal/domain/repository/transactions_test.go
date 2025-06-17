package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCountryRevenue(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	rows := sqlmock.NewRows([]string{"country", "product_name", "total_revenue", "transaction_count"}).
		AddRow("USA", "Product A", 1000.0, 10).
		AddRow("Canada", "Product B", 800.0, 8)

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT country, product_name, total_revenue, transaction_count
        FROM transaction_summary
        ORDER BY total_revenue DESC
        LIMIT ? OFFSET ?;
    `)).
		WithArgs(2, 0).
		WillReturnRows(rows)

	result, err := repo.GetCountryRevenue(context.Background(), 2, 0)

	assert.NoError(t, err)
	assert.Len(t, *result, 2)
	assert.Equal(t, "USA", (*result)[0].Country)
	assert.Equal(t, "Product A", (*result)[0].ProductName)
}

func TestGetTopProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	rows := sqlmock.NewRows([]string{"product_id", "product_name", "total_purchased", "stock_quantity"}).
		AddRow(1, "Product A", 150, 20).
		AddRow(2, "Product B", 100, 10)

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT 
    product_id,
    product_name,
    total_purchased,
    stock_quantity
    FROM product_purchase_summary
    ORDER BY total_purchased DESC
        LIMIT ?;
    `)).
		WithArgs(2).
		WillReturnRows(rows)

	result, err := repo.GetTopProducts(context.Background(), 2)

	assert.NoError(t, err)
	assert.Len(t, *result, 2)
	assert.Equal(t, "Product A", (*result)[0].ProductName)
}

func TestGetMonthlySalesVolume(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	rows := sqlmock.NewRows([]string{"month", "total_sales"}).
		AddRow("2024-01", 300).
		AddRow("2024-02", 250)

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT DATE_FORMAT(transaction_date, '%Y-%m') AS month, SUM(quantity) AS total_sales
        FROM transactions
        GROUP BY month
        ORDER BY total_sales DESC
        LIMIT ?;
    `)).
		WithArgs(2).
		WillReturnRows(rows)

	result, err := repo.GetMonthlySalesVolume(context.Background(), 2)

	assert.NoError(t, err)
	assert.Len(t, *result, 2)
	assert.Equal(t, "2024-01", (*result)[0].Month)
}

func TestGetTopRegions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	rows := sqlmock.NewRows([]string{"region", "total_revenue", "items_sold"}).
		AddRow("West", 2000, 100).
		AddRow("East", 1500, 80)

	mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT 
        region,
        SUM(price * quantity) AS total_revenue,
        SUM(quantity) AS items_sold
        FROM transactions
        GROUP BY region
        ORDER BY total_revenue DESC
        LIMIT ?;
    `)).
		WithArgs(2).
		WillReturnRows(rows)

	result, err := repo.GetTopRegions(context.Background(), 2)

	assert.NoError(t, err)
	assert.Len(t, *result, 2)
	assert.Equal(t, "West", (*result)[0].Region)
}

func TestRefreshSummaryTables_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	mock.ExpectExec("DROP TABLE IF EXISTS transaction_summary").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("CREATE TABLE transaction_summary AS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("DROP TABLE IF EXISTS product_purchase_summary").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("CREATE TABLE product_purchase_summary AS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.RefreshSummaryTables(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRefreshSummaryTables_DropTransactionSummaryFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	mock.ExpectExec("DROP TABLE IF EXISTS transaction_summary").
		WillReturnError(errors.New("drop failed"))

	err = repo.RefreshSummaryTables(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "drop failed", err.Error())
}

func TestRefreshSummaryTables_CreateTransactionSummaryFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	mock.ExpectExec("DROP TABLE IF EXISTS transaction_summary").
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec("CREATE TABLE transaction_summary AS").
		WillReturnError(errors.New("create failed"))

	err = repo.RefreshSummaryTables(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "create failed", err.Error())
}

func TestRefreshSummaryTables_DropProductPurchaseSummaryFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	mock.ExpectExec("DROP TABLE IF EXISTS transaction_summary").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE transaction_summary AS").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TABLE IF EXISTS product_purchase_summary").
		WillReturnError(errors.New("drop product summary failed"))

	err = repo.RefreshSummaryTables(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "drop product summary failed", err.Error())
}

func TestRefreshSummaryTables_CreateProductPurchaseSummaryFails(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewTransactionRepository(db)

	mock.ExpectExec("DROP TABLE IF EXISTS transaction_summary").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE transaction_summary AS").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("DROP TABLE IF EXISTS product_purchase_summary").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE product_purchase_summary AS").
		WillReturnError(errors.New("create product summary failed"))

	err = repo.RefreshSummaryTables(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "create product summary failed", err.Error())
}
