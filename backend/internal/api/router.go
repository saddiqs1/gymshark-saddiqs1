package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/saddiqs1/gymshark-saddiqs1/internal/packs"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /packs", getPacks)

	return mux
}

func getPacks(w http.ResponseWriter, r *http.Request) {
	itemsOrdered, err := strconv.Atoi(r.URL.Query().Get("itemsOrdered"))
	if err != nil || itemsOrdered <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "itemsOrdered must be a positive integer",
		})
		return
	}

	writeJSON(w, http.StatusOK, packs.GetPacksForOrder(itemsOrdered))
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
