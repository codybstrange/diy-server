package main

import (
  "context"
  "encoding/json"
  "net/http"
  "github.com/codybstrange/diy-server/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Password  string `json:"password"`
    Email     string `json:"email"`
  }
  
  decoder := json.NewDecoder(r.Body)
  params  := parameters{}
  if err  := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode the parameters", err)
    return
  }

  user, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
  if err != nil {
    respondWithError(w, http.StatusNotFound, "User not found", err)
    return
  }
  
  if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
    respondWithError(w, http.StatusUnauthorized, "Incorrect password or email", err)
  }
  
  u := User {
    ID: user.ID,
    CreatedAt: user.CreatedAt,
    UpdatedAt: user.UpdatedAt,
    Email: user.Email,
  }
  respondWithJSON(w, http.StatusOK, u)
}
