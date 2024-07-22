package main

import (
	"net/http"

	"log/slog"
)

func handlePing(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

func ApiRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", handlePing)

	return http.StripPrefix("/api/v1", mux)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/api/v1", ApiRoutes())
	mux.Handle("/", http.FileServer(http.Dir("./ui/build")))

	slog.Info("Server started at :4420")
	if err := http.ListenAndServe(":4420", mux); err != nil {
		slog.Error("Http Server Stopped due to", "err", err)
	}
}
