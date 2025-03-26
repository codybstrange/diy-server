package main

import (
  "github.com/codybstrange/diy-server/internal/auth"
  "github.com/codybstrange/diy-server/internal/database"
  "net/http"
  "database/sql"
  "time"
  "context"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Issue authenticating user", err)
    return
  }
  revokeParams := database.RevokeTokenParams{
    Token: token,
    RevokedAt: sql.NullTime{ Time: time.Now().UTC(), Valid: true, },
    UpdatedAt: time.Now().UTC(),
  }
  if err := cfg.db.RevokeToken(context.Background(), revokeParams); err != nil {
    respondWithError(w, http.StatusBadRequest, "Could not revoke token", err)
    return
  }
  
  respondWithJSON(w, http.StatusNoContent, nil)
}
