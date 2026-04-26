package http

import (
	"encoding/json"
	"net/http"
)

// createSnippetRequest defines the JSON structure we expect from the user.
type createSnippetRequest struct {
		Title	string	`json:"title"`
		OwnerID	string	`json:"owner_id"`
		Content string	`json:"content"`
}

func (h *Handler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
		var	req	createSnippetRequest

		// 1. Decode the JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request payload", http.StatusBadRequest)
				return
		}

		// 2. Hand off to the 'Brain' (Service Layer)
		snippet, err := h.svc.CreateSnippet(r.Context(), req.Title, req.OwnerID, req.Content)
		if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}

		// 3. Respond with the created identity
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(snippet)
}