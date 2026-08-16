package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/saddiqs1/gymshark-saddiqs1/internal/packs"
)

type router struct {
	logger *zerolog.Logger
}

func NewRouter(logger *zerolog.Logger) http.Handler {
	r := &router{logger: logger}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", r.health)
	mux.HandleFunc("GET /packs", r.getPacks)

	return mux
}

func (r *router) getPacks(w http.ResponseWriter, req *http.Request) {
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

	packSizesStr := strings.Split(req.URL.Query().Get("packSizes"), ",")
	packSizes := make([]int, len(packSizesStr))
	for i, s := range packSizesStr {
		v, err := strconv.Atoi(s)
		if err != nil || v <= 0 {
			r.logger.Error().Err(err).Msg("Invalid packSizes")
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "packSizes must be a list of positive integers",
			})
			return
		}
		packSizes[i] = v
	}

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
