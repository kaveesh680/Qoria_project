package getTopRegions

import (
	"abt-dashboard-api/internal/domain/entity"
	"abt-dashboard-api/internal/domain/repository"
	"abt-dashboard-api/internal/domain/services"
	"abt-dashboard-api/pkg/errors"
	"database/sql"
	"log"
	"net/http"
)

const logPrefix = "abt-dashboard-api.internal.application.http.getTopRegions.handler"

// GetTopRegionsHandler godoc
// @Summary Get top regions by revenue
// @Description Returns a list of regions with the highest total revenue and items sold, sorted by revenue.
// @Tags Metrics
// @Accept  json
// @Produce  json
// @Param limit query int false "Number of top regions to return" default(30)
// @Success 200 {array} entity.RegionRevenue
// @Failure 400 {object} errors.DomainError "Invalid request parameters"
// @Failure 500 {object} errors.ApplicationError "Internal server error"
// @Router /v1/metrics/top-regions [get]
func Handler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		decodedReq, err := DecodeTopRegionsRequest(ctx, r)
		if err != nil {
			log.Printf("ERROR [%s]: DecodeTopRegionsRequest.Error : %v", logPrefix, err)
			errors.EncodeError(w, errors.NewDomainError("Invalid request params", 10003))
			return
		}
		req := decodedReq.(entity.TopRegionsRequest)

		service := services.NewGetTopRegionsService(repository.NewTransactionRepository(db))
		result, err := service.GetTopRegions(ctx, req.Limit)
		if err != nil {
			log.Printf("ERROR [%s]: GetTopRegionsService.GetTopRegions.Error : %v", logPrefix, err)
			errors.EncodeError(w, errors.NewDomainError("Failed to fetch top regions", 12214))
			return
		}

		if err := EncodeTopRegionsResponse(ctx, w, result); err != nil {
			log.Printf("ERROR [%s]: failed to encode response: %v", logPrefix, err)
			errors.EncodeError(w, errors.NewApplicationError("Failed to encode response", 1234236))
			return
		}
	}
}
