package getTopRegions

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"encoding/json"
	"net/http"
)

func EncodeTopRegionsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response.(*[]entity.RegionRevenue))
}
