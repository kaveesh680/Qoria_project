package services

import (
	"abt-dashboard-api/internal/domain/boundary"
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"log"
)

const prefixGetCountryRevenue = `abt-dashboard-api.internal.domain.services.get_country_revenue`

type getCountryRevenueService struct {
	transactionsRepository boundary.TransactionsRepositoryInterface
}

func NewGetCountryRevenueService(transactionsRepository boundary.TransactionsRepositoryInterface) boundary.GetCountryRevenueService {
	return &getCountryRevenueService{
		transactionsRepository: transactionsRepository,
	}
}

func (g *getCountryRevenueService) GetCountryRevenue(ctx context.Context, limit int, offset int) (countryRevenueDetails *[]entity.CountryRevenueResponse, err error) {

	countryRevenueDetails, err = g.transactionsRepository.GetCountryRevenue(ctx, limit, offset)
	if err != nil {
		log.Printf("ERROR [%s]: g.transactionsRepository.GetCountryRevenue.Error: %v", prefixGetCountryRevenue, err)
		return nil, err
	}

	return countryRevenueDetails, nil

}
