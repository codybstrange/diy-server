package main

import (
  "github.com/codybstrange/diy-server/internal/auth"
  "net/http"
  "context"
  "github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
  accessToken, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Access token missing", err)
    return
  }
  userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
  }
  
  chirpID := r.PathValue("id")
  if chirpID == "" {
    respondWithError(w, http.StatusBadRequest, "No chirp ID provided", nil)
    return
  }
  parsedID, err := uuid.Parse(chirpID)
  if err != nil{ 
    respondWithError(w, http.StatusInternalServerError, "Didn't recognize id as uuid", err)
    return
  }
  chirp, err := cfg.db.GetChirp(context.Background(), parsedID)
  if err != nil {
    respondWithError(w, http.StatusNotFound, "Couldn't find chirp by id", err)
    return
  }

  if chirp.UserID != userID {
    respondWithError(w, http.StatusForbidden, "You are not the user who created this chirp", err)
    return
  }
  
  if err := cfg.db.DeleteChirp(context.Background(), parsedID); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
    return
  } 
  
  respondWithJSON(w, http.StatusNoContent, nil)
}
