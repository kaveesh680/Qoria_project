package services

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockMonthlySalesRepo mocks TransactionsRepositoryInterface for testing MonthlySalesService
type mockMonthlySalesRepo struct {
	mock.Mock
}

func (m *mockMonthlySalesRepo) GetMonthlySalesVolume(ctx context.Context, limit int) (*[]entity.MonthlySalesVolume, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).(*[]entity.MonthlySalesVolume), args.Error(1)
}

// Stub implementations for unused methods
func (m *mockMonthlySalesRepo) GetCountryRevenue(ctx context.Context, limit, offset int) (*[]entity.CountryRevenueResponse, error) {
	return nil, nil
}
func (m *mockMonthlySalesRepo) GetTopProducts(ctx context.Context, limit int) (*[]entity.TopProduct, error) {
	return nil, nil
}
func (m *mockMonthlySalesRepo) RefreshSummaryTables(ctx context.Context) error {
	return nil
}
func (m *mockMonthlySalesRepo) GetTopRegions(ctx context.Context, limit int) (*[]entity.RegionRevenue, error) {
	return nil, nil
}

func TestGetMonthlySalesVolume_Success(t *testing.T) {
	mockRepo := new(mockMonthlySalesRepo)
	expected := &[]entity.MonthlySalesVolume{
		{Month: "2024-01", TotalSales: 300},
		{Month: "2024-02", TotalSales: 280},
	}

	mockRepo.On("GetMonthlySalesVolume", mock.Anything, 2).Return(expected, nil)

	service := NewMonthlySalesService(mockRepo)
	result, err := service.GetMonthlySalesVolume(context.Background(), 2)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetMonthlySalesVolume_Failure(t *testing.T) {
	mockRepo := new(mockMonthlySalesRepo)
	expectedErr := errors.New("db failure")

	mockRepo.On("GetMonthlySalesVolume", mock.Anything, 3).Return((*[]entity.MonthlySalesVolume)(nil), expectedErr)

	service := NewMonthlySalesService(mockRepo)
	result, err := service.GetMonthlySalesVolume(context.Background(), 3)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}
