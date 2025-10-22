package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jmontesinos91/terrors"

	"github.com/go-chi/chi/v5/middleware"
)

// RenderJSON Render A helper function to render a JSON response
func RenderJSON(ctx context.Context, w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	// Headers
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(ctx))
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(js)
}

// RenderFile Render A helper function to render a File response
func RenderFile(ctx context.Context, w http.ResponseWriter, httpStatusCode int, payload []byte) {
	// Headers
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(ctx))
	w.Header().Set("Content-Type", "application/octet-stream")

	w.WriteHeader(httpStatusCode)
	_, _ = w.Write(payload)
}

// RenderError Renders a error with some sane defaults.
// This function receive any type of error, but is recommended use a terror
// for cases when you what to send a specific status code, because other kind
// of errors are handled as internal_errors
func RenderError(ctx context.Context, w http.ResponseWriter, err error) {
	var httpStatusCode int
	var code string
	var message string

	// Check if error can be parsed as terror
	var terr *terrors.Error
	if errors.As(err, &terr) {
		if terr.PrefixMatches(terrors.ErrPreconditionFailed) || terr.PrefixMatches(terrors.ErrBadRequest) {
			httpStatusCode = http.StatusBadRequest
		} else if terr.PrefixMatches(terrors.ErrUnauthorized) {
			httpStatusCode = http.StatusUnauthorized
		} else if terr.PrefixMatches(terrors.ErrNotFound) {
			httpStatusCode = http.StatusNotFound
		} else {
			httpStatusCode = http.StatusInternalServerError
		}
		// Assing Code and Error message
		code = terr.Code
		message = terr.Message
	}

	payload := map[string]string{
		"code":    code,
		"message": message,
	}

	RenderJSON(ctx, w, httpStatusCode, payload)
}
