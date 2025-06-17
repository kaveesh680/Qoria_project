package getMonthlySalesVolume

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"net/http"
	"strconv"
)

func DecodeMonthlySalesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()
	limit := 12 // default to 12 months

	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	return entity.MonthlySalesRequest{Limit: limit}, nil
}
