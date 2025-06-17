package boundary

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
)

type TransactionsRepositoryInterface interface {
	GetCountryRevenue(ctx context.Context, limit int, offset int) (response *[]entity.CountryRevenueResponse, err error)
	GetTopProducts(ctx context.Context, limit int) (response *[]entity.TopProduct, err error)
	RefreshSummaryTables(ctx context.Context) error
	GetMonthlySalesVolume(ctx context.Context, limit int) (*[]entity.MonthlySalesVolume, error)
	GetTopRegions(ctx context.Context, limit int) (response *[]entity.RegionRevenue, err error)
}
