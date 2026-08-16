package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
)

func TestWithRequestLoggerAddsLambdaAndAPIGatewayRequestIDs(t *testing.T) {
	var output bytes.Buffer
	baseLogger := zerolog.New(&output)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zerolog.Ctx(r.Context()).Info().Msg("request handled")
		w.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set(lambdaContextHeader, `{"request_id":"lambda-request-id"}`)
	request.Header.Set(requestContextHeader, `{"requestId":"api-gateway-request-id"}`)
	response := httptest.NewRecorder()

	WithRequestLogger(&baseLogger, next).ServeHTTP(response, request)

	var logEvent map[string]any
	if err := json.Unmarshal(output.Bytes(), &logEvent); err != nil {
		t.Fatalf("unmarshal log event: %v", err)
	}

	if actual := logEvent["lambdaRequestId"]; actual != "lambda-request-id" {
		t.Errorf("lambdaRequestId = %q, want %q", actual, "lambda-request-id")
	}

	if actual := logEvent["apiGatewayRequestId"]; actual != "api-gateway-request-id" {
		t.Errorf("apiGatewayRequestId = %q, want %q", actual, "api-gateway-request-id")
	}
}

func TestWithRequestLoggerWithoutRequestIDs(t *testing.T) {
	var output bytes.Buffer
	baseLogger := zerolog.New(&output)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zerolog.Ctx(r.Context()).Info().Msg("request handled")
		w.WriteHeader(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	WithRequestLogger(&baseLogger, next).ServeHTTP(response, request)

	var logEvent map[string]any
	if err := json.Unmarshal(output.Bytes(), &logEvent); err != nil {
		t.Fatalf("unmarshal log event: %v", err)
	}

	if _, exists := logEvent["lambdaRequestId"]; exists {
		t.Error("lambdaRequestId should not be logged when the Lambda context header is absent")
	}

	if _, exists := logEvent["apiGatewayRequestId"]; exists {
		t.Error("apiGatewayRequestId should not be logged when the API Gateway context header is absent")
	}
}
