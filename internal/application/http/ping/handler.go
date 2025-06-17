package ping

import (
	"abt-dashboard-api/pkg/errors"
	"log"
	"net/http"
)

const prefixHttpHandler = "abt-dashboard-api.internal.application.http.ping.handler"

// PingHandler godoc
// @Summary Health Check
// @Description Returns a simple ping response to verify that the server is running.
// @Tags Utility
// @Accept  json
// @Produce  json
// @Success 200 {string} string "ping"
// @Failure 500 {object} errors.ApplicationError
// @Router /ping [get]
func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	response := "ping"
	if err := EncodePingResponse(ctx, w, response); err != nil {
		log.Printf("ERROR [%s]: failed to encode ping response: %v", prefixHttpHandler, err)
		appErr := errors.NewApplicationError("Oops, Something went wrong!", 1234234)
		errors.EncodeError(w, appErr)
	}
}
