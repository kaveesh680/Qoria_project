package getCountryRevenue

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"net/http"
	"strconv"
)

func DecodeGetCountryRevenueRequest(_ context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()

	limit := 50
	offset := 0

	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	if o, err := strconv.Atoi(query.Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	return entity.GetCountryRevenueRequest{
		Limit:  limit,
		Offset: offset,
	}, nil
}
