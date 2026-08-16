package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

const (
	lambdaContextHeader  = "x-amzn-lambda-context"
	requestContextHeader = "x-amzn-request-context"
)

type lambdaContext struct {
	RequestID string `json:"request_id"`
}

type requestContext struct {
	RequestID string `json:"requestId"`
}

func WithRequestLogger(baseLogger *zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggerContext := baseLogger.With()

		var lambdaCtx lambdaContext
		if json.Unmarshal([]byte(r.Header.Get(lambdaContextHeader)), &lambdaCtx) == nil && lambdaCtx.RequestID != "" {
			loggerContext = loggerContext.Str("lambdaRequestId", lambdaCtx.RequestID)
		}

		var requestCtx requestContext
		if json.Unmarshal([]byte(r.Header.Get(requestContextHeader)), &requestCtx) == nil && requestCtx.RequestID != "" {
			loggerContext = loggerContext.Str("apiGatewayRequestId", requestCtx.RequestID)
		}

		requestLogger := loggerContext.Logger()
		ctx := requestLogger.WithContext(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
