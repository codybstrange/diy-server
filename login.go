package main

import (
  "context"
  "encoding/json"
  "net/http"
  "github.com/codybstrange/diy-server/internal/auth"
  "time"
)

const maxTime = 3600

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
  type parameters struct {
    Password  string `json:"password"`
    Email     string `json:"email"`
    ExpiresIn int    `json:"expires_in_seconds"`
  }
  
  decoder := json.NewDecoder(r.Body)
  params  := parameters{}
  if err  := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode the parameters", err)
    return
  }

  if params.ExpiresIn == 0 || params.ExpiresIn > maxTime {
    params.ExpiresIn = maxTime
  }

  user, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
  if err != nil {
    respondWithError(w, http.StatusNotFound, "User not found", err)
    return
  }
  
  if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
    respondWithError(w, http.StatusUnauthorized, "Incorrect password or email", err)
    return
  }

  token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Duration(params.ExpiresIn) * time.Second)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Token could not be created", err)
    return
  }
  
  u := User {
    ID: user.ID,
    CreatedAt: user.CreatedAt,
    UpdatedAt: user.UpdatedAt,
    Email: user.Email,
    Token: token, 
  }
  respondWithJSON(w, http.StatusOK, u)
}
