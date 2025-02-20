package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/wolv89/chirpy/internal/database"
)

type WebhookResponse struct {
	Event string            `json:"event"`
	Data  map[string]string `json:"data"`
}

func (cfg *apiConfig) PolkaWebhook(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	webhookResp := WebhookResponse{}
	err := decoder.Decode(&webhookResp)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{"Something went wrong"})
		return
	}

	if webhookResp.Event != "user.upgraded" {
		responseJSON(w, http.StatusNoContent, nil)
		return
	}

	user_id := webhookResp.Data["user_id"]

	uuid, err := uuid.Parse(user_id)
	if err != nil {
		responseJSON(w, http.StatusBadRequest, ErrorResponse{"invalid user ID"})
		return
	}

	_, err = cfg.dbQueries.GetUser(req.Context(), uuid)
	if err != nil {
		responseJSON(w, http.StatusNotFound, ErrorResponse{"user not found"})
		return
	}

	_, err = cfg.dbQueries.SetUserChirpyRedStatus(req.Context(), database.SetUserChirpyRedStatusParams{
		ID:          uuid,
		IsChirpyRed: true,
	})
	if err != nil {
		responseJSON(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	responseJSON(w, http.StatusNoContent, nil)

}
