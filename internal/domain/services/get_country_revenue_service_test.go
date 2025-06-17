package services

import (
	"abt-dashboard-api/internal/domain/entity"
	"abt-dashboard-api/pkg/errors"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockTransactionsRepository implements TransactionsRepositoryInterface for testing
type mockTransactionsRepository struct {
	mock.Mock
}

func (m *mockTransactionsRepository) GetCountryRevenue(ctx context.Context, limit int, offset int) (*[]entity.CountryRevenueResponse, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]entity.CountryRevenueResponse), args.Error(1)
}

// Stub implementations for other interface methods
func (m *mockTransactionsRepository) GetTopProducts(ctx context.Context, limit int) (*[]entity.TopProduct, error) {
	return nil, nil
}
func (m *mockTransactionsRepository) RefreshSummaryTables(ctx context.Context) error {
	return nil
}
func (m *mockTransactionsRepository) GetMonthlySalesVolume(ctx context.Context, limit int) (*[]entity.MonthlySalesVolume, error) {
	return nil, nil
}
func (m *mockTransactionsRepository) GetTopRegions(ctx context.Context, limit int) (*[]entity.RegionRevenue, error) {
	return nil, nil
}

func TestGetCountryRevenue_Success(t *testing.T) {
	mockData := &[]entity.CountryRevenueResponse{
		{Country: "USA", ProductName: "Product A", TotalRevenue: 1000.0, TransactionCount: 10},
		{Country: "Canada", ProductName: "Product B", TotalRevenue: 800.0, TransactionCount: 8},
	}

	mockRepo := new(mockTransactionsRepository)
	mockRepo.On("GetCountryRevenue", mock.Anything, 10, 0).Return(mockData, nil)

	service := NewGetCountryRevenueService(mockRepo)
	result, err := service.GetCountryRevenue(context.Background(), 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, mockData, result)
	assert.Len(t, *result, 2)
	assert.Equal(t, "USA", (*result)[0].Country)
}

func TestGetCountryRevenue_Failure(t *testing.T) {
	mockRepo := new(mockTransactionsRepository)

	// Expect the same error to bubble up from service
	expectedErr := errors.NewDomainError("mock repo error", 10099)

	// Set up the mock to return that error
	mockRepo.On("GetCountryRevenue", mock.Anything, 10, 0).Return(nil, expectedErr)

	// Create service
	service := NewGetCountryRevenueService(mockRepo)

	// Call method under test
	_, err := service.GetCountryRevenue(context.Background(), 10, 0)

	// Assertions
	assert.Error(t, err)

	domainErr, ok := err.(errors.DomainError)
	assert.True(t, ok, "error should be of type DomainError")
	assert.Equal(t, 10099, domainErr.Code)
	assert.Equal(t, "mock repo error", domainErr.Message)

	// Verify that the method was called with expected arguments
	mockRepo.AssertExpectations(t)
}
