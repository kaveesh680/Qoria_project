package services

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockTopProductsRepo mocks TransactionsRepositoryInterface
type mockTopProductsRepo struct {
	mock.Mock
}

func (m *mockTopProductsRepo) GetTopProducts(ctx context.Context, limit int) (*[]entity.TopProduct, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).(*[]entity.TopProduct), args.Error(1)
}

// Stub other methods
func (m *mockTopProductsRepo) GetCountryRevenue(ctx context.Context, limit int, offset int) (*[]entity.CountryRevenueResponse, error) {
	return nil, nil
}
func (m *mockTopProductsRepo) RefreshSummaryTables(ctx context.Context) error {
	return nil
}
func (m *mockTopProductsRepo) GetMonthlySalesVolume(ctx context.Context, limit int) (*[]entity.MonthlySalesVolume, error) {
	return nil, nil
}
func (m *mockTopProductsRepo) GetTopRegions(ctx context.Context, limit int) (*[]entity.RegionRevenue, error) {
	return nil, nil
}

func TestGetTopProducts_Success(t *testing.T) {
	mockRepo := new(mockTopProductsRepo)
	expected := &[]entity.TopProduct{
		{ProductId: "1", ProductName: "Product A", PurchaseCount: 150, AvailableStock: 20},
		{ProductId: "2", ProductName: "Product B", PurchaseCount: 100, AvailableStock: 10},
	}

	mockRepo.On("GetTopProducts", mock.Anything, 2).Return(expected, nil)

	service := NewGetTopProductsService(mockRepo)
	result, err := service.GetTopProducts(context.Background(), 2)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Len(t, *result, 2)
	assert.Equal(t, "Product A", (*result)[0].ProductName)

	mockRepo.AssertExpectations(t)
}

func TestGetTopProducts_Failure(t *testing.T) {
	mockRepo := new(mockTopProductsRepo)
	expectedErr := errors.New("repository failure")

	mockRepo.On("GetTopProducts", mock.Anything, 2).Return((*[]entity.TopProduct)(nil), expectedErr)

	service := NewGetTopProductsService(mockRepo)
	result, err := service.GetTopProducts(context.Background(), 2)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}
