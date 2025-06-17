package getTopProducts

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"net/http"
	"strconv"
)

func DecodeTopProductsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()

	limit := 20 // default
	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	return entity.TopProductsRequest{
		Limit: limit,
	}, nil
}
