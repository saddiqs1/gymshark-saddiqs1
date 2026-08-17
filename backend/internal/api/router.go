package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/packs"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/packsizes"
)

type router struct {
	logger    *zerolog.Logger
	packSizes packsizes.Store
}

func NewRouter(logger *zerolog.Logger, packSizes packsizes.Store) http.Handler {
	r := &router{logger: logger, packSizes: packSizes}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", r.health)
	mux.HandleFunc("GET /packs-for-items-ordered", r.getPacksForItemsOrdered)
	mux.HandleFunc("GET /pack-sizes", r.getPackSizes)
	mux.HandleFunc("POST /pack-sizes", r.addPackSize)
	mux.HandleFunc("DELETE /pack-sizes/{size}", r.removePackSize)

	return mux
}

func (r *router) removePackSize(w http.ResponseWriter, req *http.Request) {
	r.logger.Info().Msgf("Received request: %s %s", req.Method, req.URL.String())
	r.logger.Debug().Msgf("Path parameter: %v", req.PathValue("size"))

	size, err := strconv.Atoi(req.PathValue("size"))
	if err != nil || size <= 0 {
		r.logger.Error().Err(err).Msg("Invalid size")
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "size must be a positive integer",
		})
		return
	}

	if err := r.packSizes.Remove(req.Context(), size); err != nil {
		if errors.Is(err, packsizes.ErrNotFound) {
			r.logger.Error().Err(err).Msg("Size not found")
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "pack size not found",
			})
			return
		}
		r.logger.Error().Err(err).Msg("Failed to remove pack size")
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to remove pack size",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (r *router) addPackSize(w http.ResponseWriter, req *http.Request) {
	r.logger.Info().Msgf("Received request: %s %s", req.Method, req.URL.String())

	var input struct {
		Size int `json:"size"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil || input.Size <= 0 {
		r.logger.Error().Err(err).Msg("Error decoding request body")
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "size must be a positive integer",
		})
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		r.logger.Error().Err(err).Msg("Error decoding request body")
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "request body must contain a single JSON object",
		})
		return
	}

	r.logger.Debug().Msgf("Input: %v", input)
	if err := r.packSizes.Add(req.Context(), input.Size); err != nil {
		if errors.Is(err, packsizes.ErrAlreadyExists) {
			r.logger.Error().Err(err).Msg("Error creating pack size")
			writeJSON(w, http.StatusConflict, map[string]string{
				"error": "pack size already exists",
			})
			return
		}
		r.logger.Error().Err(err).Msg("Failed to add pack size")
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to add pack size",
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]int{"size": input.Size})
}

func (r *router) getPackSizes(w http.ResponseWriter, req *http.Request) {
	r.logger.Info().Msgf("Received request: %s %s", req.Method, req.URL.String())

	packSizes, err := r.packSizes.List(req.Context())
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to retrieve pack sizes")
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve pack sizes",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string][]int{"packSizes": packSizes})
}

func (r *router) getPacksForItemsOrdered(w http.ResponseWriter, req *http.Request) {
	r.logger.Info().Msgf("Received request: %s %s", req.Method, req.URL.String())
	r.logger.Debug().Msgf("Query parameters: %v", req.URL.Query())

	itemsOrdered, err := strconv.Atoi(req.URL.Query().Get("itemsOrdered"))
	if err != nil || itemsOrdered <= 0 {
		r.logger.Error().Err(err).Msg("Invalid itemsOrdered")
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "itemsOrdered must be a positive integer",
		})
		return
	}

	packSizes, err := r.packSizes.List(req.Context())
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to retrieve pack sizes")
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to retrieve pack sizes",
		})
		return
	}
	r.logger.Debug().Msgf("packSizes found: %v", packSizes)

	writeJSON(w, http.StatusOK, packs.GetPacksForOrder(r.logger, itemsOrdered, packSizes))
}

func (r *router) health(w http.ResponseWriter, _ *http.Request) {
	r.logger.Info().Msg("Received health check request")

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
