package main

import (
  "net/http"
  "encoding/json"
  "strings"
  "github.com/google/uuid"
  "fmt"
  "time"
  "context"
  "github.com/codybstrange/diy-server/internal/database"
)
const maxChars = 140

type Chirp struct {
  ID        uuid.UUID `json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  Body      string    `json:"body"`
  UserID    uuid.UUID `json:"user_id"`
}

type ChirpReq struct {
  Body string `json:"body"`
  UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
  id := r.PathValue("id")
  if id == "" {
    respondWithError(w, http.StatusBadRequest, "No chirp ID provided", nil)
    return
  }

  parsedID, err := uuid.Parse(id)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Didn't recognize id as uuid", err)
  }
  chirp, err := cfg.db.GetChirp(context.Background(), parsedID)
  if err != nil {
    respondWithError(w, http.StatusNotFound, "Couldn't find chirp by id", err)
    return
  }
  c := Chirp {
    ID: chirp.ID,
    CreatedAt: chirp.CreatedAt,
    UpdatedAt: chirp.UpdatedAt,
    Body: chirp.Body,
    UserID: chirp.UserID,
  }
  respondWithJSON(w, http.StatusOK, c)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
  chirps, err := cfg.db.GetAllChirps(context.Background())
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Issue with fetching all chirps", err)
    return
  }
  
  output := []Chirp{}
  for _, c := range chirps {
    output = append(output, Chirp {
      ID: c.ID,
      CreatedAt: c.CreatedAt,
      UpdatedAt: c.UpdatedAt,
      Body: c.Body,
      UserID: c.UserID,
    })
  }
  respondWithJSON(w, http.StatusOK, output)
}

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
  chirpData, err := validateChirp(w, r)
  if err != nil {
    return
  }

  params := database.CreateChirpParams{
    Body: chirpData.Body,
    UserID: chirpData.UserID,
  }
  chirp, err := cfg.db.CreateChirp(context.Background(), params)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "Issue with creating chirp", err)
  }
  c := Chirp{
    ID: chirp.ID,
    CreatedAt: chirp.CreatedAt,
    UpdatedAt: chirp.UpdatedAt,
    Body: chirp.Body,
    UserID: chirp.UserID,
  }
  respondWithJSON(w, http.StatusCreated, c)
}



func validateChirp(w http.ResponseWriter, r *http.Request) (ChirpReq, error) {

	decoder := json.NewDecoder(r.Body)
  chirp := ChirpReq{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return ChirpReq{}, err
	}

	if len(chirp.Body) > maxChars {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return ChirpReq{}, fmt.Errorf("Chirp is too long")
	}

  badWords := map[string] struct{}{
    "kerfuffle": {},
    "sharbert": {},
    "fornax": {},
  }
  cleaned := getCleanedBody(chirp.Body, badWords)
  chirp.Body = cleaned
  return chirp, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
