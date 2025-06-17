package boundary

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
)

type GetCountryRevenueService interface {
	GetCountryRevenue(ctx context.Context, limit int, offset int) (countryRevenueDetails *[]entity.CountryRevenueResponse, err error)
}

type GetTopProductsService interface {
	GetTopProducts(ctx context.Context, limit int) (products *[]entity.TopProduct, err error)
}

type MonthlySalesService interface {
	GetMonthlySalesVolume(ctx context.Context, limit int) (monthlySalesVolumes *[]entity.MonthlySalesVolume, err error)
}

type GetTopRegionsService interface {
	GetTopRegions(ctx context.Context, limit int) (topRegions *[]entity.RegionRevenue, err error)
}
