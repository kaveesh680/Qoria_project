package getTopProducts

import (
	"abt-dashboard-api/internal/domain/entity"
	"abt-dashboard-api/internal/domain/repository"
	"abt-dashboard-api/internal/domain/services"
	"abt-dashboard-api/pkg/errors"
	"database/sql"
	"log"
	"net/http"
)

const prefixHttpHandler = "abt-dashboard-api.internal.application.http.getTopProducts.handler"

// GetTopProductsHandler godoc
// @Summary Get top frequently purchased products
// @Description Returns a list of top frequently purchased products with their total purchased count and current available stock.
// @Tags Metrics
// @Accept  json
// @Produce  json
// @Param limit query int false "Number of top products to return" default(20)
// @Success 200 {array} entity.TopProduct
// @Failure 400 {object} errors.DomainError "Invalid request parameters"
// @Failure 500 {object} errors.ApplicationError "Internal server error"
// @Router /v1/metrics/top-products [get]
func Handler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		decodedReq, err := DecodeTopProductsRequest(ctx, r)
		if err != nil {
			log.Printf("ERROR [%s]: DecodeTopProductsRequest.Error : %v", prefixHttpHandler, err)
			errors.EncodeError(w, errors.NewDomainError("Invalid request params", 10002))
			return
		}
		req := decodedReq.(entity.TopProductsRequest)

		service := services.NewGetTopProductsService(
			repository.NewTransactionRepository(dbConn),
		)

		products, err := service.GetTopProducts(ctx, req.Limit)
		if err != nil {
			log.Printf("ERROR [%s]: GetTopProductsService.GetTopProducts.Error : %v", prefixHttpHandler, err)
			errors.EncodeError(w, errors.NewDomainError("Failed to fetch top products", 12213))
			return
		}

		if err := encodeGetTopProductsResponse(ctx, w, products); err != nil {
			log.Printf("ERROR [%s]: failed to encode top products response: %v", prefixHttpHandler, err)
			errors.EncodeError(w, errors.NewApplicationError("Failed to encode response", 1234235))
			return
		}
	}
}
