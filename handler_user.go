package main

import (
  "net/http"
  "encoding/json"
  "github.com/google/uuid"
  "time"
  "context"
  "github.com/codybstrange/diy-server/internal/auth"
  "github.com/codybstrange/diy-server/internal/database"
)

type User struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
  type parameters struct {
    Password  string `json:"password"`
    Email     string `json:"email"`
  }
  decoder := json.NewDecoder(r.Body)
  params  := parameters{}
  if err  := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
    return
  }
  
  hash, err := auth.HashPassword(params.Password)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "Password not secure enough", err)
    return
  }

  userParams := database.CreateUserParams{
    HashedPassword: hash,
    Email: params.Email,
  }

  u, err := cfg.db.CreateUser(context.Background(), userParams)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Issue with creating user", err)
    return
  }
  
  user := User {
    u.ID,
    u.CreatedAt,
    u.UpdatedAt,
    u.Email,
  }
  respondWithJSON(w, http.StatusCreated, user)
}
