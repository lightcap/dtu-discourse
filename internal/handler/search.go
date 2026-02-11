package handler

import (
	"net/http"

	"github.com/lightcap/dtu-discourse/internal/store"
)

type SearchHandler struct {
	Store *store.Store
}

// GET /search
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		q = r.URL.Query().Get("term")
	}
	if q == "" {
		writeError(w, http.StatusBadRequest, "search term is required")
		return
	}
	result := h.Store.Search(q)
	writeJSON(w, http.StatusOK, result)
}
