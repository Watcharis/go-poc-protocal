package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/restful_api/ratelimit/models"
)

func (h *restFulAPIHandlers) CreateUserProfile(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
		logger.Info(ctx, "handler - CreateUserProfile")
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.ProifleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := pkg.ValidateStruct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := h.services.CreateUserProfile(ctx, req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w = pkg.SetContentType(w, "application/json")
		w = pkg.SetHttpStatusCode(w, http.StatusOK)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *restFulAPIHandlers) GetUserProfile(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx = r.Context()
		logger.Info(ctx, "handler - GetUserProfile")

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get path parameter from request /{uuid}
		// uuid := r.URL.Query().Get("uuid")
		uuid := r.Header.Get("uuid")
		if uuid == "" {
			http.Error(w, "uuid is required", http.StatusBadRequest)
			return
		}

		// get user profile
		result, err := h.services.GetUserProfile(ctx, uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w = pkg.SetContentType(w, "application/json")
		w = pkg.SetHttpStatusCode(w, http.StatusOK)

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
