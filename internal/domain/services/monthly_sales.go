package services

import (
	"abt-dashboard-api/internal/domain/boundary"
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"log"
)

type monthlySalesService struct {
	transactionsRepository boundary.TransactionsRepositoryInterface
}

func NewMonthlySalesService(transactionsRepository boundary.TransactionsRepositoryInterface) boundary.MonthlySalesService {
	return &monthlySalesService{transactionsRepository: transactionsRepository}
}

func (s *monthlySalesService) GetMonthlySalesVolume(ctx context.Context, limit int) (monthlySalesVolumes *[]entity.MonthlySalesVolume, err error) {
	monthlySalesVolumes, err = s.transactionsRepository.GetMonthlySalesVolume(ctx, limit)
	if err != nil {
		log.Printf("ERROR [services.monthly_sales]: %v", err)
		return nil, err
	}
	return monthlySalesVolumes, nil
}
