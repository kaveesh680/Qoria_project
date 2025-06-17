package getCountryRevenue

import (
	"abt-dashboard-api/internal/domain/entity"
	"abt-dashboard-api/internal/domain/repository"
	"abt-dashboard-api/internal/domain/services"
	"abt-dashboard-api/pkg/errors"
	"database/sql"
	"log"
	"net/http"
)

const prefixHttpHandler = "abt-dashboard-api.internal.application.http.getCountryRevenue.handler"

// GetCountryRevenueHandler godoc
// @Summary Get country-level revenue data
// @Description Returns a paginated list of countries with product revenue and transaction count, sorted by total revenue.
// @Tags Metrics
// @Accept  json
// @Produce  json
// @Param limit query int false "Number of records to return" default(50)
// @Param offset query int false "Number of records to skip" default(0)
// @Success 200 {array} entity.CountryRevenueResponse
// @Failure 400 {object} errors.DomainError "Invalid request parameters"
// @Failure 500 {object} errors.ApplicationError "Internal server error"
// @Router /v1/metrics/country-revenue [get]
func Handler(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		decodedReq, err := DecodeGetCountryRevenueRequest(ctx, r)
		if err != nil {
			log.Printf("ERROR [%s]: DecodeGetCountryRevenueRequest.Error : %v", prefixHttpHandler, err)
			errors.EncodeError(w, errors.NewDomainError("Invalid request params", 10001))
			return
		}
		req := decodedReq.(entity.GetCountryRevenueRequest)

		getCountryRevenueService := services.NewGetCountryRevenueService(
			repository.NewTransactionRepository(dbConn),
		)

		countryRevenues, err := getCountryRevenueService.GetCountryRevenue(ctx, req.Limit, req.Offset)
		if err != nil {
			log.Printf("ERROR [%s]: getCountryRevenyeService.GetCountryRevenue.Error : %v", prefixHttpHandler, err)
			domainErr := errors.NewDomainError("Oops, Something went wrong!", 12212)
			errors.EncodeError(w, domainErr)
			return
		}

		if err := EncodeGetCountryRevenueResponse(ctx, w, countryRevenues); err != nil {
			log.Printf("ERROR [%s]: failed to encode country revenue response: %v", prefixHttpHandler, err)
			appErr := errors.NewApplicationError("Oops, Something went wrong!", 1234234)
			errors.EncodeError(w, appErr)
			return
		}
	}
}
