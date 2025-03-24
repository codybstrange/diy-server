package main

import (
  "net/http"
  "context"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
  if cfg.platform != "dev" {
    respondWithError(w, http.StatusForbidden, "Do not have the permissions to delete all users", nil)
  }
  if err := cfg.db.DeleteAllUsers(context.Background()); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't delete all users", err)
  }
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

