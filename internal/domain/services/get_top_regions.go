package services

import (
	"abt-dashboard-api/internal/domain/boundary"
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"log"
)

type getTopRegionsService struct {
	repo boundary.TransactionsRepositoryInterface
}

func NewGetTopRegionsService(repo boundary.TransactionsRepositoryInterface) boundary.GetTopRegionsService {
	return &getTopRegionsService{repo: repo}
}

func (s *getTopRegionsService) GetTopRegions(ctx context.Context, limit int) (topRegions *[]entity.RegionRevenue, err error) {
	topRegions, err = s.repo.GetTopRegions(ctx, limit)
	if err != nil {
		log.Printf("ERROR [services.get_top_regions]: %v", err)
		return nil, err
	}
	return topRegions, nil
}
