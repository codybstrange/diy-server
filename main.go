package main

import (
  "net/http"
)

func main() {
  newHandler := http.NewServeMux()
  server := http.Server{
    Handler: newHandler,
    Addr: ":8080",
  }
  server.ListenAndServe()
  return
}
