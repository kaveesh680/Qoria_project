package getTopRegions

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"net/http"
	"strconv"
)

func DecodeTopRegionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()

	limit := 30 // default
	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	return entity.TopRegionsRequest{Limit: limit}, nil
}
