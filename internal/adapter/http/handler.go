package http

import (
	"encoding/json"
	"net/http"

	"github.com/kevshouse/uber-sieben-brucken/internal/core"
)

// Handler is the 'Voice' of our bridge.
type Handler struct {
		service *core.SnippetService
}

func NewHandler(s *core.SnippetService) *Handler {
		return &Handler{service: s}
}

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
		snippet, err := h.service.CreateSnippet(r.Context(), req.Title, req.OwnerID, req.Content)
		if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}

		// 3. Respond with the created identity
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(snippet)
}

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
			err := h.service.CiteSnippet(r.Context(), req.SourceID, req.TargetID, req.Context)
			if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
			}

			// 3. Respond with success (204 No Content is standard for successful commands with no body)
			w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) SearchSnippets(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q") // Get the "q" parameter from the URL (e.g., /search?q=Euler)

		snippets, err := h.service.SearchSnippets(r.Context(), query) // Call the Brain
		if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}

		w.Header().Set("Content-Type", "application/json") // Responding with list of snippets found
		json.NewEncoder(w).Encode(snippets)
}