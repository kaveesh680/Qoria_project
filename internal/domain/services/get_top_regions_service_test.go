package services

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTopRegionsRepository struct {
	mock.Mock
}

func (m *mockTopRegionsRepository) GetTopRegions(ctx context.Context, limit int) (*[]entity.RegionRevenue, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).(*[]entity.RegionRevenue), args.Error(1)
}

// Stub implementations (not used in this test)
func (m *mockTopRegionsRepository) GetCountryRevenue(ctx context.Context, limit, offset int) (*[]entity.CountryRevenueResponse, error) {
	return nil, nil
}
func (m *mockTopRegionsRepository) GetTopProducts(ctx context.Context, limit int) (*[]entity.TopProduct, error) {
	return nil, nil
}
func (m *mockTopRegionsRepository) RefreshSummaryTables(ctx context.Context) error {
	return nil
}
func (m *mockTopRegionsRepository) GetMonthlySalesVolume(ctx context.Context, limit int) (*[]entity.MonthlySalesVolume, error) {
	return nil, nil
}

func TestGetTopRegions_Success(t *testing.T) {
	mockRepo := new(mockTopRegionsRepository)
	expected := &[]entity.RegionRevenue{
		{Region: "West", TotalRevenue: 10000, ItemsSold: 200},
		{Region: "East", TotalRevenue: 8000, ItemsSold: 150},
	}

	mockRepo.On("GetTopRegions", mock.Anything, 2).Return(expected, nil)

	service := NewGetTopRegionsService(mockRepo)
	result, err := service.GetTopRegions(context.Background(), 2)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetTopRegions_Failure(t *testing.T) {
	mockRepo := new(mockTopRegionsRepository)
	expectedErr := errors.New("db query failed")

	mockRepo.On("GetTopRegions", mock.Anything, 5).Return((*[]entity.RegionRevenue)(nil), expectedErr)

	service := NewGetTopRegionsService(mockRepo)
	result, err := service.GetTopRegions(context.Background(), 5)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}
