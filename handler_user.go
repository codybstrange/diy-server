package main

import (
  "net/http"
  "encoding/json"
  "github.com/google/uuid"
  "time"
  "context"
)

type User struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
  type parameters struct {
    Email string `json:"email"`
  }
  decoder :=json.NewDecoder(r.Body)
  params := parameters{}
  if err := decoder.Decode(&params); err != nil {
    respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
  }

  u, err := cfg.db.CreateUser(context.Background(), params.Email)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Issue with creating user", err)
  }
  
  user := User {
    u.ID,
    u.CreatedAt,
    u.UpdatedAt,
    u.Email,
  }
  respondWithJSON(w, http.StatusCreated, user)
}
