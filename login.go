package main

import (
  "context"
  "encoding/json"
  "net/http"
  "github.com/codybstrange/diy-server/internal/auth"
  "github.com/codybstrange/diy-server/internal/database"
  "time"
)

const maxTime = 24*60

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
    return
  }

  token, err := auth.MakeJWT(user.ID, cfg.tokenSecret)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Token could not be created", err)
    return
  }

  refreshToken := auth.MakeRefreshToken()
  // Add to database
  expiresAt := time.Now().UTC().Add(maxTime * time.Hour)
  tokenParams := database.CreateRefreshTokenParams{
    Token: refreshToken,
    UserID: user.ID,
    ExpiresAt: expiresAt,
  }
  if err := cfg.db.CreateRefreshToken(context.Background(), tokenParams); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Refresh token couldn't be added to database", err)
    return
  }

  u := User {
    ID: user.ID,
    CreatedAt: user.CreatedAt,
    UpdatedAt: user.UpdatedAt,
    Email: user.Email,
    Token: token,
    RefreshToken: refreshToken, 
  }
  respondWithJSON(w, http.StatusOK, u)
}
