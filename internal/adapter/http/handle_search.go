package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) SearchSnippets(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q") // Get the "q" parameter from the URL (e.g., /search?q=Euler)

		snippets, err := h.svc.SearchSnippets(r.Context(), query) // Call the Brain
		if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}

		w.Header().Set("Content-Type", "application/json") // Responding with list of snippets found
		json.NewEncoder(w).Encode(snippets)
}