package getMonthlySalesVolume

import (
	"abt-dashboard-api/internal/domain/entity"
	"abt-dashboard-api/internal/domain/repository"
	"abt-dashboard-api/internal/domain/services"
	"abt-dashboard-api/pkg/errors"
	"database/sql"
	"log"
	"net/http"
)

const prefix = "abt-dashboard-api.internal.application.http.getMonthlySalesVolume.handler"

// GetMonthlySalesVolumeHandler godoc
// @Summary Get monthly sales volume
// @Description Returns a list of months with their total item sales, sorted by highest sales volume.
// @Tags Metrics
// @Accept  json
// @Produce  json
// @Param limit query int false "Number of top months to return" default(12)
// @Success 200 {array} entity.MonthlySalesVolume
// @Failure 400 {object} errors.DomainError "Invalid request parameters"
// @Failure 500 {object} errors.ApplicationError "Internal server error"
// @Router /v1/metrics/monthly-sales [get]
func Handler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		decodedReq, err := DecodeMonthlySalesRequest(ctx, r)
		if err != nil {
			log.Printf("ERROR [%s]: DecodeMonthlySalesRequest: %v", prefix, err)
			errors.EncodeError(w, errors.NewDomainError("Invalid request params", 10002))
			return
		}
		req := decodedReq.(entity.MonthlySalesRequest)

		service := services.NewMonthlySalesService(repository.NewTransactionRepository(db))

		data, err := service.GetMonthlySalesVolume(ctx, req.Limit)
		if err != nil {
			log.Printf("ERROR [%s]: GetMonthlySalesVolume: %v", prefix, err)
			errors.EncodeError(w, errors.NewDomainError("Failed to fetch sales data", 12214))
			return
		}

		if err := EncodeMonthlySalesResponse(ctx, w, data); err != nil {
			log.Printf("ERROR [%s]: Encode response: %v", prefix, err)
			errors.EncodeError(w, errors.NewApplicationError("Failed to encode response", 1234236))
			return
		}
	}
}
