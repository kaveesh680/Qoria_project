package getMonthlySalesVolume

import (
	"abt-dashboard-api/internal/domain/entity"
	"context"
	"encoding/json"
	"net/http"
)

func EncodeMonthlySalesResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(*[]entity.MonthlySalesVolume)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(r)
}
