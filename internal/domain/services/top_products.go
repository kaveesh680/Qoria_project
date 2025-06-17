package services

import (
	"abt-dashboard-api/internal/domain/boundary"
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"log"
)

type getTopProductsService struct {
	topProductsRepo boundary.TransactionsRepositoryInterface
}

func NewGetTopProductsService(repo boundary.TransactionsRepositoryInterface) boundary.GetTopProductsService {
	return &getTopProductsService{topProductsRepo: repo}
}

func (s *getTopProductsService) GetTopProducts(ctx context.Context, limit int) (products *[]entity.TopProduct, err error) {
	products, err = s.topProductsRepo.GetTopProducts(ctx, limit)
	if err != nil {
		log.Printf("ERROR [services.get_top_products]: %v", err)
		return nil, err
	}
	return products, nil
}
