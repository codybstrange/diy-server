package main

import (
  "encoding/json"
  "net/http"
  "github.com/google/uuid"
  "context"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Event string `json:"event"`
    Data struct {
      UserID string `json:"user_id"`
    } `json:"data"`
  }
  
  params := parameters{}
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusNoContent, "Could not decode body", err)
    return
  }
  if params.Event != "user.upgraded" {
    respondWithJSON(w, http.StatusNoContent, nil)
    return
  }
  
  userID, err := uuid.Parse(params.Data.UserID)
  if err != nil {
    respondWithError(w, http.StatusNoContent, "Couldn't parse user id", err)
    return
  }
  if err := cfg.db.UpgradeUser(context.Background(), userID); err != nil {
    respondWithError(w, http.StatusNotFound, "User could not be found", err)
    return
  }
  respondWithJSON(w, http.StatusNoContent, nil)
}
