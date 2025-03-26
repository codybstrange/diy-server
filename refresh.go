package main

import (
  "github.com/codybstrange/diy-server/internal/auth"
  "context"
  "time"
  "net/http"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Issue authenticating user", err)
    return
  }
  
  refreshToken, err := cfg.db.GetRefreshToken(context.Background(), token)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
    return
  }
  if time.Now().After(refreshToken.ExpiresAt) || refreshToken.RevokedAt.Valid {
    respondWithError(w, http.StatusUnauthorized, "RefreshToken has expired", err)
    return
  }

  userID := refreshToken.UserID
  newToken, err := auth.MakeJWT(userID, cfg.tokenSecret)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Issue with creating new JWT", err)
    return
  }
  type response struct {
    Token string `json:"token"`
  } 
  respondWithJSON(w, http.StatusOK, response{Token: newToken,})
}
