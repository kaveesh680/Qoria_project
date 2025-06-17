package ping

import (
	"context"
	"encoding/json"
	"net/http"
)

func EncodePingResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	r := response.(string)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(r)
}
