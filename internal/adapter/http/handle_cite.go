package http

import (
	"encoding/json"
	"net/http"
)

// citeSnippetRequest defines the JSON structure for linking two snippets.
type	citeSnippetRequest struct {
			SourceID string `json:"source_id"`
			TargetID string `json:"target_id"`
			Context string  `json:"context"`
}

func	(h *Handler) CiteSnippet(w http.ResponseWriter, r *http.Request) {
			var req citeSnippetRequest

			// 1. Decode the JSON body
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "Invalid request payload", http.StatusBadRequest)
					return
			}

			// 2. Hand off to the 'Brain' (Service Layer)
			err := h.svc.CiteSnippet(r.Context(), req.SourceID, req.TargetID, req.Context)
			if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}

			// 3. Respond with success (204 No Content is standard for successful commands with no body)
			w.WriteHeader(http.StatusNoContent)

}