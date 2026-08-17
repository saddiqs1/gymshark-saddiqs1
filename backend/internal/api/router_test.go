package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/packsizes"
)

type fakePackSizeStore struct {
	sizes     []int
	err       error
	addErr    error
	added     []int
	removeErr error
	removed   []int
}

func (s *fakePackSizeStore) Add(_ context.Context, size int) error {
	if s.addErr != nil {
		return s.addErr
	}
	s.added = append(s.added, size)
	return nil
}

func (s *fakePackSizeStore) Remove(_ context.Context, size int) error {
	if s.removeErr != nil {
		return s.removeErr
	}
	s.removed = append(s.removed, size)
	return nil
}

func (s *fakePackSizeStore) List(context.Context) ([]int, error) {
	return s.sizes, s.err
}

func testRouter() http.Handler {
	return NewRouter(&zerolog.Logger{}, &fakePackSizeStore{})
}

func TestHealth(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	testRouter().ServeHTTP(response, request)

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

	testRouter().ServeHTTP(response, request)

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

			testRouter().ServeHTTP(response, request)

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

			testRouter().ServeHTTP(response, request)

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

func TestGetPackSizes(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/pack-sizes", nil)
	response := httptest.NewRecorder()
	router := NewRouter(&zerolog.Logger{}, &fakePackSizeStore{sizes: []int{250, 500, 1000}})

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != "{\"packSizes\":[250,500,1000]}\n" {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func TestGetPackSizesReturnsInternalServerError(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/pack-sizes", nil)
	response := httptest.NewRecorder()
	router := NewRouter(&zerolog.Logger{}, &fakePackSizeStore{err: errors.New("database unavailable")})

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusInternalServerError)
	}
	if response.Body.String() != "{\"error\":\"failed to retrieve pack sizes\"}\n" {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
}

func TestAddPackSize(t *testing.T) {
	store := &fakePackSizeStore{}
	router := NewRouter(&zerolog.Logger{}, store)
	request := httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBufferString(`{"size":750}`))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusCreated)
	}
	if response.Body.String() != "{\"size\":750}\n" {
		t.Fatalf("unexpected body: %s", response.Body.String())
	}
	if len(store.added) != 1 || store.added[0] != 750 {
		t.Fatalf("added pack sizes = %v; expected [750]", store.added)
	}
}

func TestAddPackSizeRejectsInvalidInput(t *testing.T) {
	tests := []string{"", `{}`, `{"size":0}`, `{"size":-1}`, `{"size":"250"}`, `{"size":250,"unknown":true}`, `{"size":250} {}`}
	for _, body := range tests {
		t.Run(body, func(t *testing.T) {
			response := httptest.NewRecorder()
			testRouter().ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBufferString(body)))
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status code is %d; expected %d", response.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestAddPackSizeReturnsConflictForDuplicate(t *testing.T) {
	router := NewRouter(&zerolog.Logger{}, &fakePackSizeStore{addErr: packsizes.ErrAlreadyExists})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBufferString(`{"size":250}`)))

	if response.Code != http.StatusConflict {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusConflict)
	}
}

func TestAddPackSizeReturnsInternalServerError(t *testing.T) {
	router := NewRouter(&zerolog.Logger{}, &fakePackSizeStore{addErr: errors.New("database unavailable")})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/pack-sizes", bytes.NewBufferString(`{"size":250}`)))

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusInternalServerError)
	}
}

func TestRemovePackSize(t *testing.T) {
	store := &fakePackSizeStore{}
	router := NewRouter(&zerolog.Logger{}, store)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/pack-sizes/750", nil))

	if response.Code != http.StatusNoContent {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusNoContent)
	}
	if response.Body.Len() != 0 {
		t.Fatalf("expected an empty body; got %q", response.Body.String())
	}
	if len(store.removed) != 1 || store.removed[0] != 750 {
		t.Fatalf("removed pack sizes = %v; expected [750]", store.removed)
	}
}

func TestRemovePackSizeRejectsInvalidSize(t *testing.T) {
	for _, target := range []string{"/pack-sizes/not-a-number", "/pack-sizes/0", "/pack-sizes/-1"} {
		t.Run(target, func(t *testing.T) {
			response := httptest.NewRecorder()
			testRouter().ServeHTTP(response, httptest.NewRequest(http.MethodDelete, target, nil))
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status code is %d; expected %d", response.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestRemovePackSizeReturnsNotFound(t *testing.T) {
	router := NewRouter(&zerolog.Logger{}, &fakePackSizeStore{removeErr: packsizes.ErrNotFound})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/pack-sizes/750", nil))

	if response.Code != http.StatusNotFound {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusNotFound)
	}
}

func TestRemovePackSizeReturnsInternalServerError(t *testing.T) {
	router := NewRouter(&zerolog.Logger{}, &fakePackSizeStore{removeErr: errors.New("database unavailable")})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/pack-sizes/750", nil))

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status code is %d; expected %d", response.Code, http.StatusInternalServerError)
	}
}
