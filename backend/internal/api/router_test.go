package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	actualResponseCode := response.Code
	expectedResponseCode := http.StatusOK
	if actualResponseCode != expectedResponseCode {
		t.Fatalf("status code is %v; expected %v", actualResponseCode, expectedResponseCode)
	}

	actualContentType := response.Header().Get("Content-Type")
	expectedContentType := "application/json"
	if actualContentType != expectedContentType {
		t.Fatalf("Content-Type is %v; expected %v", actualContentType, expectedContentType)
	}

	actualBody := response.Body.String()
	expectedBody := "{\"status\":\"ok\"}\n"
	if actualBody != expectedBody {
		t.Fatalf("body is %v; expected %v", actualBody, expectedBody)
	}
}

func TestGetPacks(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/packs?itemsOrdered=501&packSizes=250,500,1000,2000,5000", nil)
	response := httptest.NewRecorder()

	NewRouter().ServeHTTP(response, request)

	actualResponseCode := response.Code
	expectedResponseCode := http.StatusOK
	if actualResponseCode != expectedResponseCode {
		t.Fatalf("status code is %v; expected %v", actualResponseCode, expectedResponseCode)
	}

	actualBody := response.Body.String()
	expectedBody := "{\"250\":1,\"500\":1}\n"
	if actualBody != expectedBody {
		t.Fatalf("body is %v; expected %v", actualBody, expectedBody)
	}
}

func TestGetPacksRejectsInvalidItemsOrdered(t *testing.T) {
	tests := []string{
		"/packs?packSizes=1",
		"/packs?itemsOrdered=invalid&packSizes=1",
		"/packs?itemsOrdered=0&packSizes=1",
		"/packs?itemsOrdered=-1&packSizes=1",
	}

	for _, target := range tests {
		t.Run(target, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, target, nil)
			response := httptest.NewRecorder()

			NewRouter().ServeHTTP(response, request)

			actualResponseCode := response.Code
			expectedResponseCode := http.StatusBadRequest
			if actualResponseCode != expectedResponseCode {
				t.Fatalf("status code is %v; expected %v", actualResponseCode, expectedResponseCode)
			}

			actualBody := response.Body.String()
			expectedBody := "{\"error\":\"itemsOrdered must be a positive integer\"}\n"
			if actualBody != expectedBody {
				t.Fatalf("body is %v; expected %v", actualBody, expectedBody)
			}
		})
	}
}

func TestGetPacksRejectsInvalidPackSizes(t *testing.T) {
	tests := []string{
		"/packs?itemsOrdered=1",
		"/packs?itemsOrdered=1&packSizes=invalid",
		"/packs?itemsOrdered=1&packSizes=0",
		"/packs?itemsOrdered=1&packSizes=10,0,15",
		"/packs?itemsOrdered=1&packSizes=10,-1,15",
	}

	for _, target := range tests {
		t.Run(target, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, target, nil)
			response := httptest.NewRecorder()

			NewRouter().ServeHTTP(response, request)

			actualResponseCode := response.Code
			expectedResponseCode := http.StatusBadRequest
			if actualResponseCode != expectedResponseCode {
				t.Fatalf("status code is %v; expected %v", actualResponseCode, expectedResponseCode)
			}

			actualBody := response.Body.String()
			expectedBody := "{\"error\":\"packSizes must be a list of positive integers\"}\n"
			if actualBody != expectedBody {
				t.Fatalf("body is %v; expected %v", actualBody, expectedBody)
			}
		})
	}
}
