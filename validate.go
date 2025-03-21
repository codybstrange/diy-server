package main

import (
  "net/http"
  "encoding/json"
  "strings"
)
const maxChars = 140

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > maxChars {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

  badWords := map[string] struct{}{
    "kerfuffle": {},
    "sharbert": {},
    "fornax": {},
  }
  cleaned := getCleanedBody(params.Body, badWords)
  respondWithJSON(w, http.StatusOK, returnVals{
    CleanedBody: cleaned,
	})
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
