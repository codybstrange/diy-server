package main

import (
  "github.com/codybstrange/diy-server/internal/auth"
  "github.com/codybstrange/diy-server/internal/database"
  "net/http"
  "encoding/json"
  "context"
)

func (cfg *apiConfig) handlerUpdateDetails(w http.ResponseWriter, r *http.Request) {
  accessToken, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "No access header provided", err)
    return
  }

  type parameters struct {
    Password string `json:"password"`
    Email    string `json:"email"`
  }
  decoder := json.NewDecoder(r.Body)
  params  := parameters{}
  if err := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode the parameters", err)
    return
  }
  
  hash, err := auth.HashPassword(params.Password)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Password not secure enough", err)
    return
  }

  userID, err := auth.ValidateJWT(accessToken, cfg.tokenSecret)
  if err != nil {
    respondWithError(w, http.StatusUnauthorized, "Access token incorrect or missing", err)
    return
  }

  updateParams := database.UpdateUserPasswordParams{
    ID: userID,
    HashedPassword: hash,
    Email: params.Email,
  }
  newUser, err := cfg.db.UpdateUserPassword(context.Background(), updateParams)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't update the password", err)
    return
  }
  u := User{
    ID: newUser.ID,
    CreatedAt: newUser.CreatedAt,
    UpdatedAt: newUser.UpdatedAt,
    Email: newUser.Email,
    IsChirpyRed: newUser.IsChirpyRed,
  }

  respondWithJSON(w, http.StatusOK, u)
}
