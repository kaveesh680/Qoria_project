package ping

import (
	"context"
	"net/http"
)

func DecodePingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}
